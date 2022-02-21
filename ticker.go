package commons

import (
	"time"
)

const (
	TradeUpdate = iota + 1
	Book1Update
)

// Ticker holds all the relevant data for a single trading pair.
type Ticker struct {
	Time   time.Time
	Symbol string

	// Trade data.
	TradeID       int64
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

	s.Time = x.Time
	s.TradeID = x.TradeID
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

	s.Time = x.Time
	s.AskPrice = x.AskPrice
	s.AskQuantity = x.AskQuantity
	s.BidPrice = x.BidPrice
	s.BidQuantity = x.BidQuantity
	s.Last = Book1Update

	return *s
}
