package amf0

import (
	"io"
)

type NullType struct{}

var _ AmfType = &NullType{}

func NewNull() *NullType {
	return &NullType{}
}

// Decode implements AmfType.Decode
func (n *NullType) Decode(r io.Reader) error {
	return nil
}

// DecodeFrom implements AmfType.DecodeFrom
func (n *NullType) DecodeFrom(slice []byte, pos int) (int, error) {
	return 0, nil
}

// Encode implements AmfType.Encode
func (n *NullType) Encode(w io.Writer) {
}

// EncodeTo implements AmfType.EncodeTo
func (n *NullType) EncodeTo(slice []byte, pos int) {
}

// EncodeBytes implements AmfType.EncodeBytes
func (n *NullType) EncodeBytes() []byte {
	return make([]byte, 0)
}

func (n *NullType) Marker() byte {
	return MARKER_NULL
}
