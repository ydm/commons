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

	// Indicates whether the buyer is also the maker. Make sure
	// given exchange supports that before using it.
	BuyerIsMaker bool
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
	// [1] API calls.

	// CreateOrder accepts arguments of the following formats:
	//
	// - side: buy or sell,
	// - orderType: market (just that for now).
	CreateOrder(symbol, side, orderType, quantityStr string, reduceOnly bool) (CreateOrderResponse, error)

	// ChangeMarginType should be invoked with marginType set to
	// "crossed" or "isolated".
	ChangeMarginType(symbol, marginType string) error
	ChangeLeverage(symbol string, leverage int) error

	// [2] Streams.
	Book1(ctx context.Context, symbol string) chan Book1
	Trade(ctx context.Context, symbol string) chan Trade
}
