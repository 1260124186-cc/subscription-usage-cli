package output

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/1260124186-cc/subscription-usage-cli/internal/domain"
)

func TestWriteText(t *testing.T) {
	var output bytes.Buffer
	report := domain.Report{
		Accounts: []domain.AccountReport{{
			AccountID: "acme", Plan: "starter", UsedUnits: 125,
			IncludedUnits: 100, OverageUnits: 25, ChargeCents: 175,
		}},
		TotalChargeCents: 175,
	}

	if err := WriteText(&output, report); err != nil {
		t.Fatalf("WriteText() error = %v", err)
	}
	const want = "account=acme plan=starter used=125 included=100 overage=25 charge_cents=175\ntotal_charge_cents=175\n"
	if got := output.String(); got != want {
		t.Fatalf("WriteText() = %q, want %q", got, want)
	}
}

func TestWriteTextContextStopsAfterCancellation(t *testing.T) {
	report := domain.Report{Accounts: []domain.AccountReport{{
		AccountID: "acme", Plan: "starter", UsedUnits: 125,
		IncludedUnits: 100, OverageUnits: 25, ChargeCents: 175,
	}}}

	if err := WriteTextContext(&cancelAfterInitialCheck{}, &bytes.Buffer{}, report); err == nil {
		t.Fatal("WriteTextContext() error = nil, want context cancellation")
	}
}

func TestContextWriteCloserStopsActiveWriteAtDeadline(t *testing.T) {
	report := domain.Report{Accounts: []domain.AccountReport{{
		AccountID: "acme", Plan: "starter", UsedUnits: 125,
		IncludedUnits: 100, OverageUnits: 25, ChargeCents: 175,
	}}}
	writer := newBlockingWriteCloser()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	err := WriteTextContext(ctx, NewContextWriteCloser(writer), report)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("WriteTextContext() error = %v, want context deadline exceeded", err)
	}
	select {
	case <-writer.exited:
	default:
		t.Fatal("write was still running after WriteTextContext() returned")
	}
	if writer.completed {
		t.Fatal("write completed after the context deadline")
	}
}

type cancelAfterInitialCheck struct {
	checks int
}

func (c cancelAfterInitialCheck) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (c cancelAfterInitialCheck) Done() <-chan struct{} {
	return nil
}

func (c *cancelAfterInitialCheck) Err() error {
	c.checks++
	if c.checks > 1 {
		return context.Canceled
	}
	return nil
}

func (c cancelAfterInitialCheck) Value(any) any {
	return nil
}

type blockingWriteCloser struct {
	closeOnce sync.Once
	closed    chan struct{}
	exited    chan struct{}
	completed bool
}

func newBlockingWriteCloser() *blockingWriteCloser {
	return &blockingWriteCloser{
		closed: make(chan struct{}),
		exited: make(chan struct{}),
	}
}

func (w *blockingWriteCloser) Write(p []byte) (int, error) {
	defer close(w.exited)

	timer := time.NewTimer(200 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-w.closed:
		return 0, errors.New("writer closed")
	case <-timer.C:
		w.completed = true
		return len(p), nil
	}
}

func (w *blockingWriteCloser) Close() error {
	w.closeOnce.Do(func() {
		close(w.closed)
	})
	return nil
}
