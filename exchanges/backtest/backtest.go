package backtest

import (
	"context"

	"github.com/ydm/commons"
)

type Backtest struct {
	archives []string
}

func New(archives []string) *Backtest {
	return &Backtest{archives}
}

// [1] API calls.

// CreateOrder accepts arguments of the following formats:
//
// - side: buy or sell,
// - orderType: market (just that for now).
func (b *Backtest) CreateOrder(
	symbol string,
	side string,
	orderType string,
	priceStr string,
	quantityStr string,
	clientOrderID string,
	reduceOnly bool,
) (commons.CreateOrderResponse, error) {
	panic("not implemented")
}

func (b *Backtest) CancelOrder(symbol string, clientOrderID string) error {
	panic("not implemented")
}

// ChangeMarginType should be invoked with marginType set to
// "crossed" or "isolated".
func (b *Backtest) ChangeMarginType(symbol, marginType string) error {
	panic("not implemented")
}

func (b *Backtest) ChangeLeverage(symbol string, leverage int) error {
	panic("not implemented")
}

// [2] Streams.

func (b *Backtest) Book1(ctx context.Context, symbol string) <-chan commons.Book1 {
	panic("not implemented")
}

func (b *Backtest) Trade(ctx context.Context, symbol string) <-chan commons.Trade {
	return ReadTrades(b.archives, symbol)
}

// [3] Events.

func (b *Backtest) OnOrderUpdate(clientOrderID string, callback commons.OrderUpdateCallback) {
	panic("not implemented")
}
