package binances

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/adshao/go-binance/v2"
	"github.com/rs/zerolog/log"
	"github.com/ydm/commons"
	"github.com/ydm/commons/exchanges/bb"
)

// +----------------+
// | BinanceFutures |
// +----------------+

var (
	ErrLeverageNotSet = errors.New("leverage not set")
	ErrNotImplemented = errors.New("not implemented")
)

type BinanceSpot struct {
	client               *binance.Client
	streamer             *bb.Streamer
	orderUpdateCallbacks commons.CircularArray
	orderUpdateMutex     sync.Mutex
}

func New(ctx context.Context, apikey, secret string) *BinanceSpot {
	client := binance.NewClient(apikey, secret)

	// Synchronize with server by adjusting an internal time offset.
	service := client.NewSetServerTimeService()
	timeOffset, err := service.Do(context.Background()) //nolint:contextcheck

	if err != nil {
		panic(err)
	} else {
		commons.What(log.Info().Int64("timeOffset", timeOffset), "initialized server time offset")
	}

	var streamer *bb.Streamer
	if apikey != "" && secret != "" {
		streamer = bb.NewStreamer(bb.NewSpotStreamService(client))
	}

	ans := &BinanceSpot{
		client:               client,
		streamer:             streamer,
		orderUpdateCallbacks: commons.NewCircularArray(256),
		orderUpdateMutex:     sync.Mutex{},
	}

	if ans.streamer != nil {
		// Handle streamed events in a separate goroutine.
		go func() {
			commons.Checker.Push()
			defer commons.Checker.Pop()

			ans.handleEvents()
		}()

		// Start event loop.  It runs until context is cancelled.
		ans.streamer.Loop(ctx)
	}

	return ans
}

// +----------+
// | REST API |
// +----------+

func binanceSide(side string) binance.SideType {
	return binance.SideType(
		bb.SwitchSideString(
			side,
			string(binance.SideTypeBuy),
			string(binance.SideTypeSell),
		),
	)
}

func binanceOrderType(orderType string) binance.OrderType {
	return binance.OrderType(
		bb.SwitchTypeString(
			orderType,
			string(binance.OrderTypeMarket),
			string(binance.OrderTypeLimit),
			string(binance.OrderTypeStopLoss), // TODO: Is this stop market?
		),
	)
}

//nolint:funlen
func (b *BinanceSpot) CreateOrder(
	symbol string,
	side string,
	orderType string,
	priceStr string,
	quantityStr string,
	clientOrderID string,
	reduceOnly bool,
) (commons.CreateOrderResponse, error) {
	futuresSide := binanceSide(side)
	futuresOrderType := binanceOrderType(orderType)

	service := b.client.NewCreateOrderService().
		Symbol(symbol).
		Side(futuresSide).
		Type(futuresOrderType).
		// TimeInForce(binance.TimeInForceTypeFOK).
		Quantity(quantityStr).
		// QuoteOrderQty
		// Price
		// NewClientOrderID(clientOrderID).
		// StopPrice(stopPrice)
		// IcebergQuantity(icebergQuantity)
		NewOrderRespType(binance.NewOrderRespTypeRESULT)

	if clientOrderID != "" {
		service = service.NewClientOrderID(clientOrderID)
	}

	switch orderType {
	case bb.OrderTypeLimit:
		service = service.
			TimeInForce(binance.TimeInForceTypeGTC).
			Price(priceStr)
	case bb.OrderTypeStopMarket:
		service = service.
			TimeInForce(binance.TimeInForceTypeGTC).
			StopPrice(priceStr)
	}

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
		OrderID:          strconv.FormatInt(res.OrderID, 10),
		ClientOrderID:    res.ClientOrderID,
		ExecutedQuantity: res.ExecutedQuantity,
		AvgPrice:         res.Price, // TODO: Check if this field contains the average price.
	}, nil
}

func (b *BinanceSpot) CancelOrder(symbol string, clientOrderID string) error {
	service := b.client.NewCancelOrderService().
		Symbol(symbol).
		OrigClientOrderID(clientOrderID)

	resp, err := service.Do(context.Background())
	if err != nil && err.Error() != "<APIError> code=-2011, msg=Unknown order sent." {
		return bb.Wrap(err, "cancel order failed")
	}

	if resp.ClientOrderID != clientOrderID {
		commons.Msg(
			log.Fatal().
				Str("clientOrderID", clientOrderID).
				Str("resp.ClientOrderID", resp.ClientOrderID),
		)
	}

	return nil
}

func (b *BinanceSpot) ChangeMarginType(symbol, marginType string) error {
	return ErrNotImplemented
}

func (b *BinanceSpot) ChangeLeverage(symbol string, leverage int) error {
	return ErrNotImplemented
}

// +-----------+
// | Websocket |
// +-----------+

func (b *BinanceSpot) Book1(ctx context.Context, symbol string) <-chan commons.Book1 {
	return SubscribeBookTicker(ctx, symbol)
}

func (b *BinanceSpot) Trade(ctx context.Context, symbol string) <-chan commons.Trade {
	return SubscribeAggTrade(ctx, symbol)
}

// +---------+
// | Events  |
// +---------+

type orderUpdateNode struct {
	clientOrderID string
	callback      commons.OrderUpdateCallback
}

func (b *BinanceSpot) OnOrderUpdate(clientOrderID string, callback commons.OrderUpdateCallback) {
	node := orderUpdateNode{clientOrderID, callback}

	// Make sure nobody else modifies orderUpdateCallbacks.
	b.orderUpdateMutex.Lock()
	defer b.orderUpdateMutex.Unlock()

	// Push to orderUpdateCallbacks.
	b.orderUpdateCallbacks.Push(node)
}

func (b *BinanceSpot) handleEvents() {
	for pointer := range b.streamer.Events {
		event, ok := pointer.(*binance.WsUserDataEvent)
		if !ok {
			panic("TODO")
		}

		if event.Event == binance.UserDataEventTypeExecutionReport {
			var (
				node  orderUpdateNode
				found = false
			)

			// Beginning of locked section.
			b.orderUpdateMutex.Lock()

			// Iterate over all callbacks.
			n := b.orderUpdateCallbacks.Len()
			for i := 0; i < n; i++ {
				item := b.orderUpdateCallbacks.At(i)

				node, ok = item.(orderUpdateNode)
				if !ok {
					panic("TODO")
				}

				if node.clientOrderID == event.OrderUpdate.ClientOrderId {
					// We found the proper callback for this order.
					// It has to be executed outside of the locked
					// section, so we just mark it as found.
					found = true

					break
				}
			}

			b.orderUpdateMutex.Unlock()

			if found {
				update := commons.OrderUpdate{
					ClientOrderID: event.OrderUpdate.ClientOrderId,
					Status:        string(event.OrderUpdate.Status),
				}
				node.callback(update)
			}
		}
	}
}
