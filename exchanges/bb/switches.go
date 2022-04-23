package bb

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

func SwitchSideString(side string, buy, sell string) string {
	switch side {
	case "buy":
		return buy
	case "sell":
		return sell
	default:
		panic(side)
	}
}

func SwitchTypeString(orderType, market, limit, stopMarket string) string {
	switch orderType {
	case OrderTypeMarket:
		return market
	case OrderTypeLimit:
		return limit
	case OrderTypeStopMarket:
		return stopMarket
	}

	panic(orderType)
}

func SwitchMarginTypeString(marginType, crossed, isolated string) string {
	switch marginType {
	case "crossed":
		return crossed
	case "isolated":
		return isolated
	default:
		panic(marginType)
	}
}
