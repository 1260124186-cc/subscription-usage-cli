package output

import (
	"context"
	"fmt"
	"io"

	"github.com/1260124186-cc/subscription-usage-cli/internal/domain"
)

// ContextWriter stops an in-progress write when its context is canceled.
type ContextWriter interface {
	io.Writer
	WriteContext(context.Context, []byte) (int, error)
}

type contextWriteCloser struct {
	io.WriteCloser
}

func NewContextWriteCloser(writer io.WriteCloser) ContextWriter {
	return contextWriteCloser{WriteCloser: writer}
}

func (w contextWriteCloser) WriteContext(ctx context.Context, p []byte) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	type result struct {
		n   int
		err error
	}
	done := make(chan result, 1)
	go func() {
		n, err := w.Write(p)
		done <- result{n: n, err: err}
	}()

	select {
	case result := <-done:
		return result.n, result.err
	case <-ctx.Done():
		// Closing the endpoint lets a blocked write return before we wait for it.
		_ = w.Close()
		<-done
		return 0, ctx.Err()
	}
}

func WriteText(writer io.Writer, report domain.Report) error {
	return WriteTextContext(context.Background(), writer, report)
}

func WriteTextContext(ctx context.Context, writer io.Writer, report domain.Report) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for _, account := range report.Accounts {
		if err := writefContext(
			ctx,
			writer,
			"account=%s plan=%s used=%d included=%d overage=%d charge_cents=%d\n",
			account.AccountID,
			account.Plan,
			account.UsedUnits,
			account.IncludedUnits,
			account.OverageUnits,
			account.ChargeCents,
		); err != nil {
			return err
		}
	}
	return writefContext(ctx, writer, "total_charge_cents=%d\n", report.TotalChargeCents)
}

func writefContext(ctx context.Context, writer io.Writer, format string, args ...any) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	text := fmt.Sprintf(format, args...)
	if contextWriter, ok := writer.(ContextWriter); ok {
		n, err := contextWriter.WriteContext(ctx, []byte(text))
		if err != nil {
			return err
		}
		if n != len(text) {
			return io.ErrShortWrite
		}
	} else {
		if _, err := io.WriteString(writer, text); err != nil {
			return err
		}
	}
	return ctx.Err()
}
