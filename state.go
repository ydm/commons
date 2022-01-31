package commons

import (
	"sync"
	"time"
)

// +--------+
// | Ticker |
// +--------+

const (
	TradeUpdate = iota + 1
	Book1Update
)

// Ticker holds all the relevant data for a single trading pair.
type Ticker struct {
	Now    time.Time
	Symbol string

	// Trade data.
	TradePrice    float64
	TradeQuantity float64
	BuyerIsMaker  bool

	// Order book level 1 data.
	BidPrice    float64
	BidQuantity float64
	AskPrice    float64
	AskQuantity float64

	// Type of the last update unit.
	Last int
}

func (s *Ticker) ApplyTrade(x Trade) Ticker {
	if s.Symbol != x.Symbol {
		panic(x.Symbol)
	}

	s.Now = x.Time
	s.TradePrice = x.Price
	s.TradeQuantity = x.Quantity
	s.BuyerIsMaker = x.BuyerIsMaker
	s.Last = TradeUpdate

	return *s
}

func (s *Ticker) ApplyBook1(x Book1) Ticker {
	if s.Symbol != x.Symbol {
		panic(x.Symbol)
	}

	s.Now = x.Time
	s.AskPrice = x.AskPrice
	s.AskQuantity = x.AskQuantity
	s.BidPrice = x.BidPrice
	s.BidQuantity = x.BidQuantity
	s.Last = Book1Update

	return *s
}

// +-------+
// | State |
// +-------+

type State struct {
	Now     time.Time
	Tickers map[string]*Ticker
}

func (s *State) ApplyTrade(x Trade) Ticker {
	s.Now = x.Time

	return s.Tickers[x.Symbol].ApplyTrade(x)
}

func (s *State) ApplyBook1(x Book1) Ticker {
	s.Now = x.Time

	return s.Tickers[x.Symbol].ApplyBook1(x)
}

// +-------------+
// | StateKeeper |
// +-------------+

type StateKeeper struct {
	// We yield a new State through this channel after each
	// update.
	Channel   chan Ticker
	state     State
	locks     map[string]*sync.Mutex
	waitGroup sync.WaitGroup
}

func NewStateKeeper(numChannels int, symbols ...string) *StateKeeper {
	k := &StateKeeper{
		Channel: make(chan Ticker, 16),

		state: State{
			Now:     time.Time{},
			Tickers: make(map[string]*Ticker),
		},
		locks:     make(map[string]*sync.Mutex),
		waitGroup: sync.WaitGroup{},
	}

	// We need to know in advance how many input channels this keeper will consume.
	// Kind of stupid, but it is what it is.  When all of these input channels get
	// closed, we close StateKeeper.Channel too.
	k.waitGroup.Add(numChannels)

	go func() {
		Checker.Push()
		defer Checker.Pop()

		k.waitGroup.Wait()
		close(k.Channel)
	}()

	// Initialize the state kept.  We need to know in advance how many symbols we'll
	// be managing.  Also, each symbol has an associated lock (used for state updates)
	// and a Ticker (where individual symbol states are kept).
	for _, s := range symbols {
		k.state.Tickers[s] = &Ticker{
			Symbol:        s,
			Now:           time.Time{},
			TradePrice:    0,
			TradeQuantity: 0,
			BuyerIsMaker:  false,
			BidPrice:      0,
			BidQuantity:   0,
			AskPrice:      0,
			AskQuantity:   0,
			Last:          0,
		}
		k.locks[s] = &sync.Mutex{}
	}

	return k
}

func (k *StateKeeper) ConsumeTrade(xs chan Trade) {
	go func() {
		Checker.Push()
		defer Checker.Pop()

		defer k.waitGroup.Done()

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

		defer k.waitGroup.Done()

		for x := range xs {
			m := k.locks[x.Symbol]
			m.Lock()
			ticker := k.state.ApplyBook1(x)
			m.Unlock()
			k.Channel <- ticker
		}
	}()
}
