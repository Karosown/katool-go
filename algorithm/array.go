package algorithm

// MergeSortedArrayWithEntity 使用自定义比较函数合并两个有序数组
// MergeSortedArrayWithEntity merges two sorted arrays using a custom comparison function
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

// MergeSortedArrayWithPrimaryData 使用哈希函数合并两个有序数组
// MergeSortedArrayWithPrimaryData merges two sorted arrays using a hash function
func MergeSortedArrayWithPrimaryData[T any](desc bool, orderBy HashComputeFunction[T]) func(any, any) any {
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

// MergeSortedArrayWithPrimaryId 使用ID函数合并两个有序数组
// MergeSortedArrayWithPrimaryId merges two sorted arrays using an ID function
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
