package algorithm

func NumOfTwoMultiply(n int) int {
	i := 1
	for n > 1 {
		n >>= 1
		i++
	}
	return i
}

func NumOfOneInBin(n int) int {
	i := 0
	for n > 0 {
		n &= n - 1
	}
	return i
}
