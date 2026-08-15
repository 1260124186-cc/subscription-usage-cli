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
	report := domain.Report{Accounts: []domain.AccountReport{{
		AccountID: "acme", Plan: "starter", UsedUnits: 125,
		IncludedUnits: 100, OverageUnits: 25, ChargeCents: 175,
	}}}

	if err := WriteText(&failFirstWriter{}, report); err == nil {
		t.Fatal("WriteText() error = nil, want the account write error")
	}
}

type failFirstWriter struct {
	writes int
}

func (w *failFirstWriter) Write(p []byte) (int, error) {
	w.writes++
	if w.writes == 1 {
		return 0, errors.New("detail sink unavailable")
	}
	return len(p), nil
}
