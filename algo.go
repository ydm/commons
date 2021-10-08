package commons

type Algo interface {
	Run(s SingleState) bool
}
