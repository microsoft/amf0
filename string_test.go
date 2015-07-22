package amf0

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringBuildsAndEncodes(t *testing.T) {
	s := NewString()
	s.SetBody("hello")
	assert.Equal(t, []byte{0x0, 0x5, 0x68, 0x65, 0x6C, 0x6C, 0x6F}, s.EncodeBytes())
	s.SetBody("こんにちは")
	assert.Equal(t, []byte{0x0, 0xf, 0xe3, 0x81, 0x93, 0xe3, 0x82,
		0x93, 0xe3, 0x81, 0xab, 0xe3, 0x81, 0xa1, 0xe3, 0x81, 0xaf}, s.EncodeBytes())
}

func TestStringDecodes(t *testing.T) {
	s := NewString()
	s.SetBody("こんにちは")
	bytes := s.EncodeBytes()

	o := NewString()
	err := o.Decode(&reluctantReader{src: bytes})
	assert.Nil(t, err)
	assert.Equal(t, "こんにちは", o.GetBody())

	o = NewString()
	n, err := o.DecodeFrom(append([]byte{0, 1, 2, 3}, bytes...), 4)
	assert.Nil(t, err)
	assert.Equal(t, 17, n)
	assert.Equal(t, "こんにちは", o.GetBody())
}

func TestStringBounces(t *testing.T) {
	// Maybe a bit surperfluous
	in := "こんにちは"
	for i := 0; i < 100; i++ {
		in = bounceString(in)
	}

	assert.Equal(t, "こんにちは", in)
}

func bounceString(str string) string {
	in := NewString()
	in.SetBody(str)
	out := NewString()
	out.DecodeFrom(in.EncodeBytes(), 0)
	return out.GetBody()
}

func BenchmarkStringDecode(b *testing.B) {
	in := NewString()
	in.SetBody("hello")
	bytes := in.EncodeBytes()
	out := NewString()

	for i := 0; i < b.N; i++ {
		out.DecodeFrom(bytes, 0)
	}
}

func BenchmarkStringEncode(b *testing.B) {
	in := NewString()
	in.SetBody("hello")

	for i := 0; i < b.N; i++ {
		in.EncodeBytes()
	}
}

func TestLongStringBuildsAndEncodes(t *testing.T) {
	s := NewLongString()
	s.SetBody("hello")
	assert.Equal(t, []byte{0x0, 0x0, 0x0, 0x5, 0x68, 0x65, 0x6C, 0x6C, 0x6F}, s.EncodeBytes())
	s.SetBody("こんにちは")
	assert.Equal(t, []byte{0x0, 0x0, 0x0, 0xf, 0xe3, 0x81, 0x93, 0xe3, 0x82,
		0x93, 0xe3, 0x81, 0xab, 0xe3, 0x81, 0xa1, 0xe3, 0x81, 0xaf}, s.EncodeBytes())
}

func TestLongStringDecodes(t *testing.T) {
	s := NewLongString()
	s.SetBody("こんにちは")
	bytes := s.EncodeBytes()

	o := NewLongString()
	err := o.Decode(&reluctantReader{src: bytes})
	assert.Nil(t, err)
	assert.Equal(t, "こんにちは", o.GetBody())

	o = NewLongString()
	n, err := o.DecodeFrom(append([]byte{0, 1, 2, 3}, bytes...), 4)
	assert.Nil(t, err)
	assert.Equal(t, 19, n)
	assert.Equal(t, "こんにちは", o.GetBody())
}

func BenchmarkLongStringDecode(b *testing.B) {
	in := NewLongString()
	in.SetBody("hello")
	bytes := in.EncodeBytes()
	out := NewLongString()

	for i := 0; i < b.N; i++ {
		out.DecodeFrom(bytes, 0)
	}
}

func BenchmarkLongStringEncode(b *testing.B) {
	in := NewLongString()
	in.SetBody("hello")

	for i := 0; i < b.N; i++ {
		in.EncodeBytes()
	}
}
