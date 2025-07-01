package container

import (
	"errors"
	"testing"

	"github.com/karosown/katool-go/container/optional"
)

type Temp[T any] struct {
	Status  int
	Data    T
	Message string
}

func TestOptCase(t *testing.T) {
	o := &optional.OptSwitch[Temp[any]]{}
	i := 5
	submit, err := o.Case(i == 1, func() (*Temp[any], error) {
		return &Temp[any]{0, nil, "1"}, nil
	}, func() (*Temp[any], error) {
		t.Log("case 1 否")
		return nil, errors.New("")
	}, func() (*Temp[any], error) {
		t.Log("case 2 否")
		return nil, nil
	}).Case(i == 2, func() (*Temp[any], error) {
		return &Temp[any]{0, nil, "2"}, nil
	}).Case(i == 3, func() (*Temp[any], error) {
		return &Temp[any]{0, nil, "3"}, nil
	}).Case(i == 5, func() (*Temp[any], error) {
		t.Log("success")
		return &Temp[any]{0, nil, "5"}, nil
	}).
		Break().
		Case(i >= 5, func() (*Temp[any], error) {
			t.Log("success1")
			return &Temp[any]{0, nil, ">=5"}, nil
		}).Default(func(res *Temp[any], err error) (*Temp[any], error) {
		return res, err
	}).Submit()
	if err != nil {
		t.Error(err)
	}
	t.Log(submit)
}

func TestOptCaseNil(t *testing.T) {
	o := &optional.OptSwitch[Temp[any]]{}
	i := 5
	submit, err := o.
		Case(i == 1, func() (*Temp[any], error) {
			return &Temp[any]{0, nil, "1"}, nil
		},
			optional.IdentityOnlyErr[*Temp[any]](errors.New("123")),
			optional.IdentityErr[*Temp[any]](nil, errors.New("456")),
			func() (*Temp[any], error) {
				t.Log("case 3 否")
				return nil, nil
			}).
		Case(i == 2, func() (*Temp[any], error) {
			return &Temp[any]{0, nil, "2"}, nil
		}).
		Case(i == 3, func() (*Temp[any], error) {
			return &Temp[any]{0, nil, "3"}, nil
		}).
		Case(i == 5, func() (*Temp[any], error) {
			t.Log("success")
			return &Temp[any]{0, nil, "5"}, nil
		}).
		Break().
		Case(i >= 5, func() (*Temp[any], error) {
			t.Log("success1")
			return &Temp[any]{0, nil, ">=5"}, nil
		}).Default(func(res *Temp[any], err error) (*Temp[any], error) {
		return res, err
	}).Submit()
	if err != nil {
		t.Error(err)
	}
	t.Log(submit)
}
