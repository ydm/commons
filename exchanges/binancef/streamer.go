package binancef

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/rs/zerolog/log"
	"github.com/ydm/commons"
)

func wrap(err error, message string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", message, err)
	}

	return nil
}

// +-------------------+
// | BaseStreamService |
// +-------------------+

type BaseStreamService interface {
	Start(ctx context.Context) (listenKey string, err error)
	Close(ctx context.Context, listenKey string) error
	Keepalive(ctx context.Context, listenKey string) error
}

// +--------------------------+
// | BaseFuturesStreamService |
// +--------------------------+

type BaseFuturesStreamService struct {
	Client *futures.Client
}

func (s BaseFuturesStreamService) Start(ctx context.Context) (listenKey string, err error) {
	listenKey, err = s.Client.NewStartUserStreamService().Do(ctx)

	return
}

func (s BaseFuturesStreamService) Close(ctx context.Context, listenKey string) error {
	err := s.Client.NewCloseUserStreamService().ListenKey(listenKey).Do(ctx)

	return wrap(err, "close failed")
}

func (s BaseFuturesStreamService) Keepalive(ctx context.Context, listenKey string) error {
	err := s.Client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(ctx)

	return wrap(err, "keepalive failed")
}

// +-----------------------+
// | BaseSpotStreamService |
// +-----------------------+

type BaseSpotStreamService struct {
	client *binance.Client
}

func (s BaseSpotStreamService) Start(ctx context.Context) (listenKey string, err error) {
	listenKey, err = s.client.NewStartUserStreamService().Do(ctx)

	return listenKey, wrap(err, "start stream failed")
}

func (s BaseSpotStreamService) Close(ctx context.Context, listenKey string) error {
	err := s.client.NewCloseUserStreamService().ListenKey(listenKey).Do(ctx)

	return wrap(err, "close stream failed")
}

func (s BaseSpotStreamService) Keepalive(ctx context.Context, listenKey string) error {
	err := s.client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(ctx)

	return wrap(err, "keepalive failed")
}

// +---------------+
// | StreamService |
// +---------------+

type StreamService interface {
	BaseStreamService
	Feed(ctx context.Context, listenKey string, events chan interface{}) (doneC, stopC chan struct{}, err error)
}

// +----------------------+
// | FuturesStreamService |
// +----------------------+

type FuturesStreamService struct {
	base BaseFuturesStreamService
}

func NewFuturesStreamService(client *futures.Client) FuturesStreamService {
	return FuturesStreamService{BaseFuturesStreamService{client}}
}

func (s FuturesStreamService) Start(ctx context.Context) (listenKey string, err error) {
	listenKey, err = s.base.Start(ctx)

	return listenKey, wrap(err, "start failed")
}

func (s FuturesStreamService) Close(ctx context.Context, listenKey string) error {
	err := s.base.Close(ctx, listenKey)

	return wrap(err, "close failed")
}

func (s FuturesStreamService) Keepalive(ctx context.Context, listenKey string) error {
	err := s.base.Keepalive(ctx, listenKey)

	return wrap(err, "keepalive failed")
}

func (s FuturesStreamService) Feed(ctx context.Context, listenKey string, events chan interface{}) (
	doneC,
	stopC chan struct{},
	err error,
) {
	doneC, stopC, err = futures.WsUserDataServe(
		listenKey,
		func(event *futures.WsUserDataEvent) {
			events <- event
		},
		func(inner error) {
			commons.Msg(log.Error().Err(inner))
		},
	)

	err = wrap(err, "serve failed")

	return
}

// +-------------------+
// | SpotStreamService |
// +-------------------+

type SpotStreamService struct {
	base BaseSpotStreamService
}

func NewSpotStreamService(client *binance.Client) SpotStreamService {
	return SpotStreamService{BaseSpotStreamService{client}}
}

func (s SpotStreamService) Start(ctx context.Context) (listenKey string, err error) {
	return s.base.Start(ctx)
}

func (s SpotStreamService) Close(ctx context.Context, listenKey string) error {
	return s.base.Close(ctx, listenKey)
}

func (s SpotStreamService) Keepalive(ctx context.Context, listenKey string) error {
	return s.base.Keepalive(ctx, listenKey)
}

func (s SpotStreamService) Feed(ctx context.Context, listenKey string, events chan interface{}) (
	doneC,
	stopC chan struct{},
	err error,
) {
	doneC, stopC, err = binance.WsUserDataServe(listenKey,
		func(event *binance.WsUserDataEvent) {
			events <- event
		},
		func(inner error) {
			commons.Msg(log.Error().Err(inner))
		},
	)

	err = wrap(err, "serve failed")

	return doneC, stopC, err
}

// +----------+
// | Streamer |
// +----------+

type Streamer struct {
	service StreamService
	Events  chan interface{}
}

func NewStreamer(service StreamService) *Streamer {
	streamer := &Streamer{
		service: service,
		Events:  make(chan interface{}),
	}

	return streamer
}

func (s *Streamer) Loop(ctx context.Context) {
	go func() {
		commons.Checker.Push()
		defer commons.Checker.Pop()

		if err := s.loop(ctx); err != nil {
			commons.Msg(log.Error().Err(err))
		}
	}()
}

func (s *Streamer) loop(ctx context.Context) (err error) {
	defer close(s.Events)

	var previousListenKey string

	for ctx.Err() == nil {
		// For the Do() method I'm not using ctx, because in case of a closed
		// context, it panics...  TODO: Commit a patch!
		listenKey, err := s.service.Start(ctx)
		if err != nil {
			return fmt.Errorf("service start failed: %w", err)
		}

		// This is an ugly workaround for a bug (in Binance's API) I'm too lazy to
		// debug right now.  Basically the listenKey returned is the same.  As of
		// 2021-02-17 many Binance Futures bugs I encountered in the past are no
		// longer present, but this fix should stay just in case.
		if listenKey == previousListenKey {
			continue
		}

		previousListenKey = listenKey
		commons.What(
			log.Info().Str("listenKey", listenKey),
			"starting user stream",
		)

		done, stop, err := s.service.Feed(ctx, listenKey, s.Events)
		if err != nil {
			commons.Msg(log.Error().Err(err))
			time.Sleep(15 * time.Second)

			continue
		}

		go func() {
			commons.Checker.Push()
			defer commons.Checker.Pop()

			s.closeWhenDone(ctx, done, stop, listenKey)
		}()

		s.keepalive(ctx, done, listenKey)
	}

	return err
}

func (s *Streamer) closeWhenDone(ctx context.Context, done, stop chan struct{}, listenKey string) {
	select {
	case <-ctx.Done():
		close(stop)
	case <-done:
	}

	commons.What(log.Info().Str("listenKey", listenKey), "closing user stream")

	if err := s.service.Close(ctx, listenKey); err != nil {
		commons.Msg(log.Error().Err(err))
	}
}

func (s *Streamer) keepalive(ctx context.Context, done <-chan struct{}, listenKey string) {
	ticker := time.NewTicker(20 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.service.Keepalive(ctx, listenKey)
			if err != nil {
				commons.Msg(log.Warn().Err(err).Str("listenKey", listenKey))
			} else {
				commons.Msg(log.Info().Str("listenKey", listenKey))
			}
		case <-ctx.Done():
			return
		case <-done:
			return
		}
	}
}
