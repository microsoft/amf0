package amf0

import (
	"io"
	"reflect"
)

type Bool bool

var _ AmfType = new(Bool)

// Implements AmfType.Marker
func (b *Bool) Marker() byte { return 0x01 }

// Implements AmfType.Native
func (b *Bool) Native() reflect.Type { return reflect.TypeOf(false) }

// Implements AmfType.Decode
func (b *Bool) Decode(r io.Reader) error {
	var buf [1]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return err
	}

	*b = Bool(buf[0] != 0)

	return nil
}

// Implements AmfType.Encode
func (b *Bool) Encode(w io.Writer) (int, error) {
	var buf [1]byte
	if *b {
		buf[0] = 0x1
	}

	return w.Write(buf[:])
}
