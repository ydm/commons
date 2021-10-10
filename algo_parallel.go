package commons

type Parallel struct {
	algos []Algo
}

func (p *Parallel) Run(input AlgoContext, ticker Ticker) (output AlgoContext) {
	for i := range p.algos {
		go func(algo Algo, copied AlgoContext) {
			Checker.Push()
			defer Checker.Pop()

			algo.Run(copied, ticker)
		}(p.algos[i], input.Copy())
	}

	return True
}

func (p *Parallel) Insert(a Algo) {
	p.algos = append(p.algos, a)
}
