package commons

import "github.com/adshao/go-binance/v2/futures"

type FuturesStrategy interface {
	OnAggTrade(symbol string, event *futures.WsAggTradeEvent)
	OnBookTicker(symbol string, event *futures.WsBookTickerEvent)
}
