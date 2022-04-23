package commons

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
)

// +----------+
// | Location |
// +----------+

// Location returns the name of the caller function.
func Location() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	xs := strings.SplitAfterN(fn.Name(), "/", 3)

	return xs[len(xs)-1]
}

// Location2 returns the name of the caller of the caller function.
func Location2() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	xs := strings.SplitAfterN(fn.Name(), "/", 3)

	return xs[len(xs)-1]
}

func Location3() string {
	pc, _, _, ok := runtime.Caller(3)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	xs := strings.SplitAfterN(fn.Name(), "/", 3)

	return xs[len(xs)-1]
}

// +----------+
// | Rounding |
// +----------+

func Round2(x float64) float64 {
	return float64(int(x*100)) / 100.0
}

func Round4(x float64) float64 {
	return float64(int(x*10000)) / 10000.0
}

// +--------+
// | SIGINT |
// +--------+

func GoInterrupt(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		Checker.Push()
		defer Checker.Pop()

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

// +------+
// | Time |
// +------+

func TimeFromMs(x int64) time.Time {
	sec := x / 1000
	msec := x % 1000
	nsec := msec * 1000 * 1000
	t := time.Unix(sec, nsec)

	return t.UTC()
}
