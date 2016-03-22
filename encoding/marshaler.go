package encoding

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/WatchBeam/amf0"
)

type Marshaler struct {
	i *amf0.Identifier
}

func Marshal(v interface{}) ([]byte, error) {
	return NewMarshaler().Marshal(v)
}

func NewMarshaler() *Marshaler {
	return &Marshaler{
		i: amf0.DefaultIdentifier,
	}
}

func (m *Marshaler) Marshal(dest interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)

	value := reflect.ValueOf(dest).Elem()
	for i := 0; i < value.NumField(); i++ {
		amf, err := m.convertToAmfType(value.Field(i))
		if err != nil {
			return nil, err
		}

		if _, err = amf0.Encode(amf, buf); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (m *Marshaler) convertToAmfType(val reflect.Value) (amf0.AmfType, error) {
	amfType := m.i.AmfType(val.Interface())
	if amfType == nil {
		return nil, noMatchingType(val.Type().String())
	}

	if !val.Type().ConvertibleTo(amfType) {
		return nil, typeUnassignable{val.Type(), amfType}
	}

	v := reflect.New(amfType).Elem()
	v.Set(val.Convert(amfType))

	return v.Addr().Interface().(amf0.AmfType), nil
}

type noMatchingType string

func (e noMatchingType) Error() string {
	return fmt.Sprintf("amf0/encoding: no matching type for %s", string(e))
}

type typeUnassignable struct {
	assign reflect.Type
	to     reflect.Type
}

func (e typeUnassignable) Error() string {
	return fmt.Sprintf("amf0/encoding: cannot assign type %s to %s",
		e.assign, e.to)
}
