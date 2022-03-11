package commons

import (
	"os"
	"runtime/pprof"
)

// +----------+
// | Profiler |
// +----------+

type Profiler struct {
	Filename string
}

func NewProfiler(filename string) Profiler {
	p := Profiler{filename}
	f, err := os.Create(p.Filename)
	if err != nil {
		panic(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		panic(err)
	}
	return p
}

func (p Profiler) Stop() {
	pprof.StopCPUProfile()
}
