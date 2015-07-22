package amf0

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBooleanBuildsAndEncodes(t *testing.T) {
	s := NewBoolean()
	s.Set(true)
	assert.Equal(t, []byte{1}, s.EncodeBytes())
	s.Set(false)
	assert.Equal(t, []byte{0}, s.EncodeBytes())
}

func TestBooleanDecodes(t *testing.T) {
	bytes := []byte{1}

	o := NewBoolean()
	err := o.Decode(&reluctantReader{src: bytes})
	assert.Nil(t, err)
	assert.True(t, o.True())

	o = NewBoolean()
	n, err := o.DecodeFrom(bytes, 0)
	assert.Nil(t, err)
	assert.Equal(t, 1, n)
	assert.True(t, o.True())
}

func BenchmarkBooleanDecode(b *testing.B) {
	bytes := []byte{0}
	out := NewBoolean()

	for i := 0; i < b.N; i++ {
		out.DecodeFrom(bytes, 0)
	}
}

func BenchmarkBooleanEncode(b *testing.B) {
	in := NewBoolean()
	in.Set(true)

	for i := 0; i < b.N; i++ {
		in.EncodeBytes()
	}
}
