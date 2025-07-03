package lists

import (
	"errors"
	"sync"

	lynxSync "github.com/Tangerg/lynx/pkg/sync"
)

type Batch[T any, RT ~[]T] struct {
	SplitData []RT
}

func AvgPartition[T any](datas []T, totalPage int) Batch[T, []T] {
	total := len(datas)
	return Partition(datas, total/totalPage)
}
func Partition[T any](datas []T, pageSize int) Batch[T, []T] {
	splitData := make([][]T, 0)
	for i := 0; i < len(datas); i += pageSize {
		splitData = append(splitData, datas[i:min(i+pageSize, len(datas))])
	}

	return Batch[T, []T]{
		SplitData: splitData,
	}
}

func (b Batch[T, RT]) ForEach(solve func(pos int, automicDatas []T) error, async bool, limiter *lynxSync.Limiter) error {
	errs := make([]error, 0)
	if async {
		countDownLatch := &sync.WaitGroup{}
		//countDownLatch.Add(len(b.SplitData))
		for i, data := range b.SplitData {
			limiter.Acquire()
			countDownLatch.Add(1)
			go func(limit *lynxSync.Limiter, datas []T, pos int) {
				defer countDownLatch.Done()
				err := solve(pos, datas)
				if err != nil {
					errs = append(errs, err)
				}
				defer limit.Release()
			}(limiter, data, i)
		}
		countDownLatch.Wait()
	} else {
		for i, data := range b.SplitData {
			err := solve(i, data)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	err := errors.Join(errs...)
	return err
}
