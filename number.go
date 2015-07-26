package amf0

import (
	"io"
	"math"
)

type Number struct {
	num float64
}

var _ AmfType = &Number{}

// Creates a new Number type, with an optional initial value.
func NewNumber(num ...float64) *Number {
	n := &Number{}
	if len(num) == 1 {
		n.SetNumber(num[0])
	}

	return n
}

// Implements AmfType.Decode
func (n *Number) Decode(r io.Reader) error {
	bytes, err := readBytes(r, 8)
	if err != nil {
		return err
	}

	rawUint := getUint64(bytes, 0)
	n.num = math.Float64frombits(rawUint)
	return nil
}

// Gets the contained number
func (n *Number) GetNumber() float64 {
	return n.num
}

// Sets the contained number.
func (n *Number) SetNumber(num float64) {
	n.num = num
}

// Implements AmfType.Encode
func (n *Number) Encode(w io.Writer) (int, error) {
	return w.Write(n.EncodeBytes())
}

// Implements AmfType.EncodeBytes
func (n *Number) EncodeBytes() []byte {
	bytes := make([]byte, 9)

	bytes[0] = MARKER_NUMBER
	putUint64(bytes, 1, math.Float64bits(n.num))

	return bytes
}

// Implements AmfType.Marker
func (n *Number) Marker() byte {
	return MARKER_NUMBER
}
