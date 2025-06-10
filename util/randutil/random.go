package randutil

import (
	"math/rand/v2"

	"github.com/karosown/katool-go/sys"
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
		sys.Panic("Int: min cannot be greater than max")
	}
	return rand.IntN(max-min) + min
}

func String(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 设置随机种子
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}
