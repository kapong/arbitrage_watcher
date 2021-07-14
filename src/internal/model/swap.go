package model

type Swap struct {
	Contract    string
	InputToken  Token
	InputAmount uint64
	OutputToken Token
	URL         string

	ExpectedProfit      float64
	ExpectedProfileStep float64

	ReturnAmount     uint64
	SpreadAmount     uint64
	CommissionAmount uint64
	Changed          float64
}

type Token struct {
	Code string
	Name string
}
