package amf0

import (
	"encoding/binary"
	"io"
	"math"
	"reflect"
)

type Number float64

var _ AmfType = new(Number)

func NewNumber(n float64) *Number {
	num := new(Number)
	*num = Number(n)

	return num
}

func (n *Number) Marker() byte { return 0x00 }

func (n *Number) Native() reflect.Type { return reflect.TypeOf(float64(1)) }

// Implements AmfType.Decode
func (n *Number) Decode(r io.Reader) error {
	var b [8]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return err
	}

	*n = Number(math.Float64frombits(binary.BigEndian.Uint64(b[:])))

	return nil
}

// Implements AmfType.Encode
func (n *Number) Encode(w io.Writer) (int, error) {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, math.Float64bits(float64(*n)))

	return w.Write(bytes)
}
