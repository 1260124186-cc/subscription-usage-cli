package service_test

import "github.com/1260124186-cc/subscription-usage-cli/internal/domain"

type storeAccount domain.Account
type storeUsage domain.UsageEvent

type storeAccounts []storeAccount
type storeUsageEvents []storeUsage

func (accounts storeAccounts) domainAccounts() []*domain.Account {
	result := make([]*domain.Account, len(accounts))
	for i := range accounts {
		account := domain.Account(accounts[i])
		result[i] = &account
	}
	return result
}

func (usage storeUsageEvents) domainUsage() []domain.UsageEvent {
	result := make([]domain.UsageEvent, len(usage))
	for i := range usage {
		result[i] = domain.UsageEvent(usage[i])
	}
	return result
}
