package commons_test

import (
	"testing"
	"time"

	"github.com/ydm/commons"
)

func TestTimeFromMs(t *testing.T) {
	t.Parallel()

	f := func(input int64, want string) {
		t.Helper()

		if have := commons.TimeFromMs(input).Format(time.RFC3339Nano); have != want {
			t.Errorf("have %s, want %s", have, want)
		}
	}

	f(1609459259999, "2021-01-01T00:00:59.999Z")
	f(1609459260000, "2021-01-01T00:01:00Z")
	f(1609459260001, "2021-01-01T00:01:00.001Z")
}
