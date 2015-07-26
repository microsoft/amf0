package amf0

import (
	"io"
)

type bodylessType struct {
	marker byte
}

// Decode implements AmfType.Decode
func (b *bodylessType) Decode(r io.Reader) error {
	return nil
}

// Encode implements AmfType.Encode
func (b *bodylessType) Encode(w io.Writer) (int, error) {
	return w.Write([]byte{b.marker})
}

// EncodeBytes implements AmfType.EncodeBytes
func (b *bodylessType) EncodeBytes() []byte {
	return []byte{b.marker}
}

// Implements AmfType.Marker
func (b *bodylessType) Marker() byte {
	return b.marker
}
