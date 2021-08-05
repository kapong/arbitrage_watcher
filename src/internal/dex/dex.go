package dex

import (
	"arbitrage_monitoring/internal/model"
)

type DEXInterface interface {
	GetPrice(swap *model.Swap) (*model.Swap, error)
}

type DEX struct {
	DEXInterface
}
