package util

type ObjUtil interface {
	IsEmpty() bool
	IsNotEmpty() bool
	IsNil() bool
	IsNotNil() bool
}
type StringUtil struct {
}
