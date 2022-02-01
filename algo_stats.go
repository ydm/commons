package commons

import (
	"github.com/markcheno/go-talib"
	"github.com/rs/zerolog/log"
)

func defaultString(x, defaultValue string) string {
	if x != "" {
		return x
	}

	return defaultValue
}

// +-----+
// | RSI |
// +-----+

type StatsAlgoRSI struct {
	InTimePeriod int
	Key          string
}

func (s StatsAlgoRSI) Run(input AlgoContext, ticker Ticker) AlgoContext {
	candles, err := input.Candles(s.InTimePeriod + 1)
	if err != nil {
		Msg(log.Debug().Err(err))

		return False
	}

	vwaps := VWAPs(candles)
	ans := talib.Rsi(vwaps, s.InTimePeriod)
	key := defaultString(s.Key, "rsi")
	last := len(ans) - 1
	output := input.Copy()
	output.Floats[key] = ans[last]

	return output
}

// +-----+
// | BOP |
// +-----+

type StatsAlgoBOP struct {
	Key string
}

func (s StatsAlgoBOP) Run(input AlgoContext, ticker Ticker) AlgoContext {
	candles, err := input.Candles(1)
	if err != nil {
		Msg(log.Debug().Err(err))

		return False
	}

	opens := Opens(candles)
	highs := Highs(candles)
	lows := Lows(candles)
	closes := Closes(candles)
	ans := talib.Bop(opens, highs, lows, closes)
	key := defaultString(s.Key, "bop")
	last := len(ans) - 1
	output := input.Copy()
	output.Floats[key] = ans[last]

	return output
}

// +------------+
// | Volume EMA |
// +------------+

type StatsAlgoVolumeEMA struct {
	InTimePeriod int
	Key          string
}

func (s StatsAlgoVolumeEMA) Run(input AlgoContext, ticker Ticker) AlgoContext {
	candles, err := input.Candles(s.InTimePeriod)
	if err != nil {
		Msg(log.Debug().Err(err))

		return False
	}

	volumes := Volumes(candles)
	ans := talib.Ma(volumes, s.InTimePeriod, talib.SMA)
	key := defaultString(s.Key, "volume_ema")
	last := len(ans) - 1
	output := input.Copy()
	output.Floats[key] = ans[last]

	// PrintArrayF64(volumes, 2)
	// PrintArrayF64(ans, 2)
	// fmt.Printf("v=[%f %f] x=%f t=%t\n", volumes[0], volumes[1], ans[1], (volumes[0]+volumes[1])/2 == ans[1])

	return output
}
