package backtest

import (
	"github.com/ydm/commons"
)

// func SymbolFromFilename(filename string) string {
// 	split := strings.SplitN(filename, "-", 2)
// 	if len(split) != 2 {
// 		panic("")
// 	}
//
// 	return split[0]
// }

func ReadTrades(archives []string, symbol string) <-chan commons.Trade {
	trades := make(chan commons.Trade)

	go func() {
		commons.Checker.Push()
		defer commons.Checker.Pop()

		defer close(trades)

		for _, archive := range archives {
			for row := range ReadCSVZipArchive(archive) {
				// symbol := SymbolFromFilename(row.Origin)
				trade := commons.TradeFromStrings(symbol, row.Values)
				trades <- trade
			}
		}
	}()

	return trades
}
