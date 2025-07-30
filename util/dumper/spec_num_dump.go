package dumper

import (
	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/net/format"
	"github.com/karosown/katool-go/sys"
)

type SpecNumUtil[T any] struct {
	specNums     []stream.Entry[int, int]
	SyncMode     bool
	ExcludeEmpty bool
}

type SplitDumpTask[T any] struct {
	Util[T]
	stream.Entry[int, int]
}

func (d *SpecNumUtil[T]) Exec(exec func(start, end int) []T, dumpNode ...format.EnDeCodeFormat) *Util[T] {
	toStream := stream.ToStream(&d.specNums)
	if d.SyncMode {
		toStream.Parallel()
	}
	list := stream.FromAnySlice[[]T, [][]T](toStream.Map(func(i stream.Entry[int, int]) any {
		ts := exec(i.Key, i.Value)
		d2 := &SplitDumpTask[T]{
			Util[T]{
				ts,
				nil,
				d.SyncMode,
				d.ExcludeEmpty,
			},
			i,
		}

		dump, err := d2.Dump(dumpNode...)
		if err != nil {
			return err
		}
		t, ok := dump.([]T)
		if !ok {
			sys.Panic("The Exec Handler Back Type Need Consistent Of SpecNumUtil[T]，also []T")
			return nil
		}
		return t
	}).ToList()).Reduce([]T{}, func(cntValue any, nxt []T) any {
		return append(cntValue.([]T), nxt...)
	}, func(sum1, sum2 any) any {
		return append(sum1.([]T), sum2.([]T)...)
	}).([]T)
	return &Util[T]{
		list,
		nil,
		d.SyncMode,
		d.ExcludeEmpty,
	}
}

func (d *SpecNumUtil[T]) Sync() *SpecNumUtil[T] {
	d.SyncMode = true
	return d
}

func (d *SplitDumpTask[T]) Dump(dumpNode ...format.EnDeCodeFormat) (any, error) {
	if (d.ExcludeEmpty && d.data == nil) || cutil.IsEmpty(dumpNode) && nil == d.dumpChain {
		return d.data, nil
	}
	if !cutil.IsEmpty(dumpNode) {
		head := &format.EmptyEnDecodeFormatNode{}
		reduce := stream.ToStream(&dumpNode).Reduce(head, func(cntValue any, nxt format.EnDeCodeFormat) any {
			iter := cntValue.(format.EnDeCodeFormat)
			iter.Then(nxt)
			return nxt
		}, nil)
		if reduce == nil {
			sys.Panic("reduce is nil")
		}
		d.dumpChain = head
	}
	if d.dumpChain != nil {
		return d.dumpChain.SystemEncode(d.dumpChain, d)
	}
	return d.data, nil
}
