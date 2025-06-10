package dumper

import (
	"time"

	"github.com/karosown/katool-go/util/dateutil"

	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/net/format"
	"github.com/karosown/katool-go/sys"
)

type SpecTimeUtil[T any] struct {
	specTimes    []*dateutil.PeriodTime
	SyncMode     bool
	ExcludeEmpty bool
}

type TimeDumpTask[T any] struct {
	Util[T]
	dateutil.PeriodTime
}

func (d *SpecTimeUtil[T]) Exec(exec func(start, end time.Time) []T, dumpNode ...format.EnDeCodeFormat) *Util[any] {
	toStream := stream.ToStream(&d.specTimes)
	if d.SyncMode {
		toStream.Parallel()
	}
	return &Util[any]{
		toStream.Map(func(i *dateutil.PeriodTime) any {
			ts := exec(i.Start, i.End)
			d2 := &TimeDumpTask[T]{
				Util[T]{
					ts,
					nil,
					d.SyncMode,
					d.ExcludeEmpty,
				},
				*i,
			}

			dump, err := d2.Dump(dumpNode...)
			if err != nil {
				return err
			}
			return dump
		}).ToList(),
		nil,
		d.SyncMode,
		d.ExcludeEmpty,
	}
}

func (d *SpecTimeUtil[T]) Sync() *SpecTimeUtil[T] {
	d.SyncMode = true
	return d
}

func (d *TimeDumpTask[T]) Dump(dumpNode ...format.EnDeCodeFormat) (any, error) {
	if d.ExcludeEmpty && d.data == nil {
		return nil, nil
	}
	if cutil.IsEmpty(dumpNode) && nil == d.dumpChain {
		sys.Panic("Must Not Nil of Dump Node Chain.")
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
