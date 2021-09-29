package commons

import (
	"time"
)

type Trade struct {
	Time      time.Time
	Symbol    string
	Price     float64
	Quantity  float64
	TradeTime time.Time
}

type BookTicker struct {
	Time            time.Time
	TransactionTime time.Time
	Symbol          string
	BestBidPrice    float64
	BestBidQty      float64
	BestAskPrice    float64
	BestAskQty      float64
}

type Strategy interface {
	OnAggTrade(symbol string, event Trade)
	OnBookTicker(symbol string, event BookTicker)
}
