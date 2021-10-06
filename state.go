package commons

import (
	"sync"
	"time"
)

type State struct {
	Symbol string
	Now    time.Time

	LastPrice    float64
	LastQuantity float64

	BidPrice    float64
	BidQuantity float64
	AskPrice    float64
	AskQuantity float64
}

func (s *State) ApplyTrade(x Trade) {
	s.LastPrice = x.Price
	s.LastQuantity = x.Quantity
}

func (s *State) ApplyBook1(x Book1) {
	s.AskPrice = x.AskPrice
	s.AskQuantity = x.AskQuantity
	s.BidPrice = x.BidPrice
	s.BidQuantity = x.BidQuantity
}

type StateKeeper struct {
	s State
	m sync.Mutex
}

func NewStateKeeper() (k StateKeeper) {
	return StateKeeper{
		s: State{
			Symbol:       "",
			Now:          time.Time{},
			LastPrice:    0,
			LastQuantity: 0,
			BidPrice:     0,
			BidQuantity:  0,
			AskPrice:     0,
			AskQuantity:  0,
		},
		m: sync.Mutex{},
	}
}

func (k *StateKeeper) ApplyTrade(x Trade) {
	k.m.Lock()
	k.s.ApplyTrade(x)
	k.m.Unlock()
}

func (k *StateKeeper) ApplyBook1(x Book1) {
	k.m.Lock()
	k.s.ApplyBook1(x)
	k.m.Unlock()
}
