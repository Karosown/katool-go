package pool

type Worker[T any, R any, F func(t T) (R, error) | func(t T) R, IN []T | []*T] interface {
	Reader(IN) Worker[T, R, F, IN]
	Work(F)
}

type RunWorker[T any, R any, IN []T | []*T] struct {
}

func (d *RunWorker[T, R, IN]) Work(func(a T) R) {

}
func (d *RunWorker[T, R, IN]) Reader(a IN) Worker[T, R, func(a T) R, IN] {
	return d
}

type RunBackErrWorker[T any, R any, IN []T | []*T] struct {
}

func (d *RunBackErrWorker[T, R, IN]) Work(func(a T) (R, error)) {

}
func (d *RunBackErrWorker[T, R, IN]) Reader(a IN) Worker[T, R, func(a T) (R, error), IN] {
	return d
}
