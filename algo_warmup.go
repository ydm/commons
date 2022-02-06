package commons

import "github.com/rs/zerolog/log"

// +------------+
// | WarmUpAlgo |
// +------------+

type WarmUpAlgo struct {
	warmUp int
}

func NewWarmUpAlgo(warmUp int) WarmUpAlgo {
	return WarmUpAlgo{warmUp}
}

func (a WarmUpAlgo) Run(input AlgoContext, _ Ticker) AlgoContext {
	if input.CandlesLen() < a.warmUp {
		What(
			log.Debug().
				Int("candles", input.CandlesLen()).
				Int("warmup", a.warmUp),
			"still warming up",
		)

		return False
	}

	return input
}
