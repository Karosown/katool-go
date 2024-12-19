package collect

import (
	"github.com/karosown/katool/collect/lists"
	"github.com/karosown/katool/container/stream"
)

func PartitionToStream[T any, RT ~[]T](pattion lists.Batch[T, RT]) *stream.Stream[RT, []RT] {
	return stream.ToStream(&pattion.SplitData)
}
