package commons

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
)

func AggTrade(ctx context.Context, symbol string, s FuturesStrategy) {
	done, stop, err := futures.WsAggTradeServe(
		symbol,
		func(event *futures.WsAggTradeEvent) {
			s.OnAggTrade(symbol, event)
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

func BookTicker(ctx context.Context, symbol string, s FuturesStrategy) {
	done, stop, err := futures.WsBookTickerServe(
		symbol,
		func(event *futures.WsBookTickerEvent) {
			s.OnBookTicker(symbol, event)
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
