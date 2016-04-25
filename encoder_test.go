package amf0_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/WatchBeam/amf0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSuccesfulEncodingWritesMarkerAndPayload(t *testing.T) {
	buf := new(bytes.Buffer)

	b := amf0.Bool(false)

	n, err := amf0.Encode(&b, buf)

	assert.Equal(t, 2, n)
	assert.Nil(t, err)

	assert.Equal(t, byte(0x1), buf.Bytes()[0],
		"amf0/encoder: did not write header byte")
	assert.Equal(t, byte(0x0), buf.Bytes()[1],
		"amf0/encoder: did not write type payload")
}

func TestUnsuccessfulEncodingReturnsError(t *testing.T) {
	buf := new(bytes.Buffer)

	typ := new(MockAmfType)
	typ.On("Marker").Return(byte(0x00))
	typ.On("Encode", mock.Anything).Return(0, errors.New("test"))

	n, err := amf0.Encode(typ, buf)

	assert.Equal(t, 0, n)
	assert.Equal(t, "test", err.Error())
}

func TestEncodeToBytesDoesNotMangleOutput(t *testing.T) {
	s := amf0.NewString("hello world")

	buf := new(bytes.Buffer)

	n, e1 := amf0.Encode(s, buf)
	out, e2 := amf0.EncodeToBytes(s)

	assert.Len(t, out, n)
	assert.Equal(t, buf.Bytes(), out)
	assert.Equal(t, e1, e2)
}
