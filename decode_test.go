package amf0

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getEncoded() []byte {
	s := NewString()
	s.SetBody("こんにちは")
	n := NewNumber()
	n.SetNumber(42)
	return append(s.EncodeBytes(), n.EncodeBytes()...)
}

func TestDecoderComplete(t *testing.T) {
	bytes := getEncoded()

	r := &reluctantReader{src: bytes}
	kind, err := Decode(r)
	assert.Nil(t, err)
	assert.Equal(t, "こんにちは", kind.(*String).GetBody())
	kind, err = Decode(r)
	assert.Nil(t, err)
	assert.Equal(t, float64(42), kind.(*Number).GetNumber())
	_, err = Decode(r)
	assert.NotNil(t, err)
}

func BenchmarkDecoder(b *testing.B) {
	s := NewBoolean()
	s.Set(true)
	data := s.EncodeBytes()

	for i := 0; i < b.N; i++ {
		Decode(bytes.NewReader(data))
	}
}
