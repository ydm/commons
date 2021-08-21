package commons_test

import (
	"testing"
	"time"

	"github.com/ydm/commons"
)

func TestSmartCopy(t *testing.T) {
	t.Parallel()

	type Input struct {
		A string
		B string
		C string
		D int64
	}

	input := Input{
		A: "1.0",
		B: "2.0",
		C: "something",
		D: 1609459260000,
	}

	type Output struct {
		A float64
		B float32
		C string
		D time.Time
		E int
	}

	output := Output{64, 64, "yep", time.Time{}, 64}

	if err := commons.SmartCopy(&output, &input); err != nil {
		t.Error(err)
	}

	if output.A != 1 {
		t.Errorf("have %f, want 1", output.A)
	}

	if output.B != 2 {
		t.Errorf("have %f, want 2", output.B)
	}

	if output.C != "something" {
		t.Errorf("have %s, want something", output.C)
	}

	if have := output.D.Format(time.RFC3339Nano); have != "2021-01-01T00:01:00Z" {
		t.Errorf("have %s, want 2021-01-01T00:01:00Z", have)
	}

	if output.E != 64 {
		t.Errorf("have %d, want 64", output.E)
	}
}
