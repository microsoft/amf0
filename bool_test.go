package amf0

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func bytesOf(v AmfType) []byte {
	buf := new(bytes.Buffer)
	v.Encode(buf)

	return buf.Bytes()
}

func TestBooleanBuildsAndEncodes(t *testing.T) {
	b1, b2 := new(Bool), new(Bool)

	*b1 = true
	*b2 = false

	assert.Equal(t, []byte{1}, bytesOf(b1))
	assert.Equal(t, []byte{0}, bytesOf(b2))
}

func TestBooleanDecodes(t *testing.T) {
	o := new(Bool)
	err := o.Decode(bytes.NewReader([]byte{1}))

	assert.Nil(t, err)
	assert.True(t, bool(*o))
}

func BenchmarkBooleanDecode(b *testing.B) {
	data := []byte{0}
	out := new(Bool)

	for i := 0; i < b.N; i++ {
		out.Decode(bytes.NewReader(data))
	}
}

func BenchmarkBooleanEncode(b *testing.B) {
	in := Bool(true)

	for i := 0; i < b.N; i++ {
		in.Encode(ioutil.Discard)
	}
}
