package format

import (
	"fmt"
	"reflect"

	"github.com/karosown/katool-go/sys"
	"github.com/karosown/katool-go/xlog"
)

type BackData any
type EnDeCodeFormat interface {
	SetLogger(logger xlog.Logger)
	GetLogger() xlog.Logger
	ValidEncode(obj any) (bool, error)
	ValidDecode(encode any) (bool, error)
	Encode(obj any) (any, error)
	Decode(encode any, backDao any) (any, error)
	Then(format EnDeCodeFormat) EnDeCodeFormat
	SystemEncode(self EnDeCodeFormat, obj any) (*string, error)
	SystemDecode(self EnDeCodeFormat, encode any, back any) (any, error)
}

type DefaultEnDeCodeFormat struct {
	self    EnDeCodeFormat
	next    EnDeCodeFormat
	backDao any
	logger  xlog.Logger
}

func (e *DefaultEnDeCodeFormat) SetLogger(logger xlog.Logger) {
	e.logger = logger
}
func (e *DefaultEnDeCodeFormat) GetLogger() xlog.Logger {
	return e.logger
}

func (e *DefaultEnDeCodeFormat) ValidEncode(obj any) (bool, error) {
	return true, nil
}
func (e *DefaultEnDeCodeFormat) ValidDecode(encode any) (bool, error) {
	return true, nil
}

// Then 工具链
func (e *DefaultEnDeCodeFormat) Then(format EnDeCodeFormat) EnDeCodeFormat {
	e.next = format
	return format
}

func (e *DefaultEnDeCodeFormat) Encode(obj any) (any, error) {
	return nil, fmt.Errorf("DefaultEnDeCodeFormat.Encode called without implementation")
}

func (e *DefaultEnDeCodeFormat) Decode(encode any, back any) (any, error) {
	err := fmt.Errorf("DefaultEnDeCodeFormat.Decode called without implementation")
	return nil, err
}

func (e *DefaultEnDeCodeFormat) SystemDecode(self EnDeCodeFormat, encode any, back any) (any, error) {
	if e.self == nil {
		e.self = self
	}
	valid, err := e.self.ValidDecode(encode)
	if !valid {
		err = fmt.Errorf("DefaultEnDeCodeFormat.SystemDecode.Valid failed: %v", err)
		if nil != e.logger {
			e.logger.Error(err)
		} else {
			sys.Warn(err.Error())
		}
		return nil, err
	}
	backData, err := e.self.Decode(encode, back)
	if err != nil {
		err := fmt.Errorf("DefaultEnDeCodeFormat.SystemDecode failed: %v", err)
		if nil != e.logger {
			e.logger.Error(err)
		} else {
			sys.Warn(err.Error())
		}
		return nil, err
	}
	if e.next != nil {
		res, err := e.next.SystemDecode(e.next, backData, back)
		return res, err
	}
	return backData, err
}

func (e *DefaultEnDeCodeFormat) SystemEncode(self EnDeCodeFormat, obj any) (*string, error) {
	if e.self == nil {
		e.self = self
	}
	valid, err := e.self.ValidEncode(obj)
	if !valid {
		return nil, fmt.Errorf("DefaultEnDeCodeFormat.SystemEncode.Valid failed: %v", err)
	}
	encode, err := e.self.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("DefaultEnDeCodeFormat.SystemEncode failed: %v", err)
	}
	if e.next != nil {
		return e.next.SystemEncode(e.self, encode)
	} else {
		s := encode.(string)
		return &s, nil
	}
}

type BytesDecodeFormatValid struct {
}

func (f *BytesDecodeFormatValid) ValidDecode(encode any) (bool, error) {
	if reflect.TypeOf(encode).String() != "[]uint8" {
		return false, fmt.Errorf("BytesDecodeFormatValid.ValidDecode failed: not []byte")
	}
	return true, nil
}
