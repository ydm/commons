package commons

type Algo interface {
	Run(t Ticker) bool
}
