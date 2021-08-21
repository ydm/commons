package commons

import "time"

//nolint:gomnd
func TimeFromMs(x int64) time.Time {
	sec := x / 1000
	msec := x % 1000
	nsec := msec * 1000 * 1000
	t := time.Unix(sec, nsec)

	return t.UTC()
}
