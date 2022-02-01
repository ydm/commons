package ftx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-numb/go-ftx/auth"
	"github.com/go-numb/go-ftx/realtime"
	"github.com/go-numb/go-ftx/rest"
	"github.com/go-numb/go-ftx/rest/private/account"
	"github.com/go-numb/go-ftx/rest/private/orders"
	"github.com/go-numb/go-ftx/types"
	"github.com/ydm/commons"
)

type FTX struct {
	restClient *rest.Client
}

func New(apikey, secret string) commons.Exchange {
	return &FTX{
		restClient: rest.New(auth.New(apikey, secret)),
	}
}

func (f *FTX) CreateOrder(symbol, side, orderType, quantityStr string, reduceOnly bool) (
	ans commons.CreateOrderResponse,
	err error,
) {
	var (
		quantity float64
		resp     *orders.ResponseForPlaceOrder
	)

	if quantity, err = strconv.ParseFloat(quantityStr, 64); err != nil {
		return
	}

	resp, err = f.restClient.PlaceOrder(&orders.RequestForPlaceOrder{
		ClientID:   "",
		Type:       switchTypeString(orderType, types.MARKET),
		Market:     symbol,
		Side:       switchSideString(side, types.BUY, types.SELL),
		Price:      0,
		Size:       quantity,
		ReduceOnly: reduceOnly,
		Ioc:        false,
		PostOnly:   false,
	})

	ans.AvgPrice = strconv.FormatFloat(resp.Price, 'f', 8, 64)
	ans.ExecutedQuantity = strconv.FormatFloat(resp.FilledSize, 'f', 8, 64)

	return
}

func (f *FTX) ChangeMarginType(symbol, marginType string) error {
	return nil
}

func (f *FTX) ChangeLeverage(symbol string, leverage int) error {
	_, err := f.restClient.Leverage(&account.RequestForLeverage{
		Leverage: leverage,
	})
	if err != nil {
		return fmt.Errorf("ChangeLeverage failed: %w", err)
	}

	return nil
}

func (f *FTX) Book1(ctx context.Context, symbol string) chan commons.Book1 {
	inp := make(chan realtime.Response)
	out := make(chan commons.Book1)

	go func() {
		commons.Checker.Push()
		defer commons.Checker.Push()

		err := realtime.Connect(ctx, inp, []string{"ticker"}, []string{symbol}, nil)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		commons.Checker.Push()
		defer commons.Checker.Push()
		defer close(out)

		for resp := range inp {
			if resp.Type == realtime.TICKER {
				out <- commons.Book1{
					Time:        resp.Ticker.Time.Time,
					Symbol:      symbol,
					BidPrice:    resp.Ticker.Bid,
					BidQuantity: resp.Ticker.BidSize,
					AskPrice:    resp.Ticker.Ask,
					AskQuantity: resp.Ticker.AskSize,
				}
			}
		}
	}()

	return out
}

func (f *FTX) Trade(ctx context.Context, symbol string) chan commons.Trade {
	panic("not implemented")
}
