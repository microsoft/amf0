package amf0

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullBuildsAndEncodes(t *testing.T) {
	s := NewNull()
	assert.Equal(t, make([]byte, 0), s.EncodeBytes())
}

func TestNullDecodes(t *testing.T) {
	bytes := make([]byte, 0)

	o := NewNull()
	err := o.Decode(&reluctantReader{src: bytes})
	assert.Nil(t, err)
}

func BenchmarkNullDecode(b *testing.B) {
	bytes := []byte{}
	out := NewNull()

	for i := 0; i < b.N; i++ {
		out.DecodeFrom(bytes, 0)
	}
}

func BenchmarkNullEncode(b *testing.B) {
	n := NewNull()

	for i := 0; i < b.N; i++ {
		n.EncodeBytes()
	}
}
