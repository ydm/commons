package backtest

import (
	"strings"

	"github.com/ydm/commons"
)

func originMatches(origin, symbol string) bool {
	return strings.HasPrefix(origin, symbol)
}

func ReadTrades(archives []string, symbol string) <-chan commons.Trade {
	trades := make(chan commons.Trade)

	go func() {
		commons.Checker.Push()
		defer commons.Checker.Pop()

		defer close(trades)

		for _, archive := range archives {
			for row := range ReadCSVZipArchive(archive) {
				if originMatches(row.Origin, symbol) {
					trade := commons.TradeFromStrings(symbol, row.Values)
					trades <- trade
				}
			}
		}
	}()

	return trades
}
