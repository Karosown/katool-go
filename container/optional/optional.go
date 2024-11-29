package optional

func IsTrue[T any](valid bool, isTrue T, isFalse T) T {
	if valid {
		return isTrue
	} else {
		return isFalse
	}
}
