package commons

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func DefaultString(x, defaultValue string) string {
	if x != "" {
		return x
	}

	return defaultValue
}

// RandomOrderID uses code and ideas from:
// https://stackoverflow.com/questions/32349807 and
// https://stackoverflow.com/questions/13378815 .
//
// Length of produced client order ID is encoded in the code.  See `seed`.
func RandomOrderID(prefix string) string {
	const seed = 24
	xs := make([]byte, seed)

	if _, err := rand.Read(xs); err != nil {
		panic(err)
	}

	ys := base64.URLEncoding.EncodeToString(xs)
	offset := len(prefix)
	id := fmt.Sprintf("%s%s", prefix, ys[offset:])

	return id
}
