package encoding

import (
	"io"
	"reflect"

	"github.com/WatchBeam/amf0"
)

type Unmarshaler struct {
	r io.Reader

	decode amf0.Decoder
}

func Unmarshal(r io.Reader, v interface{}) error {
	return NewUnmarshaler(r).Unmarshal(v)
}

func NewUnmarshaler(r io.Reader) *Unmarshaler {
	return &Unmarshaler{
		r:      r,
		decode: amf0.Decode,
	}
}

func (u *Unmarshaler) Unmarshal(dest interface{}) error {
	v := reflect.ValueOf(dest).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		next, err := u.decode(u.r)
		if err != nil {
			return err
		}

		val := reflect.ValueOf(next).Elem()
		field.Set(val.Convert(field.Type()))
	}

	return nil
}
