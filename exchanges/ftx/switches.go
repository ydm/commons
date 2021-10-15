package ftx

func switchSideString(side, buy, sell string) string {
	switch side {
	case "buy":
		return buy
	case "sell":
		return sell
	default:
		panic("")
	}
}

func switchTypeString(orderType, market string) string {
	switch orderType {
	case "market":
		return market
	default:
		panic("")
	}
}
