package dumper

import (
	"time"

	"github.com/karosown/katool-go/util/dateutil"

	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/net/format"
	"github.com/karosown/katool-go/sys"
)

type Util[T any] struct {
	data         []T
	dumpChain    format.EnDeCodeFormat
	SyncMode     bool
	ExcludeEmpty bool
}

func (d *Util[T]) Period(s, e time.Time, spec time.Duration) *SpecTimeUtil[T] {

	return &SpecTimeUtil[T]{
		dateutil.GetPeriods(s, e, spec),
		d.SyncMode,
		d.ExcludeEmpty,
	}
}
func (d *Util[T]) Sync() *Util[T] {
	d.SyncMode = true
	return d
}
func (d *Util[T]) Exec(exec func() []T) *Util[T] {
	d.data = exec()
	return d
}

func (d *Util[T]) Dump(dumpNode ...format.EnDeCodeFormat) (any, error) {
	if d.ExcludeEmpty && d.data == nil {
		return nil, nil
	}
	if cutil.IsEmpty(dumpNode) && nil == d.dumpChain {
		sys.Panic("Must Not Nil of Dump Node Chain.")
	}
	if !cutil.IsEmpty(dumpNode) {
		head, iter := dumpNode[0], dumpNode[0]
		reduce := stream.ToStream(&dumpNode).Reduce(iter, func(cntValue any, nxt format.EnDeCodeFormat) any {
			iter = cntValue.(format.EnDeCodeFormat)
			//iter.Decode(d.data, nil)
			iter.Then(nxt)
			iter = nxt
			return iter
		}, nil)
		if reduce == nil {
			sys.Panic("reduce is nil")
		}
		d.dumpChain = head
	}
	return d.dumpChain.SystemDecode(d.dumpChain, d.data, nil)
}

func (d *Util[T]) Result() []T {
	return d.data
}
