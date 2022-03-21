package commons

import (
	"sync"
	"time"
)

// +-------+
// | State |
// +-------+

type State struct {
	Now     time.Time
	Tickers map[string]*Ticker
	lock    sync.Mutex
}

func (s *State) ApplyTrade(x Trade) Ticker {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Now = x.Time

	return s.Tickers[x.Symbol].ApplyTrade(x)
}

func (s *State) ApplyBook1(x Book1) Ticker {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Now = x.Time

	return s.Tickers[x.Symbol].ApplyBook1(x)
}

// +-------------+
// | StateKeeper |
// +-------------+

type StateKeeper struct {
	// We yield a new Ticker, representing the current state,
	// through this channel after each update.
	Channel   chan Ticker
	state     State
	waitGroup sync.WaitGroup
}

func NewStateKeeper(numChannels int, symbols ...string) *StateKeeper {
	k := &StateKeeper{
		Channel: make(chan Ticker, 16),

		state: State{
			Now:     time.Time{},
			Tickers: make(map[string]*Ticker),
			lock:    sync.Mutex{},
		},
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
			Time:          time.Time{},
			TradeID:       0,
			TradePrice:    0,
			TradeQuantity: 0,
			BuyerIsMaker:  false,
			BidPrice:      0,
			BidQuantity:   0,
			AskPrice:      0,
			AskQuantity:   0,
			Last:          0,
		}
	}

	return k
}

// ConsumeTrade is usually called in a goroutine.
func (k *StateKeeper) ConsumeTrade(xs <-chan Trade) {
	defer k.waitGroup.Done()

	for x := range xs {
		ticker := k.state.ApplyTrade(x)
		k.Channel <- ticker
	}
}

// ConsumeBook1 is usually called in a goroutine.
func (k *StateKeeper) ConsumeBook1(xs <-chan Book1) {
	defer k.waitGroup.Done()

	for x := range xs {
		ticker := k.state.ApplyBook1(x)
		k.Channel <- ticker
	}
}
