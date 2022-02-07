package commons

import (
	"fmt"

	"github.com/markcheno/go-talib"
	"github.com/rs/zerolog/log"
)

// +-----+
// | RSI |
// +-----+

type StatsAlgoRSI struct {
	InTimePeriod int
	Key          string
	CandlesKey   string
}

func (s StatsAlgoRSI) Run(input AlgoContext, ticker Ticker) AlgoContext {
	candles, err := input.Candles(s.CandlesKey, s.InTimePeriod+1)
	if err != nil {
		Msg(log.Debug().Err(err).Int("candles", input.CandlesLen(s.CandlesKey)))

		return False
	}

	vwaps := VWAPs(candles)
	rsi := talib.Rsi(vwaps, s.InTimePeriod)

	key := DefaultString(s.Key, "rsi")
	last := len(rsi) - 1

	output := input.Copy()
	output.Floats[key] = rsi[last]

	return output
}

// +----+
// | AD |
// +----+

type StatsAlgoAD struct {
	FastPeriod int
	SlowPeriod int
	Key        string
	CandlesKey string
}

func (s StatsAlgoAD) Run(input AlgoContext, ticker Ticker) AlgoContext {
	candles, err := input.Candles(s.CandlesKey, s.SlowPeriod+1)
	if err != nil {
		Msg(
			log.Debug().
				Err(err).
				Int("want", s.SlowPeriod+1).
				Int("have", input.CandlesLen(s.CandlesKey)),
		)

		return False
	}

	highs := Highs(candles)
	lows := Lows(candles)
	closes := Closes(candles)
	volumes := Volumes(candles)
	ans := talib.AdOsc(highs, lows, closes, volumes, s.FastPeriod, s.SlowPeriod)

	key := DefaultString(s.Key, fmt.Sprintf("ad_%d_%d", s.FastPeriod, s.SlowPeriod))
	last := len(ans) - 1

	output := input.Copy()
	output.Floats[key] = ans[last]

	return output
}

// +-----+
// | BOP |
// +-----+

type StatsAlgoBOP struct {
	Key        string
	CandlesKey string
}

func (s StatsAlgoBOP) Run(input AlgoContext, ticker Ticker) AlgoContext {
	candles, err := input.Candles(s.CandlesKey, 1)
	if err != nil {
		Msg(log.Debug().Err(err).Int("candles", input.CandlesLen(s.CandlesKey)))

		return False
	}

	opens := Opens(candles)
	highs := Highs(candles)
	lows := Lows(candles)
	closes := Closes(candles)
	bop := talib.Bop(opens, highs, lows, closes)

	key := DefaultString(s.Key, "bop")
	last := len(bop) - 1

	output := input.Copy()
	output.Floats[key] = bop[last]

	return output
}

// +-----+
// | CMO |
// +-----+

type StatsAlgoCMO struct {
	InTimePeriod int
	Key          string
	CandlesKey   string
}

func (s StatsAlgoCMO) Run(input AlgoContext, ticker Ticker) AlgoContext {
	candles, err := input.Candles(s.CandlesKey, s.InTimePeriod)
	if err != nil {
		Msg(
			log.Debug().
				Err(err).
				Int("want", s.InTimePeriod).
				Int("have", input.CandlesLen(s.CandlesKey)),
		)

		return False
	}

	closes := Closes(candles)
	ans := talib.Cmo(closes, s.InTimePeriod)

	key := DefaultString(s.Key, fmt.Sprintf("cmo%d", s.InTimePeriod))
	last := len(ans) - 1

	output := input.Copy()
	output.Floats[key] = ans[last]

	return output
}

// +-----------+
// | Price EMA |
// +-----------+

type StatsAlgoPriceEMA struct {
	InTimePeriod int
	Key          string
	CandlesKey   string
}

func (s StatsAlgoPriceEMA) Run(input AlgoContext, ticker Ticker) AlgoContext {
	candles, err := input.Candles(s.CandlesKey, s.InTimePeriod)
	if err != nil {
		Msg(
			log.Debug().
				Err(err).
				Int("want", s.InTimePeriod).
				Int("have", input.CandlesLen(s.CandlesKey)),
		)

		return False
	}

	closes := Closes(candles)
	ans := talib.Ma(closes, s.InTimePeriod, talib.EMA)

	key := DefaultString(s.Key, fmt.Sprintf("ema%d", s.InTimePeriod))
	last := len(ans) - 1

	output := input.Copy()
	output.Floats[key] = ans[last]

	return output
}

// +----------+
// | Price MA |
// +----------+

type StatsAlgoPriceMA struct {
	InTimePeriod int
	Key          string
	CandlesKey   string
}

func (s StatsAlgoPriceMA) Run(input AlgoContext, ticker Ticker) AlgoContext {
	candles, err := input.Candles(s.CandlesKey, s.InTimePeriod)
	if err != nil {
		Msg(
			log.Debug().
				Err(err).
				Int("want", s.InTimePeriod).
				Int("have", input.CandlesLen(s.CandlesKey)),
		)

		return False
	}

	closes := Closes(candles)
	ans := talib.Ma(closes, s.InTimePeriod, talib.SMA)

	key := DefaultString(s.Key, fmt.Sprintf("ma%d", s.InTimePeriod))
	last := len(ans) - 1

	output := input.Copy()
	output.Floats[key] = ans[last]

	return output
}

// +------------+
// | Volume EMA |
// +------------+

type StatsAlgoVolumeMA struct {
	InTimePeriod int
	Key          string
	CandlesKey   string
}

func (s StatsAlgoVolumeMA) Run(input AlgoContext, ticker Ticker) AlgoContext {
	candles, err := input.Candles(s.CandlesKey, s.InTimePeriod)
	if err != nil {
		Msg(
			log.Debug().
				Err(err).
				Int("want", s.InTimePeriod).
				Int("have", input.CandlesLen(s.CandlesKey)),
		)

		return False
	}

	volumes := Volumes(candles)
	ans := talib.Ma(volumes, s.InTimePeriod, talib.SMA)
	key := DefaultString(s.Key, "volume_ema")
	last := len(ans) - 1
	output := input.Copy()
	output.Floats[key] = ans[last]

	return output
}
