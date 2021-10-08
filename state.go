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
	Symbol string

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
	if s.Symbol != x.Symbol {
		panic(x.Symbol)
	}

	s.TradePrice = x.Price
	s.TradeQuantity = x.Quantity

	return *s
}

func (s *Ticker) ApplyBook1(x Book1) Ticker {
	if s.Symbol != x.Symbol {
		panic(x.Symbol)
	}

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

func (s *State) ApplyTrade(x Trade) Ticker {
	s.Now = x.Time

	return s.Symbols[x.Symbol].ApplyTrade(x)
}

func (s *State) ApplyBook1(x Book1) Ticker {
	s.Now = x.Time

	return s.Symbols[x.Symbol].ApplyBook1(x)
}

// +-------------+
// | StateKeeper |
// +-------------+

type StateKeeper struct {
	Channel chan Ticker
	state   State
	locks   map[string]*sync.Mutex
}

func NewStateKeeper(numChannels int, symbols ...string) (k StateKeeper) {
	// We yield a new State through this channel after each
	// update.
	channel := make(chan Ticker, 16)

	// We need to know in advance how many input channels this
	//  keeper will consume.  Kind of stupid, but it is what it
	//  is.  When all of these input channels get closed, we close
	//  our channel too.
	var wg sync.WaitGroup

	wg.Add(numChannels)

	go func() {
		Checker.Push()
		defer Checker.Pop()

		wg.Wait()
		close(channel)
	}()

	// Initialize the state kept.  We need to know in advance how
	// many symbols we'll be managing the state for.
	state := State{
		Now:     time.Time{},
		Symbols: make(map[string]*Ticker),
	}

	// Each symbol has an associated lock (used for state updates)
	// and a ticker (where state is kept).
	locks := make(map[string]*sync.Mutex)

	for _, s := range symbols {
		state.Symbols[s] = &Ticker{
			Symbol:        s,
			TradePrice:    0,
			TradeQuantity: 0,
			BidPrice:      0,
			BidQuantity:   0,
			AskPrice:      0,
			AskQuantity:   0,
		}
		locks[s] = &sync.Mutex{}
	}

	return StateKeeper{
		Channel: channel,
		state:   state,
		locks:   locks,
	}
}

func (k *StateKeeper) ConsumeTrade(xs chan Trade) {
	go func() {
		Checker.Push()
		defer Checker.Pop()

		for x := range xs {
			m := k.locks[x.Symbol]
			m.Lock()
			ticker := k.state.ApplyTrade(x)
			m.Unlock()
			k.Channel <- ticker
		}
	}()
}

func (k *StateKeeper) ConsumeBook1(xs chan Book1) {
	go func() {
		Checker.Push()
		defer Checker.Pop()

		for x := range xs {
			m := k.locks[x.Symbol]
			m.Lock()
			ticker := k.state.ApplyBook1(x)
			m.Unlock()
			k.Channel <- ticker
		}
	}()
}
