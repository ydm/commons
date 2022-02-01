package commons

import "github.com/rs/zerolog"

// +-------------+
// | AlgoContext |
// +-------------+

type AlgoContext struct {
	Result  bool
	Bools   map[string]bool
	Floats  map[string]float64
	Ints    map[string]int
	Objects map[string]interface{}
}

func NewAlgoContext() AlgoContext {
	return AlgoContext{
		Result:  true,
		Bools:   make(map[string]bool),
		Floats:  make(map[string]float64),
		Ints:    make(map[string]int),
		Objects: make(map[string]interface{}),
	}
}

func (c AlgoContext) Copy() AlgoContext {
	ans := AlgoContext{
		Result:  c.Result,
		Bools:   make(map[string]bool),
		Floats:  make(map[string]float64),
		Ints:    make(map[string]int),
		Objects: make(map[string]interface{}),
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

	// TODO: This shit is not a deep copy and that should be
	// documented.
	for k, v := range c.Objects {
		ans.Objects[k] = v
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

	for k, v := range c.Objects {
		dict.Interface(k, v)
	}

	return dict
}

// +-----------+
// | Constants |
// +-----------+

//nolint:gochecknoglobals
var True = AlgoContext{
	Result:  true,
	Bools:   nil,
	Floats:  nil,
	Ints:    nil,
	Objects: nil,
}

//nolint:gochecknoglobals
var False = AlgoContext{
	Result:  false,
	Bools:   nil,
	Floats:  nil,
	Ints:    nil,
	Objects: nil,
}
