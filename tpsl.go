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

		What(log.Info().Interface("update", update)

		if err := exchange.CancelOrder(symbol, stopLossID); err != nil {
			What(log.Warn().Err(err), "failed to cancel stop_loss after take profit got executed")
		}
	})

	exchange.OnOrderUpdate(stopLossID, func(update OrderUpdate) {
		if update.Status != "FILLED" {
			return
		}

		if err := exchange.CancelOrder(symbol, takeProfitID); err != nil {
			What(log.Warn().Err(err), "failed to cancel stop_loss after take profit got executed")
		}
	})

	// Next, create two new orders simultaneously.

	go func() {
		_, err := exchange.CreateOrder(
			symbol,
			closeSide,
			"stop_market",
			stopLossPrice,
			quantityStr,
			stopLossID,
			true,
		)

		if err != nil {
			What(log.Warn().Err(err), "create order failed")
		}
	}()

	_, err := exchange.CreateOrder(
		symbol,
		closeSide,
		"limit",
		takeProfitPrice,
		quantityStr,
		takeProfitID,
		true,
	)

	if err != nil {
		What(log.Warn().Err(err), "create order failed")
	}
}
