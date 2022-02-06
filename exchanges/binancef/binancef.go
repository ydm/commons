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
)

// // +----------+
// // | Streamer |
// // +----------+

// type Streamer struct {
// 	ctx     context.Context
// 	service StreamService
// 	Events  chan interface{}
// }

// func NewStreamer(ctx context.Context, service StreamService) *Streamer {
// 	streamer := &Streamer{
// 		ctx:     ctx,
// 		service: service,
// 		Events:  make(chan interface{}),
// 	}
// 	return streamer
// }

// func (s *Streamer) Loop() {
// 	go func() {
// 		commons.Checker.Push()
// 		defer commons.Checker.Pop()

// 		if err := s.loop(); err != nil {
// 			commons.Msg(log.Error().Err(err))
// 		}
// 	}()
// }

// func (s *Streamer) loop() (err error) {
// 	defer close(s.Events)

// 	var previousListenKey string
// 	for s.ctx.Err() == nil {
// 		// For the Do() method I'm not using ctx, because in case of a closed
// 		// context, it panics...  TODO: Commit a patch!
// 		listenKey, err := s.service.Start(context.Background())
// 		if err != nil {
// 			return err
// 		}
// 		// This is an ugly workaround for a bug (in Binance's API) I'm too lazy to
// 		// debug right now.  Basically the listenKey returned is the same.  As of
// 		// 2021-02-17 many Binance Futures bugs I encountered in the past are no
// 		// longer present, but this fix should stay just in case.
// 		if listenKey == previousListenKey {
// 			continue
// 		}
// 		previousListenKey = listenKey

// 		log.Info().
// 			Str("what", "starting user stream").
// 			Str("listenKey", listenKey).
// 			Msg(fortuna.Location())
// 		done, stop, err := s.service.Feed(s.ctx, listenKey, s.Events)
// 		if err != nil {
// 			log.Error().Err(err).Msg(fortuna.Location())
// 			time.Sleep(15 * time.Second)
// 			continue
// 		}
// 		go func() {
// 			fortuna.CheckerPush()
// 			defer fortuna.CheckerPop()

// 			s.closeWhenDone(done, stop, listenKey)
// 		}()
// 		s.keepalive(done, listenKey)
// 	}
// 	return
// }

// func (s *Streamer) closeWhenDone(done, stop chan struct{}, listenKey string) {
// 	select {
// 	case <-s.ctx.Done():
// 		close(stop)
// 	case <-done:
// 	}

// 	log.Info().
// 		Str("what", "closing user stream").
// 		Str("listenKey", listenKey).
// 		Msg(fortuna.Location())
// 	err := s.service.Close(context.Background(), listenKey)
// 	if err != nil {
// 		log.Error().Err(err).Msg(fortuna.Location())
// 	}
// }

// func (s *Streamer) keepalive(done <-chan struct{}, listenKey string) {
// 	ticker := time.NewTicker(20 * time.Minute)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			err := s.service.Keepalive(context.Background(), listenKey)
// 			if err != nil {
// 				log.Warn().Err(err).Str("listenKey", listenKey).Msg(fortuna.Location())
// 			} else {
// 				log.Info().Str("listenKey", listenKey).Msg(fortuna.Location())
// 			}
// 		case <-s.ctx.Done():
// 			return
// 		case <-done:
// 			return
// 		}
// 	}
// }

// +----------------+
// | BinanceFutures |
// +----------------+

var ErrLeverageNotSet = errors.New("leverage not set")

type BinanceFutures struct {
	client               *futures.Client
	streamer             Streamer
	orderUpdateCallbacks commons.CircularArray
	orderUpdateMutex     sync.Mutex
}

func New(apikey, secret string) *BinanceFutures {
	client := futures.NewClient(apikey, secret)

	// Solves some time window problems.
	if _, err := client.NewSetServerTimeService().Do(context.Background()); err != nil {
		panic(err)
	}

	ans := &BinanceFutures{
		client:               client,
		streamer:             *NewStreamer(NewFuturesStreamService(client)),
		orderUpdateCallbacks: commons.NewCircularArray(256),
		orderUpdateMutex:     sync.Mutex{},
	}

	ans.streamer.Loop(context.Background())

	// TODO: Move that out!
	go func() {
		commons.Checker.Push()
		defer commons.Checker.Pop()

		for pointer := range ans.streamer.Events {
			event, ok := pointer.(*futures.WsUserDataEvent)
			if !ok {
				panic("TODO")
			}

			if event.Event == futures.UserDataEventTypeOrderTradeUpdate {
				var (
					node  orderUpdateNode
					found = false
				)

				ans.orderUpdateMutex.Lock()
				n := ans.orderUpdateCallbacks.Len()
				for i := 0; i < n; i++ {
					item := ans.orderUpdateCallbacks.At(i)
					node, ok = item.(orderUpdateNode)
					if !ok {
						panic("TODO")
					}

					if node.clientOrderID == event.OrderTradeUpdate.ClientOrderID {
						found = true
						break
					}
				}
				ans.orderUpdateMutex.Unlock()

				if found {
					update := commons.OrderUpdate{
						ClientOrderID: event.OrderTradeUpdate.ClientOrderID,
						Status:        string(event.OrderTradeUpdate.Status),
					}
					node.callback(update)
				}
			}
		}
	}()

	return ans
}

// +----------+
// | REST API |
// +----------+

func binanceSide(side string) futures.SideType {
	return futures.SideType(
		switchSideString(
			side,
			string(futures.SideTypeBuy),
			string(futures.SideTypeSell),
		),
	)
}

func binanceOrderType(orderType string) futures.OrderType {
	return futures.OrderType(
		switchTypeString(
			orderType,
			string(futures.OrderTypeMarket),
			string(futures.OrderTypeLimit),
			string(futures.OrderTypeStopMarket),
		),
	)
}

func (b *BinanceFutures) CreateOrder(
	symbol string,
	side string,
	orderType string,
	priceStr string,
	quantityStr string,
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
		// ClientOrderID
		// StopPrice
		// WorkingType
		WorkingType(futures.WorkingTypeContractPrice).
		// ActivationPrice
		// CallbackRate (TODO)
		PriceProtect(true).
		NewOrderResponseType(futures.NewOrderRespTypeRESULT)
	//     ClosePosition(false)

	switch orderType {
	case orderTypeLimit:
		service = service.Price(priceStr)
	case orderTypeStopMarket:
		service = service.StopPrice(priceStr)
	}

	if orderType == orderTypeLimit || orderType == orderTypeStopMarket {
		service = service.TimeInForce(futures.TimeInForceTypeGTC)
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

func (b *BinanceFutures) Book1(ctx context.Context, symbol string) chan commons.Book1 {
	return SubscribeBookTicker(ctx, symbol)
}

func (b *BinanceFutures) Trade(ctx context.Context, symbol string) chan commons.Trade {
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
	b.orderUpdateMutex.Lock()
	b.orderUpdateCallbacks.Push(node)
	b.orderUpdateMutex.Unlock()
}
