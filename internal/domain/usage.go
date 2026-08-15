package domain

import (
	"fmt"
	"strings"
)

type UsageEvent struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Units     int64  `json:"units"`
}

func (u UsageEvent) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return fmt.Errorf("usage event id is required")
	}
	if strings.TrimSpace(u.AccountID) == "" {
		return fmt.Errorf("usage event %q account id is required", u.ID)
	}
	if u.Units < 0 {
		return fmt.Errorf("usage event %q units cannot be negative", u.ID)
	}
	return nil
}
