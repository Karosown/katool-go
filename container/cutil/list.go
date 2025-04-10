package cutil

func IsEmpty[T any, R []T](list R) bool {
	return list == nil || len(list) == 0
}
