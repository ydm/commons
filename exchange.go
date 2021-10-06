package commons

import (
	"context"
	"time"
)

type Trade struct {
	Time     time.Time
	Symbol   string
	Price    float64
	Quantity float64
}

// Book1 is the level one order book data.
type Book1 struct {
	Time        time.Time
	Symbol      string
	BidPrice    float64
	BidQuantity float64
	AskPrice    float64
	AskQuantity float64
}

type CreateOrderResponse struct {
	ExecutedQuantity string
	AvgPrice         string
}

type Exchange interface {
	// API calls.
	CreateOrder(symbol, side, orderType, quantityStr string, reduceOnly bool) (CreateOrderResponse, error)

	// Streams.
	Book1(ctx context.Context, symbol string) chan Book1
	Trade(ctx context.Context, symbol string) chan Trade
}
