package commons

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

func DefaultString(x, defaultValue string) string {
	if x != "" {
		return x
	}

	return defaultValue
}

// +----------------+
// | Time utilities |
// +----------------+

func AlignTime(t time.Time, interval time.Duration) time.Time {
	uni := t.Unix()
	sec := int64(interval.Seconds())

	return time.Unix(uni-uni%sec, 0).UTC()
}

func TimeFromTimestamp(x int64) time.Time {
	sec := x / 1000
	msec := x % 1000
	nsec := msec * 1000 * 1000
	t := time.Unix(sec, nsec)

	return t.UTC()
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

// +--------------+
// | Constructors |
// +--------------+

// TradeFromStrings expects an array of the following format:
// - [0] TradeID,
// - [1] Price,
// - [2] Quantity,
// - [3] Cost,
// - [4] Time,
// - [5] Buyer is maker.
func TradeFromStrings(symbol string, xs []string) Trade {
	// 0: TradeID.
	tradeID, err := strconv.ParseInt(xs[0], 10, 64)
	if err != nil {
		panic(err)
	}

	// 1: Price.
	price, err := strconv.ParseFloat(xs[1], 64)
	if err != nil {
		panic(err)
	}

	// 2: Quantity.
	quantity, err := strconv.ParseFloat(xs[2], 64)
	if err != nil {
		panic(err)
	}

	// 3: Cost.

	// 4: Time.
	timestamp, err := strconv.ParseInt(xs[4], 10, 64)
	if err != nil {
		panic(err)
	}

	tradeTime := TimeFromTimestamp(timestamp)

	// 5: Buyer is maker.
	buyerIsMaker, err := strconv.ParseBool(xs[5])
	if err != nil {
		panic(err)
	}

	trade := Trade{
		TradeID:      tradeID,
		Time:         tradeTime,
		Symbol:       symbol,
		Price:        price,
		Quantity:     quantity,
		BuyerIsMaker: buyerIsMaker,
	}

	return trade
}
