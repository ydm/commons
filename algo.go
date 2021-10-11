package commons

type Algo interface {
	Run(input AlgoContext, ticker Ticker) (output AlgoContext)
}
