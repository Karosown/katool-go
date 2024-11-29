package algorithm

func NumOfTwoMultiply(n int) int {
	i := 1
	for n > 1 {
		n &= n - 1
		i++
	}
	return i
}
