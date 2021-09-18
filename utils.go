package commons

import (
	"runtime"
	"strings"
)

// Location returns the name of the caller function.
func Location() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	xs := strings.SplitAfterN(fn.Name(), "/", 3) // nolint: gomnd

	return xs[len(xs)-1]
}

// Location2 returns the name of the caller of the caller function.
func Location2() string {
	pc, _, _, ok := runtime.Caller(2) // nolint:gomnd
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	xs := strings.SplitAfterN(fn.Name(), "/", 3) // nolint: gomnd

	return xs[len(xs)-1]
}

func Location3() string {
	pc, _, _, ok := runtime.Caller(3) // nolint:gomnd
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	xs := strings.SplitAfterN(fn.Name(), "/", 3) // nolint: gomnd

	return xs[len(xs)-1]
}

//nolint:gomnd
func Round2(x float64) float64 {
	return float64(int(x*100)) / 100.0
}

//nolint:gomnd
func Round4(x float64) float64 {
	return float64(int(x*10000)) / 10000.0
}
