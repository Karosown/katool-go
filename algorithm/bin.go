package algorithm

// NumOfTwoMultiply 计算2的幂次乘法数量
// NumOfTwoMultiply calculates the number of power-of-two multiplications
func NumOfTwoMultiply(n int) int {
	i := 1
	for n > 1 {
		n >>= 1
		i++
	}
	return i
}

// NumOfOneInBin 计算二进制中1的个数
// NumOfOneInBin counts the number of 1s in binary representation
func NumOfOneInBin(n int) int {
	i := 0
	for n > 0 {
		n &= n - 1
	}
	return i
}
