package amf0

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var arrTestData = []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x03, 0x66,
	0x6f, 0x6f, 0x02, 0x00, 0x03, 0x62, 0x61, 0x72, 0x00, 0x00, 0x09}

func TestArrayDecodes(t *testing.T) {
	o := NewArray()
	err := o.Decode(bytes.NewBuffer(arrTestData))

	assert.Nil(t, err)
	assert.Equal(t, 1, o.Len())

	s, _ := o.String("foo")
	assert.Equal(t, "bar", string(*s))

	_, err = o.Bool("app")
	assert.Equal(t, WrongTypeError, err)
	_, err = o.Bool("foo")
	assert.Equal(t, NotFoundError, err)
}

func TestArrayBuildsAndEncodes(t *testing.T) {
	s := NewArray()
	s.Add("foo", NewString("bar"))

	assert.Equal(t, arrTestData, s.EncodeBytes())
}

func BenchmarkArrayDecode(b *testing.B) {
	out := NewArray()

	for i := 0; i < b.N; i++ {
		out.Decode(bytes.NewReader(arrTestData))
	}
}

func BenchmarkArrayLookup(b *testing.B) {
	out := NewArray()
	out.Decode(bytes.NewReader(arrTestData))

	for i := 0; i < b.N; i++ {
		out.String("foo")
	}
}

func BenchmarkArrayBuild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := NewArray()
		a.Add("app", NewString("myapp"))
		a.Add("type", NewString("nonprivate"))
		a.Add("flashVer", NewString("FMLE/3.0 (compatible; FMSc/1.0)"))
		a.Add("swfUrl", NewString("rtmp://localhost/myapp"))
		a.Add("tcUrl", NewString("rtmp://localhost/myapp"))
		a.EncodeBytes()
	}
}
