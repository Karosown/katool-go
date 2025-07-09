package randutil

import (
	"math/rand/v2"

	"github.com/karosown/katool-go/sys"
)

// Int 返回[min, max)范围内的随机整数（包含min，不包含max）
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

// String 生成指定长度的随机字符串
// String generates a random string of specified length
func String(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 设置随机种子
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}
