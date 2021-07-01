package dex

import (
	"arbitrage_monitoring/internal/model"
)

type DEXInterface interface {
	GetPrice(swap *model.Swap) *model.Swap
}

type DEX struct {
	DEXInterface
}