package commons

import "github.com/rs/zerolog/log"

func oppositeSide(side string) string {
	switch side {
	case "buy":
		return "sell"
	case "sell":
		return "buy"
	default:
		panic("TODO")
	}
}

func create(exchange Exchange, symbol, side, orderType, price, quantity, id string) {
	resp, err := exchange.CreateOrder(symbol, side, orderType, price, quantity, id, true)

	if err != nil {
		What(log.Warn().Err(err), "create order failed")
	} else {
		What(
			log.Info().
				Str("symbol", symbol).
				Str("side", side).
				Str("type", orderType).
				Str("price", price).
				Str("quantity", quantity).
				Str("clientOrderID", id).
				Interface("resp", resp),
			"created aTPSL order",
		)
	}
}

func TPSL(
	exchange Exchange,
	symbol string,
	openSide string,
	quantityStr string,
	takeProfitPrice,
	stopLossPrice string,
) {
	closeSide := oppositeSide(openSide)
	takeProfitID := RandomOrderID("")
	stopLossID := RandomOrderID("")

	// First make sure we handle order filled events.

	exchange.OnOrderUpdate(takeProfitID, func(update OrderUpdate) {
		if update.Status != "FILLED" {
			return
		}

		What(
			log.Info().Interface("update", update).Str("stopLossID", stopLossID),
			"take profit order filled, will now cancel stop loss",
		)

		if err := exchange.CancelOrder(symbol, stopLossID); err != nil {
			What(log.Warn().Err(err), "failed to cancel stop_loss after take profit got executed")
		}
	})

	exchange.OnOrderUpdate(stopLossID, func(update OrderUpdate) {
		if update.Status != "FILLED" {
			return
		}

		What(
			log.Info().Interface("update", update).Str("takeProfitID", takeProfitID),
			"stop loss order filled, will now cancel take profit",
		)

		if err := exchange.CancelOrder(symbol, takeProfitID); err != nil {
			What(log.Warn().Err(err), "failed to cancel stop_loss after take profit got executed")
		}
	})

	// Next, create two new orders simultaneously.
	go create(exchange, symbol, closeSide, "stop_market", stopLossPrice, quantityStr, stopLossID)
	create(exchange, symbol, closeSide, "limit", takeProfitPrice, quantityStr, takeProfitID)
}
