package commons

import (
	"fmt"
	"time"

	"github.com/alexander-yu/stream/minmax"
)

// +--------+
// | Candle |
// +--------+

// Candle types
//
// NB: These constants should be aligned with baser/constants.py.
const (
	DollarCandle = iota + 1
	TickCandle
	TimeCandle
	VolumeCandle
)

// Candle symbols
//
// NB: These constants should be aligned with baser/constants.py.
const (
	BTCUSDT = iota + 1
	ETHUSDT
	BCHUSDT
	LINKUSDT
	LTCUSDT
)

type Candle struct {
	Type                     int
	Symbol                   string
	Open                     float64
	High                     float64
	Low                      float64
	Close                    float64
	Volume                   float64
	VWAP                     float64
	QuoteAssetVolume         float64
	NumberOfTrades           int
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
	OpenTime                 time.Time
	CloseTime                time.Time
	LastTradeID              int64
}

func (c *Candle) TypeSymbol() string {
	prefix := "_"
	switch c.Type {
	case DollarCandle:
		prefix = "d"
	case TickCandle:
		prefix = "k"
	case VolumeCandle:
		prefix = "v"
	}

	return prefix + c.Symbol
}

// +---------------+
// | CandleBuilder |
// +---------------+

type CandleBuilder struct {
	Open                     float64
	High                     minmax.Max
	Low                      minmax.Min
	Close                    float64
	Volume                   float64
	QuoteAssetVolume         float64
	NumberOfTrades           int
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
	OpenTime                 time.Time
	CloseTime                time.Time
	LastTradeID              int64
}

func NewCandleBuilder() *CandleBuilder {
	builder := &CandleBuilder{
		Open:                     0,
		High:                     *minmax.NewGlobalMax(),
		Low:                      *minmax.NewGlobalMin(),
		Close:                    0,
		Volume:                   0,
		QuoteAssetVolume:         0,
		NumberOfTrades:           0,
		TakerBuyBaseAssetVolume:  0,
		TakerBuyQuoteAssetVolume: 0,
		OpenTime:                 time.Time{},
		CloseTime:                time.Time{},
		LastTradeID:              0,
	}
	builder.reset()

	return builder
}

func (b *CandleBuilder) Push(t Ticker) {
	if b.NumberOfTrades <= 0 {
		b.Open = t.TradePrice
		b.OpenTime = t.Time
	}

	if err := b.High.Push(t.TradePrice); err != nil {
		panic(err)
	}

	if err := b.Low.Push(t.TradePrice); err != nil {
		panic(err)
	}

	b.Close = t.TradePrice
	b.Volume += t.TradeQuantity
	quoteAssetVolume := t.TradePrice * t.TradeQuantity
	b.QuoteAssetVolume += quoteAssetVolume
	b.NumberOfTrades++
	b.CloseTime = t.Time
	b.LastTradeID = t.TradeID

	if !t.BuyerIsMaker {
		b.TakerBuyBaseAssetVolume += t.TradeQuantity
		b.TakerBuyQuoteAssetVolume += quoteAssetVolume
	}
}

func (b *CandleBuilder) Seal() (candle Candle, err error) {
	high, err := b.High.Value()
	if err != nil {
		return candle, fmt.Errorf("high.Value() failed: %w", err)
	}

	low, err := b.Low.Value()
	if err != nil {
		return candle, fmt.Errorf("low.Value() failed: %w", err)
	}

	// Populate return object.
	candle.Open = b.Open
	candle.High = high
	candle.Low = low
	candle.Close = b.Close
	candle.Volume = b.Volume

	candle.VWAP = 0
	if b.Volume != 0 {
		candle.VWAP = b.QuoteAssetVolume / b.Volume
	}

	candle.QuoteAssetVolume = b.QuoteAssetVolume
	candle.NumberOfTrades = b.NumberOfTrades
	candle.TakerBuyBaseAssetVolume = b.TakerBuyBaseAssetVolume
	candle.TakerBuyQuoteAssetVolume = b.TakerBuyQuoteAssetVolume
	candle.OpenTime = b.OpenTime.UTC()
	candle.CloseTime = b.CloseTime.UTC()
	candle.LastTradeID = b.LastTradeID

	// After all the data is snapped, reset object.
	b.reset()

	return candle, nil
}

func (b *CandleBuilder) reset() {
	b.Open = 0
	b.High.Clear()
	b.Low.Clear()
	b.Close = 0
	b.Volume = 0
	b.QuoteAssetVolume = 0
	b.NumberOfTrades = 0
	b.TakerBuyBaseAssetVolume = 0
	b.TakerBuyQuoteAssetVolume = 0
	b.OpenTime = time.Time{}
	b.CloseTime = time.Time{}
	b.LastTradeID = 0
}

// +-----------+
// | Utilities |
// +-----------+

type Extractor func(c Candle) float64

func MapCandles(f Extractor, candles []Candle) []float64 {
	n := len(candles)
	ans := make([]float64, n)

	for i := 0; i < n; i++ {
		ans[i] = f(candles[i])
	}

	return ans
}

func Opens(candles []Candle) []float64 {
	return MapCandles(func(c Candle) float64 {
		return c.Open
	}, candles)
}

func Highs(candles []Candle) []float64 {
	return MapCandles(func(c Candle) float64 {
		return c.High
	}, candles)
}

func Lows(candles []Candle) []float64 {
	return MapCandles(func(c Candle) float64 {
		return c.Low
	}, candles)
}

func Closes(candles []Candle) []float64 {
	return MapCandles(func(c Candle) float64 {
		return c.Close
	}, candles)
}

func Volumes(candles []Candle) []float64 {
	return MapCandles(func(c Candle) float64 {
		return c.Volume
	}, candles)
}

func VWAPs(candles []Candle) []float64 {
	return MapCandles(func(c Candle) float64 {
		return c.VWAP
	}, candles)
}
