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
	afterTradeID int64
	builder      *CandleBuilder
	candleType   int
	key          string
	predicate    Predicate
	symbol       string
}

func NewCandlesAlgo(afterTradeID int64, candleType int, key string, predicate Predicate, symbol string) *CandlesAlgo {
	key = DefaultString(key, "candles")

	return &CandlesAlgo{
		Candles:      NewCircularArray(256),
		afterTradeID: afterTradeID,
		builder:      NewCandleBuilder(),
		candleType:   candleType,
		key:          key,
		predicate:    predicate,
		symbol:       symbol,
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
		candle, err := a.builder.Seal()
		if err != nil {
			return False
		}

		// Assign extra attributes.
		candle.Type = a.candleType
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

	return NewCandlesAlgo(afterTradeID, TickCandle, key, predicate, symbol)
}

// +----------------+
// | Volume candles |
// +----------------+

func NewVolumeCandlesAlgo(volumeThreshold float64, key string, symbol string, afterTradeID int64) *CandlesAlgo {
	predicate := func(b *CandleBuilder) bool {
		return b.Volume >= volumeThreshold
	}

	return NewCandlesAlgo(afterTradeID, VolumeCandle, key, predicate, symbol)
}

// +-----------+
// | $ candles |
// +-----------+

func NewDollarCandlesAlgo(threshold float64, key string, symbol string, afterTradeID int64) *CandlesAlgo {
	predicate := func(b *CandleBuilder) bool {
		return b.QuoteAssetVolume >= threshold
	}

	return NewCandlesAlgo(afterTradeID, DollarCandle, key, predicate, symbol)
}

// +--------------+
// | Time candles |
// +--------------+

type TimeCandlesAlgo struct {
	Candles         CircularArray
	afterTradeID    int64
	builder         *CandleBuilder
	duration        time.Duration
	key             string
	lastAlignedTime time.Time
	lastTradeTime   time.Time
	symbol          string
}

func NewTimeCandlesAlgo(dur time.Duration, key string, symbol string, afterTradeID int64) *TimeCandlesAlgo {
	return &TimeCandlesAlgo{
		Candles:         NewCircularArray(256),
		afterTradeID:    afterTradeID,
		builder:         NewCandleBuilder(),
		duration:        dur,
		key:             key,
		lastAlignedTime: time.Time{},
		lastTradeTime:   time.Time{},
		symbol:          symbol,
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
		candle, err := a.builder.Seal()
		if err != nil {
			return False
		}

		// Assign symbol.
		candle.Type = TimeCandle
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
