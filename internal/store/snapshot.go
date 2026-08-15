package store

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/zhangchengcheng/subscription-usage-cli/internal/domain"
)

type Snapshot struct {
	Accounts []domain.Account    `json:"accounts"`
	Usage    []domain.UsageEvent `json:"usage"`
}

func LoadSnapshot(reader io.Reader) (Snapshot, error) {
	var snapshot Snapshot
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&snapshot); err != nil {
		return Snapshot{}, fmt.Errorf("decode snapshot: %w", err)
	}
	if err := snapshot.Validate(); err != nil {
		return Snapshot{}, err
	}
	return snapshot, nil
}

func (s Snapshot) Validate() error {
	knownAccounts := make(map[string]struct{}, len(s.Accounts))
	for _, account := range s.Accounts {
		if err := account.Validate(); err != nil {
			return err
		}
		if _, exists := knownAccounts[account.ID]; exists {
			return fmt.Errorf("duplicate account %q", account.ID)
		}
		knownAccounts[account.ID] = struct{}{}
	}
	for _, usage := range s.Usage {
		if err := usage.Validate(); err != nil {
			return err
		}
		if _, exists := knownAccounts[usage.AccountID]; !exists {
			return fmt.Errorf("usage event %q references unknown account %q", usage.ID, usage.AccountID)
		}
	}
	return nil
}
