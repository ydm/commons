package binancef

import (
	"context"
	"errors"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/rs/zerolog/log"
	"github.com/ydm/commons"
)

var ErrLeverageNotSet = errors.New("leverage not set")

type BinanceFutures struct {
	client *futures.Client
}

func New(apikey, secret string) commons.Exchange {
	client := futures.NewClient(apikey, secret)

	return &BinanceFutures{client}
}

// +----------+
// | REST API |
// +----------+

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

	service := b.client.NewCreateOrderService().
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

func (b *BinanceFutures) ChangeMarginType(symbol, marginType string) error {
	m := futures.MarginType(switchMarginTypeString(
		marginType,
		string(futures.MarginTypeCrossed),
		string(futures.MarginTypeIsolated),
	))

	err := b.client.NewChangeMarginTypeService().
		Symbol(symbol).
		MarginType(m).
		Do(context.Background())

	if err != nil && err.Error() == "<APIError> code=-4046, msg=No need to change margin type." {
		return nil
	}

	return err
}

func (b *BinanceFutures) ChangeLeverage(symbol string, leverage int) error {
	resp, err := b.client.NewChangeLeverageService().
		Symbol(symbol).
		Leverage(leverage).
		Do(context.Background())

	if err != nil {
		return err
	}

	if resp.Leverage != leverage {
		return ErrLeverageNotSet
	}

	return nil
}

// +-----------+
// | Websocket |
// +-----------+

func (b *BinanceFutures) Book1(ctx context.Context, symbol string) chan commons.Book1 {
	return SubscribeBookTicker(ctx, symbol)
}

func (b *BinanceFutures) Trade(ctx context.Context, symbol string) chan commons.Trade {
	return SubscribeAggTrade(ctx, symbol)
}
