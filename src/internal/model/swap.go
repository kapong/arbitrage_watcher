package model

type Swap struct {
	Contract string
	InputToken string
	InputAmount uint64
	OutputToken string
	URL string

	ReturnAmount uint64
	SpreadAmount uint64
	CommissionAmount uint64
	Changed float64
}