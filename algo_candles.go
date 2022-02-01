package commons

// Pass a copy to make sure it's not changed.  MB this isn't smart.
// As long it's not a bottleneck, I don't give a fuck.  Plus, I
// wouldn't be surprised if the runtime actually optimizes this.
type Predicate func(*CandleBuilder) bool

type CandlesAlgo struct {
	Candles   CircularArray
	builder   CandleBuilder
	predicate Predicate
}

func NewCandlesAlgo(predicate Predicate) *CandlesAlgo {
	return &CandlesAlgo{
		Candles:   NewCircularArray(256),
		builder:   NewCandleBuilder(),
		predicate: predicate,
	}
}

func (a *CandlesAlgo) Run(ctx AlgoContext, ticker Ticker) AlgoContext {
	if ticker.Last == TradeUpdate {
		a.builder.Push(ticker)

		// Indicates whether we should produce a new candle
		// and clear the builder.
		clear := a.predicate(&a.builder)

		// Produce and publish candle.
		if clear {
			candle := a.builder.Clear()
			a.Candles.Push(candle)

			ans := ctx.Copy()
			ans.Objects["candles"] = &a.Candles
			return ans
		}
	}

	return False
}

// +--------------+
// | Tick candles |
// +--------------+

func NewTickCandlesAlgo(numTicks int) *CandlesAlgo {
	predicate := func(b *CandleBuilder) bool {
		return b.NumberOfTrades >= numTicks
	}

	return NewCandlesAlgo(predicate)
}

// +----------------+
// | Volume candles |
// +----------------+

func NewVolumeCandlesAlgo(volumeThreshold float64) *CandlesAlgo {
	predicate := func(b *CandleBuilder) bool {
		return b.Volume >= volumeThreshold
	}

	return NewCandlesAlgo(predicate)
}

// +-----------+
// | $ candles |
// +-----------+

func NewDollarCandlesAlgo(threshold float64) *CandlesAlgo {
	predicate := func(b *CandleBuilder) bool {
		return b.QuoteAssetVolume >= threshold
	}

	return NewCandlesAlgo(predicate)
}
