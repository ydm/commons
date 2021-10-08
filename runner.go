package commons

func Run(keeper *StateKeeper, algo Algo) {
	go func() {
		Checker.Push()
		defer Checker.Pop()

		for ticker := range keeper.Channel {
			algo.Run(ticker)
		}
	}()
}
