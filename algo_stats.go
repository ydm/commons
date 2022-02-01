package commons

import (
	"fmt"

	"github.com/markcheno/go-talib"
)

type StatsRSI struct {
	InTimePeriod int
}

func (s StatsRSI) Run(input AlgoContext, ticker Ticker) (output AlgoContext) {
	candles, ok := input.Objects["candles"].(*CircularArray)
	if !ok {
		panic("not candles")
	}

	if candles.Len() < (s.InTimePeriod + 1) {
		fmt.Printf("NOT ENOUGH CANDLES YET: %d\n", candles.Len())

		return False
	}

	in := VWAPs(candles, s.InTimePeriod)
	for i := range in {
		fmt.Printf("in[%d]=%f\n", i, in[i])
	}

	out := talib.Rsi(in, s.InTimePeriod)
	// for i := range out {
	// 	fmt.Printf("out[%d]=%f\n", i, out[i])
	// }
	ans := input.Copy()
	ans.Floats["rsi"] = out[len(out)-1]

	return ans
}
