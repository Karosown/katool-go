package dumper

import (
	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/net/format"
	"github.com/karosown/katool-go/sys"
	"github.com/karosown/katool-go/util/splitutil"
)

type SpecNumUtil[R any, T int | int8 | int16 | int32 | int64 | float32 | float64 | byte | rune] struct {
	specNums     splitutil.Segments[T]
	SyncMode     bool
	ExcludeEmpty bool
}

type SplitDumpTask[R any, T int | int8 | int16 | int32 | int64 | float32 | float64 | byte | rune] struct {
	Util[R]
	splitutil.Segment[T]
}

func (d *SpecNumUtil[R, T]) Exec(exec func(start, end T) []R, dumpNode ...format.EnDeCodeFormat) *Util[R] {
	toStream := stream.ToStream(&d.specNums)
	if d.SyncMode {
		toStream.Parallel()
	}
	list := stream.Cast[[]R](toStream.Map(func(i splitutil.Segment[T]) any {
		ts := exec(i.Begin, i.End)
		d2 := &SplitDumpTask[R, T]{
			Util[R]{
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
		t, ok := dump.([]R)
		if !ok {
			sys.Panic("The Exec Handler Back Type Need Consistent Of SpecNumUtil[T]，also []T")
			return nil
		}
		return t
	})).Reduce([]R{}, func(cntValue any, nxt []R) any {
		return append(cntValue.([]R), nxt...)
	}, func(sum1, sum2 any) any {
		return append(sum1.([]R), sum2.([]R)...)
	}).([]R)
	return &Util[R]{
		list,
		nil,
		d.SyncMode,
		d.ExcludeEmpty,
	}
}

func (d *SpecNumUtil[R, T]) Sync() *SpecNumUtil[R, T] {
	d.SyncMode = true
	return d
}

func (d *SplitDumpTask[R, T]) Dump(dumpNode ...format.EnDeCodeFormat) (any, error) {
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
