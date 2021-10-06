package commons

type Algo interface {
	Run(s State) bool
}
