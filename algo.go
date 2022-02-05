package commons

type Algo interface {
	Run(input AlgoContext, ticker Ticker) (output AlgoContext)
}

// type AlgoDependencies interface {
// 	Dependencies() []Algo
// }

type AlgoContainer interface {
	Algo
	Insert(a Algo)
}
