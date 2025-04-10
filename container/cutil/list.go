package cutil

func IsEmpty[T any, R []T | []*T](list R) bool {
	return list == nil || len(list) == 0
}
