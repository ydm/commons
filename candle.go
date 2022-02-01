package commons

import (
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
	if b.NumberOfTrades == 0 {
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

func (b *CandleBuilder) Clear() Candle {
	high, err := b.High.Value()
	if err != nil {
		panic(err)
	}

	low, err := b.Low.Value()
	if err != nil {
		panic(err)
	}

	candle := Candle{
		Open:                     b.Open,
		High:                     high,
		Low:                      low,
		Close:                    b.Close,
		Volume:                   b.Volume,
		VWAP:                     b.QuoteAssetVolume / b.Volume,
		QuoteAssetVolume:         b.QuoteAssetVolume,
		NumberOfTrades:           b.NumberOfTrades,
		TakerBuyBaseAssetVolume:  b.TakerBuyBaseAssetVolume,
		TakerBuyQuoteAssetVolume: b.TakerBuyQuoteAssetVolume,
	}

	b.Open = 0
	b.High.Clear()
	b.Low.Clear()
	b.Close = 0
	b.Volume = 0
	b.QuoteAssetVolume = 0
	b.NumberOfTrades = 0
	b.TakerBuyBaseAssetVolume = 0
	b.TakerBuyQuoteAssetVolume = 0

	return candle
}

// +-----------+
// | Utilities |
// +-----------+

type Extractor func(c Candle) float64

func MapCandles(f Extractor, a *CircularArray, m int) []float64 {
	n := a.Len()
	ans := make([]float64, m)

	for i := n - m; i < n; i++ {
		candle, ok := a.At(i).(Candle)
		if !ok {
			panic("not a candle")
		}

		ans[i-n+m] = f(candle)
	}

	return ans
}

func Closes(a *CircularArray, m int) []float64 {
	return MapCandles(func(c Candle) float64 {
		return c.Close
	}, a, m)
}

func VWAPs(a *CircularArray, m int) []float64 {
	return MapCandles(func(c Candle) float64 {
		return c.VWAP
	}, a, m)
}
