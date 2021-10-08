package commons

// +---------------+
// | AlgoContainer |
// +---------------+

type AlgoContainer interface {
	Algo
	Insert(a Algo)
}

// +-------------+
// | SingleStack |
// +-------------+

type SingleStack struct {
	algos []Algo
	stack []AlgoContainer
}

func NewSingleStack() *SingleStack {
	return &SingleStack{
		algos: make([]Algo, 0, 8),
		stack: make([]AlgoContainer, 0, 2),
	}
}

// Run implements Algo.
func (s *SingleStack) Run(input AlgoContext, ticker Ticker) (output AlgoContext) {
	ctx := input

	for _, a := range s.algos {
		if !ctx.Result {
			return ctx
		}

		ctx = a.Run(ctx, ticker)
	}

	output = ctx

	return
}

// Insert implements AlgoContainer.
func (s *SingleStack) Insert(a Algo) {
	for i := len(s.stack) - 1; i >= 0; i-- {
		s.stack[i].Insert(a)
		a = s.stack[i]
	}

	s.algos = append(s.algos, a)
	s.stack = make([]AlgoContainer, 0, 2)
}
