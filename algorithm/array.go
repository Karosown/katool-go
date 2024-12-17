package algorithm

func MergeSortedArrayWithEntity[T any](orderBy func(a, b T) bool) func(any, any) any {
	return func(cntValue any, nxt any) any {
		ts := cntValue.([]any)
		nxts := nxt.([]any)
		if len(nxts) == 0 {
			return ts
		}
		lenLast := len(ts)
		lenNxt := len(nxts)
		rest := make([]any, 0)
		l := 0
		r := 0
		for l < lenLast && r < lenNxt {
			total := ts[l].(T)
			current := nxts[r].(T)
			if orderBy(total, current) {
				rest = append(rest, total)
				l++
			} else {
				rest = append(rest, current)
				r++
			}
		}
		if l < lenLast {
			rest = append(rest, ts[l:lenLast]...)
		}
		if r < lenNxt {
			rest = append(rest, nxts[r:lenNxt]...)
		}
		return rest
	}
}
func MergeSortedArrayWithPrimaryData[T any](desc bool, orderBy HashComputeFunction) func(any, any) any {
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

func MergeSortedArrayWithPrimaryId[T any](desc bool, orderBy IDComputeFunction) func(any, any) any {
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
