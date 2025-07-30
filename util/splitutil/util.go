package splitutil

type Segment[T any] struct {
	Begin T
	End   T
}
type Segments[T any] []Segment[T]

func NumberSplit[T int | int8 | int16 | int32 | int64 | float32 | float64 | byte | rune](b, e, spec T) Segments[T] {
	var res Segments[T]
	for i := b; i < e; i += spec {
		if i+spec < e {
			res = append(res, Segment[T]{Begin: i, End: i + spec - 1})
		} else {
			res = append(res, Segment[T]{Begin: i, End: e})
		}
	}
	return res
}
