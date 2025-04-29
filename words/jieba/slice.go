package jieba

import (
	"strings"

	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/container/xmap"
)

type SplitStrings []string

func (s SplitStrings) String() string {
	return strings.Join(s, ",")
}
func (s SplitStrings) ToStream() *stream.Stream[string, SplitStrings] {
	return stream.ToStream(&s)
}
func (s SplitStrings) Distinct() SplitStrings {
	return stream.ToStream(&s).Distinct().ToList()
}
func (s SplitStrings) Frequency() *xmap.SortedMap[string, int64] {
	newMap := xmap.NewSortedMap[string, int64]()
	s.Distinct().ToStream().ForEach(func(item string) {
		newMap.Set(item, s.ToStream().Filter(func(i string) bool {
			return i == item
		}).Count())
	})
	return newMap
}
