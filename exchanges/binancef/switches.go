package binancef

// func switchSideFloat64(side string, buy, sell float64) float64 {
// 	switch side {
// 	case "buy":
// 		return buy
// 	case "sell":
// 		return sell
// 	default:
// 		panic("invalid side")
// 	}
// }

func switchSideString(side string, buy, sell string) string {
	switch side {
	case "buy":
		return buy
	case "sell":
		return sell
	default:
		panic(side)
	}
}

func switchTypeString(orderType, market, limit, stopMarket string) string {
	switch orderType {
	case orderTypeMarket:
		return market
	case orderTypeLimit:
		return limit
	case orderTypeStopMarket:
		return stopMarket
	}

	panic(orderType)
}

func switchMarginTypeString(marginType, crossed, isolated string) string {
	switch marginType {
	case "crossed":
		return crossed
	case "isolated":
		return isolated
	default:
		panic(marginType)
	}
}
