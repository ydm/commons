package commons

import "github.com/rs/zerolog"

// +-------------+
// | AlgoContext |
// +-------------+

type AlgoContext struct {
	Result bool
	Bools  map[string]bool
	Floats map[string]float64
	Ints   map[string]int
}

func NewAlgoContext() AlgoContext {
	return AlgoContext{
		Result: true,
		Bools:  make(map[string]bool),
		Floats: make(map[string]float64),
		Ints:   make(map[string]int),
	}
}

func (c AlgoContext) Copy() AlgoContext {
	ans := AlgoContext{
		Result: c.Result,
		Bools:  make(map[string]bool),
		Floats: make(map[string]float64),
		Ints:   make(map[string]int),
	}

	for k, v := range c.Bools {
		ans.Bools[k] = v
	}

	for k, v := range c.Floats {
		ans.Floats[k] = v
	}

	for k, v := range c.Ints {
		ans.Ints[k] = v
	}

	return ans
}

func (c AlgoContext) Dict(dict *zerolog.Event) *zerolog.Event {
	if dict == nil {
		dict = zerolog.Dict()
	}

	// dict.Bool("__result__", c.Result)

	for k, v := range c.Bools {
		dict.Bool(k, v)
	}

	for k, v := range c.Floats {
		dict.Float64(k, v)
	}

	for k, v := range c.Ints {
		dict.Int(k, v)
	}

	return dict
}

// +-----------+
// | Constants |
// +-----------+

//nolint:gochecknoglobals
var True = AlgoContext{
	Result: true,
	Bools:  nil,
	Floats: nil,
	Ints:   nil,
}

//nolint:gochecknoglobals
var False = AlgoContext{
	Result: false,
	Bools:  nil,
	Floats: nil,
	Ints:   nil,
}
