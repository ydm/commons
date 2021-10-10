package commons

import "github.com/rs/zerolog"

// +-------+
// | Types |
// +-------+

type AlgoContext struct {
	Result bool
	Bools  map[string]bool
	Floats map[string]float64
}

func (c AlgoContext) Copy() AlgoContext {
	ans := AlgoContext{
		Result: c.Result,
		Bools:  make(map[string]bool),
		Floats: make(map[string]float64),
	}

	for k, v := range c.Bools {
		ans.Bools[k] = v
	}

	for k, v := range c.Floats {
		ans.Floats[k] = v
	}

	return ans
}

func (c AlgoContext) Dict(dict *zerolog.Event) *zerolog.Event {
	if dict == nil {
		dict = zerolog.Dict()
	}

	dict.Bool("__result__", c.Result)

	for k, v := range c.Bools {
		dict.Bool(k, v)
	}

	for k, v := range c.Floats {
		dict.Float64(k, v)
	}

	return dict
}

type Algo interface {
	Run(input AlgoContext, ticker Ticker) (output AlgoContext)
}

// +-----------+
// | Constants |
// +-----------+

//nolint:gochecknoglobals
var True = AlgoContext{
	Result: true,
	Bools:  nil,
	Floats: nil,
}

//nolint:gochecknoglobals
var False = AlgoContext{
	Result: false,
	Bools:  nil,
	Floats: nil,
}
