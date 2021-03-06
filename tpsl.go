package commons

import (
	"strings"

	"github.com/rs/zerolog/log"
)

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
	takeProfitID := RandomOrderID("P_")
	stopLossID := RandomOrderID("L_")

	handler := func(update OrderUpdate) {
		if strings.EqualFold(update.Status, "filled") {
			return
		}

		cancelID := ""

		switch update.ClientOrderID {
		case takeProfitID:
			cancelID = stopLossID
		case stopLossID:
			cancelID = takeProfitID
		}

		What(
			log.Info().
				Bool("stopLossFilled", update.ClientOrderID == stopLossID).
				Bool("takeProfitFilled", update.ClientOrderID == takeProfitID).
				Interface("update", update).
				Str("stopLossID", stopLossID).
				Str("takeProfitID", takeProfitID).
				Str("cancelID", cancelID),
			"one TPSL order filled, canceling the other",
		)

		if cancelID != "" {
			err := exchange.CancelOrder(symbol, cancelID)
			if err != nil {
				What(
					log.Warn().
						Err(err).
						Str("symbol", symbol).
						Str("stopLossID", stopLossID),
					"failed to cancel stop_loss after take_profit got filled",
				)
			}
		}
	}

	// First make sure we handle order filled events.
	exchange.OnOrderUpdate(takeProfitID, handler)
	exchange.OnOrderUpdate(stopLossID, handler)

	// Second, create two orders simultaneously.
	go create(exchange, symbol, closeSide, "stop_market", stopLossPrice, quantityStr, stopLossID)
	create(exchange, symbol, closeSide, "limit", takeProfitPrice, quantityStr, takeProfitID)
}
