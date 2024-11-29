package remote

import (
	"fmt"
	"reflect"

	"github.com/karosown/katool/log"
)

type BackData any
type EnDeCodeFormat interface {
	SetLogger(logger log.Logger)
	GetLogger() log.Logger
	ValidEncode(obj any) (bool, error)
	ValidDecode(encode any) (bool, error)
	Encode(obj any) (any, error)
	Decode(encode any, backDao any) (any, error)
	Then(format EnDeCodeFormat) EnDeCodeFormat
	SystemEncode(self EnDeCodeFormat, obj any) *string
	SystemDecode(self EnDeCodeFormat, encode any, back any) any
}

type DefaultEnDeCodeFormat struct {
	self    EnDeCodeFormat
	next    EnDeCodeFormat
	backDao any
	logger  log.Logger
}

func (e *DefaultEnDeCodeFormat) SetLogger(logger log.Logger) {
	e.logger = logger
}
func (e *DefaultEnDeCodeFormat) GetLogger() log.Logger {
	return e.logger
}

func (e *DefaultEnDeCodeFormat) ValidEncode(obj any) (bool, error) {
	return true, nil
}
func (e *DefaultEnDeCodeFormat) ValidDecode(encode any) (bool, error) {
	return true, nil
}

/**
 * 工具链
 */
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

func (e *DefaultEnDeCodeFormat) SystemDecode(self EnDeCodeFormat, encode any, back any) any {
	if e.self == nil {
		e.self = self
	}
	valid, err := e.self.ValidDecode(encode)
	if !valid {
		e.logger.Error(fmt.Errorf("DefaultEnDeCodeFormat.SystemDecode.Valid failed: %v", err))
		return nil
	}
	backData, err := e.self.Decode(encode, back)
	if err != nil {
		e.logger.Error(fmt.Errorf("DefaultEnDeCodeFormat.SystemDecode failed: %v", err))
		return nil
	}
	if e.next != nil {
		res := e.next.SystemDecode(e.next, backData, back)
		return res
	}
	return backData
}

func (e *DefaultEnDeCodeFormat) SystemEncode(self EnDeCodeFormat, obj any) *string {
	if e.self == nil {
		e.self = self
	}
	valid, err := e.self.ValidEncode(obj)
	if !valid {
		fmt.Errorf("DefaultEnDeCodeFormat.SystemEncode.Valid failed: %v", err)
		return nil
	}
	encode, err := e.self.Encode(obj)
	if err != nil {
		fmt.Errorf("DefaultEnDeCodeFormat.SystemEncode failed: %v", err)
		return nil
	}
	if e.next != nil {
		return e.next.SystemEncode(e.self, encode)
	} else {
		s := encode.(string)
		return &s
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
