package binancef

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/rs/zerolog/log"
	"github.com/ydm/commons"
)

type Trade struct {
	LastTradeID int64
	Time        time.Time
	Symbol      string
	Price       float64
	Quantity    float64
	TradeTime   time.Time
	Maker       bool // Whether buyer is maker.
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

func SubscribeAggTrade(ctx context.Context, symbol string) <-chan commons.Trade {
	c := make(chan commons.Trade)

	// In a parallel goroutine, start a loop that subscribes to WsAggTrade and feeds c.
	// If an error occurs, start over.
	go func() {
		// Using this as MT-safe arrays.
		stops := make(chan chan struct{}, 256)
		dones := make(chan chan struct{}, 256)

		// Make sure all stops are called and all dones are waited after context.
		go func() {
			commons.Checker.Push()
			defer commons.Checker.Pop()

			<-ctx.Done()

			for stop := range stops {
				close(stop)
			}

			for done := range dones {
				<-done
			}
		}()

		for ctx.Err() == nil {
			commons.What(log.Info().Str("symbol", symbol), "subscribing to WsAggTrade")

			done, stop, err := futures.WsAggTradeServe(
				symbol,
				func(event *futures.WsAggTradeEvent) {
					if event == nil {
						commons.What(log.Warn(), "event is nil")

						return
					}

					var x Trade
					if err := commons.SmartCopy(&x, event); err != nil {
						commons.What(log.Warn().Err(err), "SmartCopy(WsAggTradeEvent) failed")
					}

					c <- commons.Trade{
						TradeID:      x.LastTradeID,
						Time:         x.Time,
						Symbol:       x.Symbol,
						Price:        x.Price,
						Quantity:     x.Quantity,
						BuyerIsMaker: x.Maker,
					}
				},
				func(err error) {
					commons.What(log.Warn().Str("symbol", symbol).Err(err), "WsAggTradeServe invoked error handler")
				},
			)

			if err == nil {
				dones <- done
				stops <- stop
			} else {
				commons.What(log.Warn().Str("symbol", symbol).Err(err), "WsAggTradeServe failed")
			}

			<-done
		}
	}()

	return c
}

func SubscribeBookTicker(ctx context.Context, symbol string) <-chan commons.Book1 {
	c := make(chan commons.Book1)

	done, stop, err := futures.WsBookTickerServe(
		symbol,
		func(event *futures.WsBookTickerEvent) {
			var x BookTicker
			if err := commons.SmartCopy(&x, event); err != nil {
				commons.What(log.Warn().Err(err), "SmartCopy(WsBookTickerEvent) failed")
			}

			c <- commons.Book1{
				Time:        x.Time,
				Symbol:      x.Symbol,
				BidPrice:    x.BestBidPrice,
				BidQuantity: x.BestBidQty,
				AskPrice:    x.BestAskPrice,
				AskQuantity: x.BestAskQty,
			}
		},
		func(err error) {
			panic(err)
		},
	)
	if err != nil {
		panic(err)
	}

	go func() {
		commons.Checker.Push()
		defer commons.Checker.Pop()

		<-done
		close(c)
	}()

	go func() {
		commons.Checker.Push()
		defer commons.Checker.Pop()

		<-ctx.Done()
		close(stop)
		<-done
	}()

	return c
}
