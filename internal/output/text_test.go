package output

import (
	"bytes"
	"errors"
	"testing"

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

func TestWriteTextReturnsDetailWriteError(t *testing.T) {
	detailErr := errors.New("detail sink unavailable")
	report := domain.Report{Accounts: []domain.AccountReport{
		{AccountID: "acme", Plan: "starter", UsedUnits: 125, IncludedUnits: 100, OverageUnits: 25, ChargeCents: 175},
		{AccountID: "beta", Plan: "team", UsedUnits: 4, IncludedUnits: 10, OverageUnits: 0, ChargeCents: 0},
	}}

	writer := failOnWriteWriter{failAt: 2, err: detailErr}
	err := WriteText(&writer, report)
	if !errors.Is(err, detailErr) {
		t.Fatalf("WriteText() error = %v, want %v", err, detailErr)
	}
	if got, want := writer.writes, 2; got != want {
		t.Fatalf("writes = %d, want %d; WriteText should stop after the failed account detail", got, want)
	}
}

type failOnWriteWriter struct {
	writes int
	failAt int
	err    error
}

func (w *failOnWriteWriter) Write(p []byte) (int, error) {
	w.writes++
	if w.writes == w.failAt {
		return 0, w.err
	}
	return len(p), nil
}
