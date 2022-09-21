package numbers

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrInvalidEmptyInput = errors.New("invalid empty input")
	ErrNotNumber         = errors.New("not a number")
)

// DecimalToSatoshis removes the comma in a decimal string
// "12.345" => "12345"
// "0.0230" => "230"
func DecimalToSatoshis(dec string) (string, error) {
	out := strings.TrimLeft(dec, " ")
	out = strings.TrimRight(out, " ")
	out = strings.Replace(out, ".", "", 1)

	// trim left 0's but keep last
	if l := len(out); l >= 2 {
		out = strings.TrimLeft(out[:l-1], "0") + out[l-1:l]
	}

	if len(out) == 0 {
		return "", ErrInvalidEmptyInput
	}

	for _, c := range out {
		if !unicode.IsNumber(c) {
			return "", ErrNotNumber
		}
	}

	return out, nil
}
