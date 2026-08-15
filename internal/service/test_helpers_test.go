package service_test

import "github.com/zhangchengcheng/subscription-usage-cli/internal/domain"

type storeAccount domain.Account
type storeUsage domain.UsageEvent

type storeAccounts []storeAccount
type storeUsageEvents []storeUsage

func (accounts storeAccounts) domainAccounts() []domain.Account {
	result := make([]domain.Account, len(accounts))
	for i := range accounts {
		result[i] = domain.Account(accounts[i])
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
