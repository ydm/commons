package commons

import (
	"time"
)

// Pass a copy to make sure it's not changed.  MB this isn't smart.
// As long it's not a bottleneck, I don't give a fuck.  Plus, I
// wouldn't be surprised if the runtime actually optimizes this.
type Predicate func(*CandleBuilder) bool

type CandlesAlgo struct {
	Candles      CircularArray
	builder      *CandleBuilder
	predicate    Predicate
	key          string
	symbol       string
	afterTradeID int64
}

func NewCandlesAlgo(predicate Predicate, key string, symbol string, afterTradeID int64) *CandlesAlgo {
	key = DefaultString(key, "candles")

	return &CandlesAlgo{
		Candles:      NewCircularArray(256),
		builder:      NewCandleBuilder(),
		predicate:    predicate,
		key:          key,
		symbol:       symbol,
		afterTradeID: afterTradeID,
	}
}

func (a *CandlesAlgo) Run(input AlgoContext, ticker Ticker) AlgoContext {
	// If AfterTradeID is provided, make sure it is respected.
	if a.afterTradeID > 0 && a.afterTradeID >= ticker.TradeID {
		return False
	}

	// Skip all tickers that are not concerning our symbol.
	if a.symbol != "" && ticker.Symbol != a.symbol {
		return False
	}

	// Skip all tickers that are not updated after a trade.
	if ticker.Last != TradeUpdate {
		return False
	}

	a.builder.Push(ticker)

	// Indicates whether we should produce a new candle
	// and clear the builder.
	clear := a.predicate(a.builder)

	// Produce and publish candle.
	if clear {
		candle, err := a.builder.Clear()
		if err != nil {
			return False
		}

		// Assign symbol.
		candle.Symbol = ticker.Symbol

		// Add to array.
		a.Candles.Push(candle)

		// Create output context and add candles.
		output := input.Copy()
		output.Objects[a.key] = &a.Candles

		return output
	}

	return False
}

// +--------------+
// | Tick candles |
// +--------------+

func NewTickCandlesAlgo(numTicks int, key string, symbol string, afterTradeID int64) *CandlesAlgo {
	predicate := func(b *CandleBuilder) bool {
		return b.NumberOfTrades >= numTicks
	}

	return NewCandlesAlgo(predicate, key, symbol, afterTradeID)
}

// +----------------+
// | Volume candles |
// +----------------+

func NewVolumeCandlesAlgo(volumeThreshold float64, key string, symbol string, afterTradeID int64) *CandlesAlgo {
	predicate := func(b *CandleBuilder) bool {
		return b.Volume >= volumeThreshold
	}

	return NewCandlesAlgo(predicate, key, symbol, afterTradeID)
}

// +-----------+
// | $ candles |
// +-----------+

func NewDollarCandlesAlgo(threshold float64, key string, symbol string, afterTradeID int64) *CandlesAlgo {
	predicate := func(b *CandleBuilder) bool {
		return b.QuoteAssetVolume >= threshold
	}

	return NewCandlesAlgo(predicate, key, symbol, afterTradeID)
}

// +--------------+
// | Time candles |
// +--------------+

type TimeCandlesAlgo struct {
	Candles         CircularArray
	builder         *CandleBuilder
	duration        time.Duration
	lastAlignedTime time.Time
	lastTradeTime   time.Time
	key             string
	symbol          string
	afterTradeID    int64
}

func NewTimeCandlesAlgo(dur time.Duration, key string, symbol string, afterTradeID int64) *TimeCandlesAlgo {
	return &TimeCandlesAlgo{
		Candles:         NewCircularArray(256),
		builder:         NewCandleBuilder(),
		duration:        dur,
		lastAlignedTime: time.Time{},
		lastTradeTime:   time.Time{},
		key:             key,
		symbol:          symbol,
		afterTradeID:    afterTradeID,
	}
}

func (a *TimeCandlesAlgo) predicate(tradeTime time.Time) (ans bool) {
	if a.lastTradeTime.After(tradeTime) {
		panic("unexpected")
	}

	aligned := AlignTime(tradeTime, a.duration)

	if a.lastAlignedTime.IsZero() {
		a.lastAlignedTime = aligned
	}

	if !a.lastAlignedTime.Equal(aligned) {
		a.lastAlignedTime = aligned
		ans = true
	}

	return ans
}

func (a *TimeCandlesAlgo) Run(input AlgoContext, ticker Ticker) AlgoContext {
	// If AfterTradeID is provided, make sure it is respected.
	if a.afterTradeID > 0 && a.afterTradeID >= ticker.TradeID {
		return False
	}

	// Skip all tickers that do not concern our symbol.
	if a.symbol != "" && ticker.Symbol != a.symbol {
		return False
	}

	// Skip all tickers that are not updated after a trade.
	if ticker.Last != TradeUpdate {
		return False
	}

	// Indicates whether we should produce a new candle
	// and clear the builder.
	clear := a.predicate(ticker.Time)

	// Produce and publish candle.
	if clear {
		candle, err := a.builder.Clear()
		if err != nil {
			return False
		}

		// Assign symbol.
		candle.Symbol = ticker.Symbol

		// Add to array.
		a.Candles.Push(candle)

		// Create output context and add candles.
		output := input.Copy()
		output.Objects[a.key] = &a.Candles

		// Before returning, push the latest ticker.
		// Otherwise it would be just skipped.
		a.builder.Push(ticker)

		return output
	}

	// We don't have to build a candle.  Just push latest ticker.
	a.builder.Push(ticker)

	return False
}
