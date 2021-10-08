package commons

// +-------+
// | Types |
// +-------+

type AlgoContext struct {
	Result bool
	Bools  map[string]bool
	Floats map[string]float64
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
