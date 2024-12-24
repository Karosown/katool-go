package util

import (
	"math/rand/v2"
)

// Int returns a random integer in [min, max) range (inclusive).
// Parameters:
//   - min: the lower bound of the range
//   - max: the upper bound of the range
//
// Panics if min > max.
// Returns a random integer between min and max (inclusive).
func Int(min, max int) int {
	if min > max {
		panic("Int: min cannot be greater than max")
	}
	return rand.IntN(max-min) + min
}
