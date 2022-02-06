package commons

import "github.com/rs/zerolog/log"

// +------------+
// | WarmUpAlgo |
// +------------+

type WarmUpAlgo struct {
	warmUp     int
	candlesKey string
}

func NewWarmUpAlgo(warmUp int, candlesKey string) WarmUpAlgo {
	return WarmUpAlgo{warmUp, candlesKey}
}

func (a WarmUpAlgo) Run(input AlgoContext, _ Ticker) AlgoContext {
	if input.CandlesLen(a.candlesKey) < a.warmUp {
		What(
			log.Debug().
				Int("candles", input.CandlesLen(a.candlesKey)).
				Int("warmup", a.warmUp),
			"still warming up",
		)

		return False
	}

	return input
}
