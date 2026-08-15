package domain

import (
	"fmt"
	"strings"
)

type Account struct {
	ID             string `json:"id"`
	Plan           string `json:"plan"`
	IncludedUnits  int64  `json:"included_units"`
	UnitPriceCents int64  `json:"unit_price_cents"`
}

func (a Account) Validate() error {
	if strings.TrimSpace(a.ID) == "" {
		return fmt.Errorf("account id is required")
	}
	if strings.TrimSpace(a.Plan) == "" {
		return fmt.Errorf("account %q plan is required", a.ID)
	}
	if a.IncludedUnits < 0 {
		return fmt.Errorf("account %q included units cannot be negative", a.ID)
	}
	if a.UnitPriceCents < 0 {
		return fmt.Errorf("account %q unit price cannot be negative", a.ID)
	}
	return nil
}
