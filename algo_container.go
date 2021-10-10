package commons

// +---------------+
// | AlgoContainer |
// +---------------+

type AlgoContainer interface {
	Algo
	Insert(a Algo)
}
