package amf0

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func strBytes(s AmfType) []byte {
	buf := new(bytes.Buffer)
	s.Encode(buf)

	return buf.Bytes()
}

func TestStringBuildsAndEncodes(t *testing.T) {
	s := NewString("hello")

	assert.Equal(t, []byte{0x0, 0x5, 0x68, 0x65, 0x6C,
		0x6C, 0x6F}, strBytes(s))

}

func TestStringEncodingUtf8(t *testing.T) {
	s := NewString("こんにちは")
	assert.Equal(t, []byte{0x0, 0xf, 0xe3, 0x81, 0x93,
		0xe3, 0x82, 0x93, 0xe3, 0x81, 0xab, 0xe3, 0x81, 0xa1, 0xe3,
		0x81, 0xaf}, strBytes(s))
}

func TestStringDecodes(t *testing.T) {
	o := new(String)

	err := o.Decode(bytes.NewReader([]byte{
		0x0, 0xf, 0xe3, 0x81, 0x93, 0xe3, 0x82, 0x93, 0xe3, 0x81, 0xab,
		0xe3, 0x81, 0xa1, 0xe3, 0x81, 0xaf,
	}))

	assert.Nil(t, err)
	assert.Equal(t, "こんにちは", string(*o))
}

func BenchmarkStringDecode(b *testing.B) {
	in := NewString("hello")
	data := strBytes(in)

	out := new(String)
	for i := 0; i < b.N; i++ {
		out.Decode(bytes.NewReader(data))
	}
}

func BenchmarkStringEncode(b *testing.B) {
	in := NewString("hello")

	for i := 0; i < b.N; i++ {
		in.Encode(ioutil.Discard)
	}
}

func TestLongStringBuildsAndEncodes(t *testing.T) {
	s := NewLongString("hello")
	assert.Equal(t, []byte{0x0, 0x0, 0x0, 0x5, 0x68,
		0x65, 0x6C, 0x6C, 0x6F}, strBytes(s))
}

func TestLongStringEncodingUtf8(t *testing.T) {
	s := NewLongString("こんにちは")
	assert.Equal(t, []byte{
		0x00, 0x00, 0x00, 0xf, 0xe3, 0x81, 0x93, 0xe3, 0x82, 0x93, 0xe3,
		0x81, 0xab, 0xe3, 0x81, 0xa1, 0xe3, 0x81, 0xaf,
	}, strBytes(s))
}

func TestLongStringDecodes(t *testing.T) {
	o := new(LongString)
	err := o.Decode(bytes.NewReader([]byte{
		0x00, 0x00, 0x00, 0xf, 0xe3, 0x81, 0x93, 0xe3, 0x82, 0x93, 0xe3,
		0x81, 0xab, 0xe3, 0x81, 0xa1, 0xe3, 0x81, 0xaf,
	}))

	assert.Nil(t, err)
	assert.Equal(t, "こんにちは", string(*o))
}

func BenchmarkLongStringDecode(b *testing.B) {
	data := []byte{
		0x00, 0x00, 0x00, 0xf, 0xe3, 0x81, 0x93, 0xe3, 0x82, 0x93, 0xe3,
		0x81, 0xab, 0xe3, 0x81, 0xa1, 0xe3, 0x81, 0xaf,
	}

	out := new(LongString)
	for i := 0; i < b.N; i++ {
		out.Decode(bytes.NewReader(data))
	}
}

func BenchmarkLongStringEncode(b *testing.B) {
	in := NewLongString("hello")
	for i := 0; i < b.N; i++ {
		in.Encode(ioutil.Discard)
	}
}
