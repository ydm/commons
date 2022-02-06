package commons

import (
	"errors"

	"github.com/rs/zerolog"
)

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

	dict.Bool("__result__", c.Result)

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

// +------+
// | Misc |
// +------+

var (
	ErrKeyNotFound      = errors.New("key not found")
	ErrNotEnoughCandles = errors.New("not enough candles")
	ErrNotCandle        = errors.New("object is not a candle")
)

func (c AlgoContext) Float(key string) float64 {
	x, ok := c.Floats[key]
	if !ok {
		panic("need float " + key)
	}

	return x
}

func (c AlgoContext) Bool(key string) bool {
	x, ok := c.Bools[key]
	if !ok {
		panic("need bool " + key)
	}

	return x
}

func (c AlgoContext) CandlesLen(key string) int {
	key = DefaultString(key, "candles")

	candles, ok := c.Objects[key].(*CircularArray)
	if !ok {
		return -1
	}

	return candles.Len()
}

func (c AlgoContext) Candles(key string, n int) ([]Candle, error) {
	key = DefaultString(key, "candles")

	candles, ok := c.Objects[key].(*CircularArray)
	if !ok {
		return nil, ErrKeyNotFound
	}

	length := candles.Len()
	if length < n {
		return nil, ErrNotEnoughCandles
	}

	ans := make([]Candle, n)

	for i := 0; i < n; i++ {
		candle, ok := candles.At(length - n + i).(Candle)
		if !ok {
			return ans, ErrNotCandle
		}

		ans[i] = candle
	}

	return ans, nil
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
