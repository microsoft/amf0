package amf0

import (
	"io"
)

type Boolean struct {
	isTrue bool
}

var _ AmfType = &Boolean{}

func NewBoolean() *Boolean {
	return &Boolean{}
}

// Implements AmfType.Decode
func (n *Boolean) Decode(r io.Reader) error {
	bytes, err := readBytes(r, 1)
	if err != nil {
		return err
	}

	n.isTrue = bytes[0] > 0
	return nil
}

// Implements AmfType.DecodeFrom
func (n *Boolean) DecodeFrom(slice []byte, pos int) (int, error) {
	if len(slice) == 0 {
		return 0, io.EOF
	}

	n.isTrue = slice[0] > 0
	return 1, nil
}

// Gets the contained boolean
func (n *Boolean) True() bool {
	return n.isTrue
}

// Sets the contained boolean.
func (n *Boolean) Set(isTrue bool) {
	n.isTrue = isTrue
}

// Implements AmfType.Encode
func (n *Boolean) Encode(w io.Writer) {
	w.Write(n.EncodeBytes())
}

// Implements AmfType.EncodeTo
func (n *Boolean) EncodeTo(slice []byte, pos int) {
	if n.isTrue {
		slice[pos] = 1
	} else {
		slice[pos] = 0
	}
}

// Implements AmfType.EncodeBytes
func (n *Boolean) EncodeBytes() []byte {
	if n.isTrue {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

func (b *Boolean) Marker() byte {
	return MARKER_BOOLEAN
}
