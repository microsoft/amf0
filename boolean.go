package amf0

import (
	"io"
)

type Bool bool

var _ AmfType = new(Bool)

// Implements AmfType.Decode
func (b *Bool) Decode(r io.Reader) error {
	bytes, err := readBytes(r, 1)
	if err != nil {
		return err
	}

	*b = Bool(bytes[0] != 0)

	return nil
}

// Implements AmfType.Encode
func (b *Bool) Encode(w io.Writer) (int, error) {
	var buf [1]byte
	if bool(*b) == true {
		buf[0] = 0x1
	}

	return w.Write(buf[:])
}

// Implements AmfType.Marker
func (b *Bool) Marker() byte { return 0x01 }
