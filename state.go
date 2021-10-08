package commons

import (
	"sync"
	"time"
)

// +--------+
// | Ticker |
// +--------+

// Ticker holds all the relevant data for a single trading pair.
type Ticker struct {
	// Trade data.
	TradePrice    float64
	TradeQuantity float64

	// Order book level 1 data.
	BidPrice    float64
	BidQuantity float64
	AskPrice    float64
	AskQuantity float64
}

func (s *Ticker) ApplyTrade(x Trade) Ticker {
	s.TradePrice = x.Price
	s.TradeQuantity = x.Quantity

	return *s
}

func (s *Ticker) ApplyBook1(x Book1) Ticker {
	s.AskPrice = x.AskPrice
	s.AskQuantity = x.AskQuantity
	s.BidPrice = x.BidPrice
	s.BidQuantity = x.BidQuantity

	return *s
}

// +-------+
// | State |
// +-------+

type State struct {
	Now     time.Time
	Symbols map[string]*Ticker
}

func (s State) ApplyTrade(symbol string, x Trade) Ticker {
	return s.Symbols[symbol].ApplyTrade(x)
}

func (s State) ApplyBook1(symbol string, x Book1) Ticker {
	return s.Symbols[symbol].ApplyBook1(x)
}

// +-------------+
// | StateKeeper |
// +-------------+

type StateKeeper struct {
	Channel chan Ticker
	state   State
}

func NewStateKeeper(numChannels int) (k StateKeeper) {
	// We yield a new State through this channel after each update.
	channel := make(chan Ticker)

	// We need to know in advance how many input channels this keeper will consume.
	var wg sync.WaitGroup
	wg.Add(numChannels)

	go func() {
		Checker.Push()
		defer Checker.Pop()

		wg.Wait()
		close(channel)
	}()

	return StateKeeper{
		Channel: channel,
		states:  sync.Map{},
	}
}

func (k *StateKeeper) ConsumeTrade(xs chan Trade) {
	go func() {
		Checker.Push()
		defer Checker.Pop()

		for x := range xs {

			k.m.Lock()
			k.s.ApplyTrade(x)
			state := k.s
			k.m.Unlock()

			k.C <- state
		}
	}()
}

func (k *StateKeeper) ConsumeBook1(xs chan Book1) {
	go func() {
		Checker.Push()
		defer Checker.Pop()

		for x := range xs {
			k.m.Lock()
			k.s.ApplyBook1(x)
			state := k.s
			k.m.Unlock()

			k.C <- state
		}
	}()
}
