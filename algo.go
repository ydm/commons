package commons

type Algo interface {
	Run(input AlgoContext, ticker Ticker) (output AlgoContext)
}

type AlgoContainer interface {
	Algo
	Insert(a Algo)
}
