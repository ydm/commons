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
type ResourceChecker struct {
	m         sync.Mutex
	resources map[string]int
}

// Checker is a default global instance of ResourceChecker.
//
//nolint:gochecknoglobals
var Checker = ResourceChecker{
	m:         sync.Mutex{},
	resources: make(map[string]int),
}

func (c *ResourceChecker) Push(xs ...string) {
	var name string

	switch len(xs) {
	case 0:
		name = Location2()
	case 1:
		name = xs[0]
	default:
		panic("invalid argument")
	}

	c.m.Lock()
	c.resources[name]++
	c.m.Unlock()
}

func (c *ResourceChecker) Pop(xs ...string) {
	var name string

	switch len(xs) {
	case 0:
		name = Location2()
	case 1:
		name = xs[0]
	default:
		panic("invalid argument")
	}

	c.m.Lock()
	c.resources[name]--
	c.m.Unlock()
}

// CheckerAssert should be defer-called in main().
func (c *ResourceChecker) Assert() {
	What(log.Debug(), "checking resources...")
	time.Sleep(1 * time.Second)

	c.m.Lock()
	defer c.m.Unlock()

	leak := false

	for k, v := range c.resources {
		if v != 0 {
			leak = true

			What(log.Warn().Int("counter", v).Str("unit", k), "leaked resource")
		}
	}

	if !leak {
		What(log.Debug(), "no leaked resources detected")
	}
}
