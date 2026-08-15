package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/1260124186-cc/subscription-usage-cli/internal/domain"
	"github.com/1260124186-cc/subscription-usage-cli/internal/store"
)

type ReportGenerator struct{}

func NewReportGenerator() ReportGenerator {
	return ReportGenerator{}
}

func (ReportGenerator) Generate(ctx context.Context, snapshot store.Snapshot) (domain.Report, error) {
	if err := ctx.Err(); err != nil {
		return domain.Report{}, err
	}

	sort.Slice(snapshot.Usage, func(i, j int) bool {
		return snapshot.Usage[i].ID < snapshot.Usage[j].ID
	})

	accounts := make(map[string]domain.Account, len(snapshot.Accounts))
	for _, account := range snapshot.Accounts {
		accounts[account.ID] = account
	}

	usedByAccount := make(map[string]int64, len(accounts))
	for _, usage := range snapshot.Usage {
		if err := ctx.Err(); err != nil {
			return domain.Report{}, err
		}
		usedByAccount[usage.AccountID] += usage.Units
	}

	reports := make([]domain.AccountReport, 0, len(snapshot.Accounts))
	for _, account := range snapshot.Accounts {
		used := usedByAccount[account.ID]
		overage := used - account.IncludedUnits
		if overage < 0 {
			overage = 0
		}
		reports = append(reports, domain.AccountReport{
			AccountID:     account.ID,
			Plan:          account.Plan,
			UsedUnits:     used,
			IncludedUnits: account.IncludedUnits,
			OverageUnits:  overage,
			ChargeCents:   overage * account.UnitPriceCents,
		})
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].AccountID < reports[j].AccountID
	})

	report := domain.Report{Accounts: reports}
	for _, accountReport := range report.Accounts {
		if accountReport.ChargeCents < 0 {
			return domain.Report{}, fmt.Errorf("negative charge for account %q", accountReport.AccountID)
		}
		report.TotalChargeCents += accountReport.ChargeCents
	}
	return report, nil
}
