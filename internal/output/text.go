package output

import (
	"context"
	"fmt"
	"io"

	"github.com/1260124186-cc/subscription-usage-cli/internal/domain"
)

func WriteText(writer io.Writer, report domain.Report) error {
	return WriteTextContext(context.Background(), writer, report)
}

func WriteTextContext(ctx context.Context, writer io.Writer, report domain.Report) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for _, account := range report.Accounts {
		if _, err := fmt.Fprintf(
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
	_, err := fmt.Fprintf(writer, "total_charge_cents=%d\n", report.TotalChargeCents)
	return err
}
