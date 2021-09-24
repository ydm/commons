package commons

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/rs/zerolog/log"
)

func SubscribeAggTrade(ctx context.Context, symbol string, s Strategy) {
	done, stop, err := futures.WsAggTradeServe(
		symbol,
		func(event *futures.WsAggTradeEvent) {
			var x Trade
			if err := SmartCopy(&x, event); err != nil {
				What(log.Warn().Err(err), "SmartCopy(WsAggTradeEvent) failed")
			}

			s.OnAggTrade(symbol, x)
		},
		func(err error) {
			panic(err)
		},
	)
	if err != nil {
		panic(err)
	}

	go func() {
		CheckerPush()

		defer CheckerPop()

		<-ctx.Done()
		close(stop)
		<-done
	}()
}

func SubscribeBookTicker(ctx context.Context, symbol string, s Strategy) {
	done, stop, err := futures.WsBookTickerServe(
		symbol,
		func(event *futures.WsBookTickerEvent) {
			var x BookTicker
			if err := SmartCopy(&x, event); err != nil {
				What(log.Warn().Err(err), "SmartCopy(WsBookTickerEvent) failed")
			}

			s.OnBookTicker(symbol, x)
		},
		func(err error) {
			panic(err)
		},
	)
	if err != nil {
		panic(err)
	}

	go func() {
		CheckerPush()

		defer CheckerPop()

		<-ctx.Done()
		close(stop)
		<-done
	}()
}
