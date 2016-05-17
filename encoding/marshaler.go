package encoding

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/WatchBeam/amf0"
)

// Marshaler is a type capable of marshaling structs into their AMF-encoded
// equivalents.
type Marshaler struct {
	i *amf0.Identifier
}

// Marshal is a method which delegates into the Marshaler type to Marshal a
// struct into its AMF-encoded equivalent.
func Marshal(v interface{}) ([]byte, error) {
	return NewMarshaler().Marshal(v)
}

// NewMarshaler constructs a new Marshaler.
func NewMarshaler() *Marshaler {
	return &Marshaler{
		i: amf0.DefaultIdentifier,
	}
}

// Marshal marshals some interface{} into its AMF-encoded equal. It passes
// through each field of a type one-by-one and marshals it by converting to its
// AMF type. If a field is already an AMF type, it marshals it directly. It does
// not recurse to embedded fields.
//
// If the field is nil (i.e., an uninitialized pointer), then an amf0.Null will
// be written, instead of the actual type.
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
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return new(amf0.Null), nil
	}

	amf := m.i.NewMatchingTypeFromValue(val)
	if amf == nil {
		return nil, noMatchingType{val.Type()}
	}

	amft := reflect.ValueOf(amf)

	var toType reflect.Type
	if val.Kind() == reflect.Ptr {
		toType = amft.Type()
	} else {
		toType = amft.Type().Elem()
	}

	if !val.Type().ConvertibleTo(toType) {
		return nil, typeUnassignable{val.Type(), amft.Type().Elem()}
	}
	amft.Elem().Set(val.Convert(toType))

	return amf, nil
}

type noMatchingType struct {
	typ reflect.Type
}

func (e noMatchingType) Error() string {
	return fmt.Sprintf("amf0/encoding: no matching type for %s", e.typ.String())
}

type typeUnassignable struct {
	assign reflect.Type
	to     reflect.Type
}

func (e typeUnassignable) Error() string {
	return fmt.Sprintf("amf0/encoding: cannot assign type %s to %s",
		e.assign, e.to)
}
