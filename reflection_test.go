package commons_test

import (
	"testing"

	"github.com/ydm/commons"
)

func TestCopyFloat64(t *testing.T) {
	t.Parallel()

	type Input struct {
		A string
		B string
		C string
	}

	type Output struct {
		A float64
		B float64
		C float64
		D float64
	}

	input := Input{
		A: "1.0",
		B: "2.0",
		C: "4.0",
	}

	output := Output{64, 64, 64, 64}

	if have := commons.ParseFields(&output, &input); have != 3 {
		t.Errorf("have %d, want 3", have)
	}

	if output.A != 1 {
		t.Errorf("have %f, want 1", output.A)
	}

	if output.B != 2 {
		t.Errorf("have %f, want 2", output.B)
	}

	if output.C != 4 {
		t.Errorf("have %f, want 4", output.C)
	}

	if output.D != 64 {
		t.Errorf("have %f, want 64", output.D)
	}
}
