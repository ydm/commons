package binancef

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/rs/zerolog/log"
	"github.com/ydm/commons"
)

type BinanceFutures struct {
	Client *futures.Client
}

func New(apikey, secret string) BinanceFutures {
	client := futures.NewClient(apikey, secret)

	return BinanceFutures{client}
}

// CreateOrder accepts arguments of the following formats:
//
// - side: buy or sell
// - orderType: market
func (b *BinanceFutures) CreateOrder(
	symbol string,
	side string,
	orderType string,
	quantityStr string,
	reduceOnly bool,
) (commons.CreateOrderResponse, error) {
	futuresSide := futures.SideType(
		switchSideString(
			side,
			string(futures.SideTypeBuy),
			string(futures.SideTypeSell),
		),
	)
	futuresOrderType := futures.OrderType(
		switchTypeString(orderType, string(futures.OrderTypeMarket)),
	)

	service := b.Client.NewCreateOrderService().
		Symbol(symbol).
		Side(futuresSide).
		PositionSide(futures.PositionSideTypeBoth).
		Type(futuresOrderType).
		// TimeInForce(futures.TimeInForceTypeFOK).
		Quantity(quantityStr).
		ReduceOnly(reduceOnly).
		// Price
		// ClientOrderID
		// StopPrice
		// WorkingType
		WorkingType(futures.WorkingTypeContractPrice).
		// ActivationPrice
		// CallbackRate (TODO)
		PriceProtect(true).
		NewOrderResponseType(futures.NewOrderRespTypeRESULT)
		// ClosePosition(false)

	res, err := service.Do(context.Background())
	if err != nil {
		commons.Msg(
			log.Fatal().
				Err(err).
				Str("symbol", symbol).
				Str("side", side).
				Str("quantity", quantityStr).
				Bool("reduceOnly", reduceOnly),
		)
	}

	return commons.CreateOrderResponse{
		ExecutedQuantity: res.ExecutedQuantity,
		AvgPrice:         res.AvgPrice,
	}, nil
}

func (b *BinanceFutures) Book1(ctx context.Context, symbol string) chan commons.Book1 {
	return SubscribeBookTicker(ctx, symbol)
}

func (b *BinanceFutures) Trade(ctx context.Context, symbol string) chan commons.Trade {
	return SubscribeAggTrade(ctx, symbol)
}
