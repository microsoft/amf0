package amf0

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func simpleEncode(n float64) []byte {
	un := math.Float64bits(n)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, un)
	return b
}

func TestNumberBuildsAndEncodes(t *testing.T) {
	s := NewNumber()
	s.SetNumber(0x4242)
	assert.Equal(t, simpleEncode(0x4242), s.EncodeBytes())
}

func TestNumberDecodes(t *testing.T) {
	s := NewNumber()
	s.SetNumber(0x4242)
	bytes := s.EncodeBytes()

	o := NewNumber()
	err := o.Decode(&reluctantReader{src: bytes})
	assert.Nil(t, err)
	assert.Equal(t, float64(0x4242), o.GetNumber())

	o = NewNumber()
	n, err := o.DecodeFrom(bytes, 0)
	assert.Nil(t, err)
	assert.Equal(t, 8, n)
	assert.Equal(t, float64(0x4242), o.GetNumber())
}

func BenchmarkNumberDecode(b *testing.B) {
	in := NewNumber()
	in.SetNumber(0x4242)
	bytes := in.EncodeBytes()
	out := NewNumber()

	for i := 0; i < b.N; i++ {
		out.DecodeFrom(bytes, 0)
	}
}

func BenchmarkNumberEncode(b *testing.B) {
	in := NewNumber()
	in.SetNumber(0x4242)

	for i := 0; i < b.N; i++ {
		in.EncodeBytes()
	}
}
