package dumper

import (
	"time"

	"github.com/karosown/katool-go/util/dateutil"

	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/net/format"
	"github.com/karosown/katool-go/sys"
)

// Util 数据转储工具
// Util is a data dumping utility
type Util[T any] struct {
	data         []T
	dumpChain    format.EnDeCodeFormat
	SyncMode     bool
	ExcludeEmpty bool
}

// Period 创建周期性时间工具
// Period creates a periodic time utility
func (d *Util[T]) Period(s, e time.Time, spec time.Duration) *SpecTimeUtil[T] {

	return &SpecTimeUtil[T]{
		dateutil.GetPeriods(s, e, spec),
		d.SyncMode,
		d.ExcludeEmpty,
	}
}

// Sync 设置同步模式
// Sync sets synchronization mode
func (d *Util[T]) Sync() *Util[T] {
	d.SyncMode = true
	return d
}

// Exec 执行数据获取函数
// Exec executes a data retrieval function
func (d *Util[T]) Exec(exec func() []T) *Util[T] {
	d.data = exec()
	return d
}

// Dump 转储数据
// Dump dumps the data
func (d *Util[T]) Dump(dumpNode ...format.EnDeCodeFormat) (any, error) {
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
	return d.dumpChain.SystemDecode(d.dumpChain, d.data, nil)
}

// Result 获取结果数据
// Result gets the result data
func (d *Util[T]) Result() []T {
	return d.data
}
