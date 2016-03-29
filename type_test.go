package amf0_test

import (
	"io"
	"reflect"

	"github.com/WatchBeam/amf0"
	"github.com/stretchr/testify/mock"
)

type MockAmfType struct {
	mock.Mock
}

var _ amf0.AmfType = new(MockAmfType)

func (t *MockAmfType) Marker() byte {
	args := t.Called()

	return args.Get(0).(byte)
}

func (t *MockAmfType) Native() reflect.Type {
	args := t.Called()

	return args.Get(0).(reflect.Type)
}

func (t *MockAmfType) Encode(w io.Writer) (int, error) {
	args := t.Called()

	return args.Int(0), args.Error(1)
}

func (t *MockAmfType) Decode(r io.Reader) error {
	args := t.Called()

	return args.Error(0)
}
