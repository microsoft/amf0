package amf0

import (
	"io"
	"math"
)

type Number struct {
	encoded []byte
	num     float64
}

var _ AmfType = &Number{}

func NewNumber() *Number {
	return &Number{}
}

// Implements AmfType.Decode
func (n *Number) Decode(r io.Reader) error {
	bytes, err := readBytes(r, 8)
	if err != nil {
		return err
	}

	_, err = n.DecodeFrom(bytes, 0)
	return err
}

// Implements AmfType.DecodeFrom
func (n *Number) DecodeFrom(slice []byte, pos int) (int, error) {
	if pos+7 >= len(slice) {
		return 0, io.EOF
	}

	bytes := getUint64(slice, pos)
	n.encoded = slice[0:7]
	n.num = math.Float64frombits(bytes)
	return 8, nil
}

// Gets the contained number
func (n *Number) GetNumber() float64 {
	return n.num
}

// Sets the contained number.
func (n *Number) SetNumber(num float64) {
	bytes := math.Float64bits(num)
	if n.encoded == nil {
		n.encoded = make([]byte, 8)
	}

	putUint64(n.encoded, 0, bytes)
	n.num = num
}

// Implements AmfType.Encode
func (n *Number) Encode(w io.Writer) {
	w.Write(n.encoded)
}

// Implements AmfType.EncodeTo
func (n *Number) EncodeTo(slice []byte, pos int) {
	copy(slice[pos:], n.encoded)
}

// Implements AmfType.EncodeBytes
func (n *Number) EncodeBytes() []byte {
	return n.encoded
}

func (n *Number) Marker() byte {
	return MARKER_NUMBER
}
