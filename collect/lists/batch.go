package lists

import (
	"errors"
	"sync"

	lynxSync "github.com/Tangerg/lynx/pkg/sync"
)

type Batch[T any, RT ~[]T] struct {
	SplitData []RT
}

func Partition[T any](datas []T, size int) Batch[T, []T] {
	splitData := make([][]T, 0)
	for i := 0; i < len(datas); i += size {
		splitData = append(splitData, datas[i:min(i+size, len(datas))])
	}

	return Batch[T, []T]{
		SplitData: splitData,
	}
}

func (b Batch[T, RT]) ForEach(solve func(automicDatas []T) error, async bool, limiter *lynxSync.Limiter) error {
	errs := make([]error, 0)
	if async {
		countDownLatch := &sync.WaitGroup{}
		countDownLatch.Add(len(b.SplitData))
		for _, data := range b.SplitData {
			limiter.Acquire()
			go func(datas []T) {
				defer countDownLatch.Done()
				err := solve(datas)
				if err != nil {
					errs = append(errs, err)
				}
				defer limiter.Release()
			}(data)
		}
		countDownLatch.Wait()
	} else {
		for _, data := range b.SplitData {
			err := solve(data)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	err := errors.Join(errs...)
	return err
}
