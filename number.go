package amf0

import (
	"encoding/binary"
	"io"
	"math"
)

type Number float64

var _ AmfType = new(Number)

func (n *Number) Marker() byte { return 0x00 }

// Implements AmfType.Decode
func (n *Number) Decode(r io.Reader) error {
	var b [8]byte
	if _, err := r.Read(b[:]); err != nil {
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
