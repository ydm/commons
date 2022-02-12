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
	key       string
}

func NewCandlesAlgo(predicate Predicate, key string) *CandlesAlgo {
	key = DefaultString(key, "candles")

	return &CandlesAlgo{
		Candles:   NewCircularArray(256),
		builder:   NewCandleBuilder(),
		predicate: predicate,
		last:      time.Now().UTC(),
		key:       key,
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
			candle, err := a.builder.Clear()
			if err != nil {
				return False
			}

			// Adjust times.
			since := time.Since(a.last)
			a.last = time.Now().UTC()

			// Log candle.
			Code(log.Trace().Interface("candle", candle).Dur("dur", since), "candle")

			// Add to array.
			a.Candles.Push(candle)

			// Create output context and add candles.
			output := input.Copy()
			output.Objects[a.key] = &a.Candles

			return output
		}
	}

	return False
}

// +--------------+
// | Tick candles |
// +--------------+

func NewTickCandlesAlgo(numTicks int, key string) *CandlesAlgo {
	predicate := func(b *CandleBuilder) bool {
		return b.NumberOfTrades >= numTicks
	}

	return NewCandlesAlgo(predicate, key)
}

// +----------------+
// | Volume candles |
// +----------------+

func NewVolumeCandlesAlgo(volumeThreshold float64, key string) *CandlesAlgo {
	predicate := func(b *CandleBuilder) bool {
		return b.Volume >= volumeThreshold
	}

	return NewCandlesAlgo(predicate, key)
}

// +-----------+
// | $ candles |
// +-----------+

func NewDollarCandlesAlgo(threshold float64, key string) *CandlesAlgo {
	predicate := func(b *CandleBuilder) bool {
		return b.QuoteAssetVolume >= threshold
	}

	return NewCandlesAlgo(predicate, key)
}
