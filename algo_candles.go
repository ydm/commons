package commons

import "sync"

type TickCandlesAlgo struct {
	Candles  CircularArray
	NumTicks int
	builder  CandleBuilder
	i        int
	mutex    sync.Mutex
}

func NewTickCandlesAlgo(numTicks int) *TickCandlesAlgo {
	return &TickCandlesAlgo{
		Candles:  NewCircularArray(256),
		NumTicks: numTicks,
		builder:  NewCandleBuilder(),
		i:        0,
		mutex:    sync.Mutex{},
	}
}

func (a *TickCandlesAlgo) Run(ctx AlgoContext, ticker Ticker) AlgoContext {
	if ticker.Last == TradeUpdate {
		a.builder.Push(ticker)
		if a.builder.NumberOfTrades >= a.NumTicks {
			candle := a.builder.Clear()
			a.Candles.Push(candle)
			ctx.Result = true
		} else {
			ctx.Result = false
		}
	}

	return ctx
}
