package amf0

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func simpleEncode(n float64) []byte {
	un := math.Float64bits(n)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, un)
	return append([]byte{MARKER_NUMBER}, b...)
}

func TestNumberBuildsAndEncodes(t *testing.T) {
	s := NewNumber()
	s.SetNumber(0x4242)
	assert.Equal(t, simpleEncode(0x4242), s.EncodeBytes())
}

func TestNumberDecodes(t *testing.T) {
	s := NewNumber()
	s.SetNumber(0x4242)
	data := s.EncodeBytes()[1:]

	o := NewNumber()
	err := o.Decode(bytes.NewReader(data))
	assert.Nil(t, err)
	assert.Equal(t, float64(0x4242), o.GetNumber())
}

func BenchmarkNumberDecode(b *testing.B) {
	in := NewNumber()
	in.SetNumber(0x4242)
	data := in.EncodeBytes()[1:]
	out := NewNumber()

	for i := 0; i < b.N; i++ {
		out.Decode(bytes.NewReader(data))
	}
}

func BenchmarkNumberEncode(b *testing.B) {
	in := NewNumber()
	in.SetNumber(0x4242)

	for i := 0; i < b.N; i++ {
		in.EncodeBytes()
	}
}
