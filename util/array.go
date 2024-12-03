package util

import (
	"github.com/karosown/katool/algorithm"
)

func MergeSortedArrayWithHash[T any](desc bool, orderBy algorithm.HashComputeFunction) func(any, any) any {
	return func(cntValue any, nxt any) any {
		ts := cntValue.([]any)
		nxts := nxt.([]any)
		if len(nxts) == 0 {
			return ts
		}
		lenRe := len(ts)
		lenNxt := len(nxts)
		ress := make([]any, 0)
		l := 0
		r := 0
		for l < lenRe && r < lenNxt {
			current := nxts[r].(T)
			total := ts[l].(T)
			if orderBy(total) > orderBy(current) {
				if desc {
					ress = append(ress, total)
					l++
				} else {
					ress = append(ress, current)
					r++
				}
			} else {
				if desc {
					ress = append(ress, current)
					r++
				} else {
					ress = append(ress, total)
					l++
				}
			}
		}

		if r < lenNxt {
			ress = append(ress, nxts[r:lenNxt]...)
		}
		if l < lenRe {
			ress = append(ress, ts[l:lenRe]...)
		}
		return ress
	}
}

func MergeSortedArrayWithId[T any](desc bool, orderBy algorithm.IDComputeFunction) func(any, any) any {
	return func(cntValue any, nxt any) any {
		ts := cntValue.([]any)
		nxts := nxt.([]any)
		if len(nxts) == 0 {
			return ts
		}
		lenRe := len(ts)
		lenNxt := len(nxts)
		ress := make([]any, 0)
		l := 0
		r := 0
		for l < lenRe && r < lenNxt {
			current := nxts[r].(T)
			total := ts[l].(T)
			if orderBy(total) > orderBy(current) {
				if desc {
					ress = append(ress, total)
					l++
				} else {
					ress = append(ress, current)
					r++
				}
			} else {
				if desc {
					ress = append(ress, current)
					r++
				} else {
					ress = append(ress, total)
					l++
				}
			}
		}

		if r < lenNxt {
			ress = append(ress, nxts[r:lenNxt]...)
		}
		if l < lenRe {
			ress = append(ress, ts[l:lenRe]...)
		}
		return ress
	}
}
