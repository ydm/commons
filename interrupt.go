package commons

import (
	"context"
	"os"
	"os/signal"
)

func GoInterrupt(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		if Interrupt(ctx) {
			cancel()
		}
	}()
}

// Interrupt returns when either (1) interrupt signal is received by
// the OS or (2) the given context is done.
func Interrupt(ctx context.Context) bool {
	appSignal := make(chan os.Signal, 1)
	signal.Notify(appSignal, os.Interrupt)
	select {
	case <-appSignal:
		// log.Info().Str("what", "Caught an interrupt signal").Msg(Location())
		return true
	case <-ctx.Done():
		// log.Info().Str("what", "Context is done").Msg(Location())
		return false
	}
}
