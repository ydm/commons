package commons

import "sync"

type ParallelAlgo struct {
	algos []Algo
}

func NewParallelAlgo() *ParallelAlgo {
	return &ParallelAlgo{
		algos: make([]Algo, 0, 8),
	}
}

func (p *ParallelAlgo) Run(input AlgoContext, ticker Ticker) (output AlgoContext) {
	var wg sync.WaitGroup

	for i := range p.algos {
		wg.Add(1)

		go func(algo Algo, copied AlgoContext) {
			Checker.Push()
			defer Checker.Pop()

			defer wg.Done()

			algo.Run(copied, ticker)
		}(p.algos[i], input.Copy())
	}

	// Wait for all parallel algos of this iteration to finish.
	wg.Wait()

	// TODO: It could be useful to return a union of all returned contexts.
	return True
}

func (p *ParallelAlgo) Insert(a Algo) {
	p.algos = append(p.algos, a)
}
