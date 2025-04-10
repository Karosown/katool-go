package collect

import (
	"github.com/karosown/katool-go/collect/lists"
	"github.com/karosown/katool-go/container/stream"
)

func PartitionToStream[T any, RT ~[]T](pattion lists.Batch[T, RT]) *stream.Stream[RT, []RT] {
	return stream.ToStream(&pattion.SplitData)
}
