package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/1260124186-cc/subscription-usage-cli/internal/domain"
	"github.com/1260124186-cc/subscription-usage-cli/internal/service"
	"github.com/1260124186-cc/subscription-usage-cli/internal/store"
)

func TestGenerateCalculatesIncludedAndOverageUsage(t *testing.T) {
	snapshot := store.Snapshot{
		Accounts: storeAccounts{
			{ID: "acme", Plan: "starter", IncludedUnits: 100, UnitPriceCents: 7},
			{ID: "beta", Plan: "team", IncludedUnits: 10, UnitPriceCents: 15},
		}.domainAccounts(),
		Usage: storeUsageEvents{
			{ID: "evt-1", AccountID: "acme", Units: 125},
			{ID: "evt-2", AccountID: "beta", Units: 4},
		}.domainUsage(),
	}

	report, err := service.NewReportGenerator().Generate(context.Background(), snapshot)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got, want := report.TotalChargeCents, int64(175); got != want {
		t.Fatalf("total charge = %d, want %d", got, want)
	}
	if got, want := report.Accounts[0].AccountID, "acme"; got != want {
		t.Fatalf("first account = %q, want %q", got, want)
	}
	if got, want := report.Accounts[0].OverageUnits, int64(25); got != want {
		t.Fatalf("acme overage = %d, want %d", got, want)
	}
}

func TestGenerateRejectsNilAccount(t *testing.T) {
	snapshot := store.Snapshot{Accounts: []*domain.Account{nil}}

	_, err := service.NewReportGenerator().Generate(context.Background(), snapshot)
	if err == nil {
		t.Fatal("Generate() error = nil, want an error for a nil account")
	}
	if got, want := err.Error(), "validate snapshot: account at index 0 is null"; got != want {
		t.Fatalf("Generate() error = %q, want %q", got, want)
	}
}

func TestGenerateRejectsUsageForUnknownAccount(t *testing.T) {
	snapshot := store.Snapshot{
		Accounts: storeAccounts{
			{ID: "acme", Plan: "starter", IncludedUnits: 100, UnitPriceCents: 7},
		}.domainAccounts(),
		Usage: storeUsageEvents{
			{ID: "evt-1", AccountID: "missing", Units: 1},
		}.domainUsage(),
	}

	_, err := service.NewReportGenerator().Generate(context.Background(), snapshot)
	if err == nil {
		t.Fatal("Generate() error = nil, want an error for usage with an unknown account")
	}
	if !strings.Contains(err.Error(), `usage event "evt-1" references unknown account "missing"`) {
		t.Fatalf("Generate() error = %q, want an unknown-account error", err)
	}
}
