package commons

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// +---------+
// | Checker |
// +---------+

// Checker is a simple tool to check if everything initialized is subsequently
// deinitialized.  Works from simple open/close calls to gourintes.
type resourceChecker struct {
	m         sync.Mutex
	resources map[string]int
}

// nolint:gochecknoglobals
var defaultResourceChecker = resourceChecker{
	m:         sync.Mutex{},
	resources: make(map[string]int),
}

func CheckerPush(xs ...string) {
	var name string

	switch len(xs) {
	case 0:
		name = Location2()
	case 1:
		name = xs[0]
	default:
		panic("invalid argument")
	}

	defaultResourceChecker.m.Lock()
	defaultResourceChecker.resources[name]++
	defaultResourceChecker.m.Unlock()
}

func CheckerPop(xs ...string) {
	var name string

	switch len(xs) {
	case 0:
		name = Location2()
	case 1:
		name = xs[0]
	default:
		panic("invalid argument")
	}

	defaultResourceChecker.m.Lock()
	defaultResourceChecker.resources[name]--
	defaultResourceChecker.m.Unlock()
}

// CheckerAssert should be defer-called in main().
func CheckerAssert() {
	What(log.Debug(), "checking resources...")
	time.Sleep(1 * time.Second)

	defaultResourceChecker.m.Lock()
	defer defaultResourceChecker.m.Unlock()

	for k, v := range defaultResourceChecker.resources {
		if v != 0 {
			What(log.Warn().Int("counter", v).Str("unit", k), "leaked resource")
		}
	}
}
