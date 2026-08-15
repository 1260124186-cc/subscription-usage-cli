package output

import (
	"bytes"
	"context"
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
