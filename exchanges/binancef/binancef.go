package binancef

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/rs/zerolog/log"
	"github.com/ydm/commons"
	"github.com/ydm/commons/exchanges/bb"
)

// +----------------+
// | BinanceFutures |
// +----------------+

var ErrLeverageNotSet = errors.New("leverage not set")

type BinanceFutures struct {
	client               *futures.Client
	streamer             *bb.Streamer
	orderUpdateCallbacks commons.CircularArray
	orderUpdateMutex     sync.Mutex
}

func New(ctx context.Context, apikey, secret string) *BinanceFutures {
	client := futures.NewClient(apikey, secret)

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
		streamer = bb.NewStreamer(bb.NewFuturesStreamService(client))
	}

	ans := &BinanceFutures{
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

func binanceSide(side string) futures.SideType {
	return futures.SideType(
		bb.SwitchSideString(
			side,
			string(futures.SideTypeBuy),
			string(futures.SideTypeSell),
		),
	)
}

func binanceOrderType(orderType string) futures.OrderType {
	return futures.OrderType(
		bb.SwitchTypeString(
			orderType,
			string(futures.OrderTypeMarket),
			string(futures.OrderTypeLimit),
			string(futures.OrderTypeStopMarket),
		),
	)
}

//nolint:funlen
func (b *BinanceFutures) CreateOrder(
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
		PositionSide(futures.PositionSideTypeBoth).
		Type(futuresOrderType).
		// TimeInForce(futures.TimeInForceTypeFOK).
		Quantity(quantityStr).
		ReduceOnly(reduceOnly).
		// Price
		// NewClientOrderID(clientOrderID).
		// StopPrice
		// WorkingType
		WorkingType(futures.WorkingTypeMarkPrice).
		// ActivationPrice
		// CallbackRate (TODO)
		PriceProtect(true).
		NewOrderResponseType(futures.NewOrderRespTypeRESULT)
	//     ClosePosition(false)

	if clientOrderID != "" {
		service = service.NewClientOrderID(clientOrderID)
	}

	switch orderType {
	case bb.OrderTypeLimit:
		service = service.
			TimeInForce(futures.TimeInForceTypeGTC).
			Price(priceStr)
	case bb.OrderTypeStopMarket:
		service = service.
			TimeInForce(futures.TimeInForceTypeGTC).
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
		AvgPrice:         res.AvgPrice,
	}, nil
}

func (b *BinanceFutures) CancelOrder(symbol string, clientOrderID string) error {
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

func (b *BinanceFutures) ChangeMarginType(symbol, marginType string) error {
	m := futures.MarginType(bb.SwitchMarginTypeString(
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

	return fmt.Errorf("ChangeMarginType: symbol=%s, marginType=%s, err=%w",
		symbol, marginType, err)
}

func (b *BinanceFutures) ChangeLeverage(symbol string, leverage int) error {
	resp, err := b.client.NewChangeLeverageService().
		Symbol(symbol).
		Leverage(leverage).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("ChangeLeverage: symbol=%s, leverage=%d, err=%w",
			symbol, leverage, err)
	}

	if resp.Leverage != leverage {
		return ErrLeverageNotSet
	}

	return nil
}

// +-----------+
// | Websocket |
// +-----------+

func (b *BinanceFutures) Book1(ctx context.Context, symbol string) <-chan commons.Book1 {
	return SubscribeBookTicker(ctx, symbol)
}

func (b *BinanceFutures) Trade(ctx context.Context, symbol string) <-chan commons.Trade {
	return SubscribeAggTrade(ctx, symbol)
}

// +---------+
// | Events  |
// +---------+

type orderUpdateNode struct {
	clientOrderID string
	callback      commons.OrderUpdateCallback
}

func (b *BinanceFutures) OnOrderUpdate(clientOrderID string, callback commons.OrderUpdateCallback) {
	node := orderUpdateNode{clientOrderID, callback}

	// Make sure nobody else modifies orderUpdateCallbacks.
	b.orderUpdateMutex.Lock()
	defer b.orderUpdateMutex.Unlock()

	// Push to orderUpdateCallbacks.
	b.orderUpdateCallbacks.Push(node)
}

func (b *BinanceFutures) handleEvents() {
	for pointer := range b.streamer.Events {
		event, ok := pointer.(*futures.WsUserDataEvent)
		if !ok {
			panic("TODO")
		}

		if event.Event == futures.UserDataEventTypeOrderTradeUpdate {
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

				if node.clientOrderID == event.OrderTradeUpdate.ClientOrderID {
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
					ClientOrderID: event.OrderTradeUpdate.ClientOrderID,
					Status:        string(event.OrderTradeUpdate.Status),
				}
				node.callback(update)
			}
		}
	}
}
