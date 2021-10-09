package commons

import (
	"github.com/rs/zerolog/log"
)

type VerboseAlgo struct{}

func (v VerboseAlgo) Run(ctx AlgoContext, ticker Ticker) AlgoContext {
	Msg(log.Info().Interface("ticker", ticker))

	return ctx
}
