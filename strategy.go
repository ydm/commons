package commons

import "github.com/adshao/go-binance/v2/futures"

type Strategy interface {
	OnAggTrade(event *futures.WsAggTradeEvent)
	OnBookTicker(event *futures.WsBookTickerEvent)
}
