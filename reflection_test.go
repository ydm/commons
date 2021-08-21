package commons_test

import (
	"testing"

	"github.com/ydm/commons"
)

func TestSmartCopy(t *testing.T) {
	t.Parallel()

	type Input struct {
		A string
		B string
		C string
	}

	input := Input{
		A: "1.0",
		B: "2.0",
		C: "something",
	}

	type Output struct {
		A float64
		B float32
		C string
		D int
	}

	output := Output{64, 64, "yep", 64}

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

	if output.D != 64 {
		t.Errorf("have %d, want 64", output.D)
	}
}
