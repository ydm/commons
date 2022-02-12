package commons

import (
	"fmt"
	"time"

	"github.com/alexander-yu/stream/minmax"
)

// +--------+
// | Candle |
// +--------+

type Candle struct {
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
	CloseTime                time.Time
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
}

func NewCandleBuilder() CandleBuilder {
	return CandleBuilder{
		Open:                     0,
		High:                     *minmax.NewGlobalMax(),
		Low:                      *minmax.NewGlobalMin(),
		Close:                    0,
		Volume:                   0,
		QuoteAssetVolume:         0,
		NumberOfTrades:           0,
		TakerBuyBaseAssetVolume:  0,
		TakerBuyQuoteAssetVolume: 0,
	}
}

func (b *CandleBuilder) Push(t Ticker) {
	if b.NumberOfTrades <= 0 {
		b.Open = t.TradePrice
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

	if !t.BuyerIsMaker {
		b.TakerBuyBaseAssetVolume += t.TradeQuantity
		b.TakerBuyQuoteAssetVolume += quoteAssetVolume
	}
}

func (b *CandleBuilder) Clear() (candle Candle, err error) {
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
	candle.CloseTime = time.Now().UTC()

	b.Open = 0
	b.High.Clear()
	b.Low.Clear()
	b.Close = 0
	b.Volume = 0
	b.QuoteAssetVolume = 0
	b.NumberOfTrades = 0
	b.TakerBuyBaseAssetVolume = 0
	b.TakerBuyQuoteAssetVolume = 0

	return candle, nil
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
