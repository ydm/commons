package commons

import (
	"time"

	"github.com/rs/zerolog/log"
)

// Pass a copy to make sure it's not changed.  MB this isn't smart.
// As long it's not a bottleneck, I don't give a fuck.  Plus, I
// wouldn't be surprised if the runtime actually optimizes this.
type Predicate func(*CandleBuilder) bool

type CandlesAlgo struct {
	Candles   CircularArray
	builder   CandleBuilder
	predicate Predicate
	last      time.Time
}

func NewCandlesAlgo(predicate Predicate) *CandlesAlgo {
	return &CandlesAlgo{
		Candles:   NewCircularArray(256),
		builder:   NewCandleBuilder(),
		predicate: predicate,
		last:      time.Now().UTC(),
	}
}

func (a *CandlesAlgo) Run(input AlgoContext, ticker Ticker) AlgoContext {
	if ticker.Last == TradeUpdate {
		a.builder.Push(ticker)

		// Indicates whether we should produce a new candle
		// and clear the builder.
		clear := a.predicate(&a.builder)

		// Produce and publish candle.
		if clear {
			candle := a.builder.Clear()

			// Adjust times.
			since := time.Since(a.last)
			a.last = time.Now().UTC()

			// Log candle.
			What(log.Debug().Interface("candle", candle).Dur("dur", since), "producing new candle")

			// Add to array.
			a.Candles.Push(candle)

			// Create output context and add candles.
			output := input.Copy()
			output.Objects["candles"] = &a.Candles

			return output
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
