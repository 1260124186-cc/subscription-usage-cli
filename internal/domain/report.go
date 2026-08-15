package domain

type AccountReport struct {
	AccountID     string
	Plan          string
	UsedUnits     int64
	IncludedUnits int64
	OverageUnits  int64
	ChargeCents   int64
}

type Report struct {
	Accounts         []AccountReport
	TotalChargeCents int64
}
