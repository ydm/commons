package commons

import (
	"sync"
	"time"
)

// +-------+
// | State |
// +-------+

type State struct {
	Symbol string
	Now    time.Time

	// Trade data.
	TradePrice    float64
	TradeQuantity float64

	// Order book level 1 data.
	BidPrice    float64
	BidQuantity float64
	AskPrice    float64
	AskQuantity float64
}

func (s *State) ApplyTrade(x Trade) {
	if s.Symbol != x.Symbol {
		panic("")
	}

	s.Now = x.Time
	s.TradePrice = x.Price
	s.TradeQuantity = x.Quantity
}

func (s *State) ApplyBook1(x Book1) {
	if s.Symbol != x.Symbol {
		panic("")
	}

	s.Now = x.Time
	s.AskPrice = x.AskPrice
	s.AskQuantity = x.AskQuantity
	s.BidPrice = x.BidPrice
	s.BidQuantity = x.BidQuantity
}

// +-------------+
// | StateKeeper |
// +-------------+

type StateKeeper struct {
	C chan State
	s State
	m sync.Mutex
}

func NewStateKeeper() (k StateKeeper) {
	return StateKeeper{
		C: make(chan State),
		s: State{
			Symbol:        "",
			Now:           time.Time{},
			TradePrice:    0,
			TradeQuantity: 0,
			BidPrice:      0,
			BidQuantity:   0,
			AskPrice:      0,
			AskQuantity:   0,
		},
		m: sync.Mutex{},
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
