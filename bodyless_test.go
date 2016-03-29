package amf0

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullBuildsAndEncodes(t *testing.T) {
	buf := new(bytes.Buffer)
	n, err := new(Null).Encode(buf)

	assert.Nil(t, err)
	assert.Equal(t, 0, n)
	assert.Empty(t, buf.Bytes())
}

func TestNullDecodes(t *testing.T) {
	o := new(Null)

	err := o.Decode(new(bytes.Buffer))
	assert.Nil(t, err)
}

func BenchmarkNullDecode(b *testing.B) {
	out := new(Null)

	for i := 0; i < b.N; i++ {
		out.Decode(bytes.NewReader([]byte{}))
	}
}

func BenchmarkNullEncode(b *testing.B) {
	n := new(Null)

	for i := 0; i < b.N; i++ {
		n.Encode(ioutil.Discard)
	}
}

func TestUndefinedBuildsAndEncodes(t *testing.T) {
	buf := new(bytes.Buffer)
	n, err := new(Undefined).Encode(buf)

	assert.Nil(t, err)
	assert.Equal(t, 0, n)
	assert.Empty(t, buf.Bytes())
}

func TestUndefinedDecodes(t *testing.T) {
	o := new(Undefined)

	err := o.Decode(new(bytes.Buffer))
	assert.Nil(t, err)
}

func BenchmarkUndefinedDecode(b *testing.B) {
	out := new(Undefined)

	for i := 0; i < b.N; i++ {
		out.Decode(bytes.NewReader([]byte{}))
	}
}

func BenchmarkUndefinedEncode(b *testing.B) {
	n := new(Undefined)

	for i := 0; i < b.N; i++ {
		n.Encode(ioutil.Discard)
	}
}
