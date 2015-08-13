package amf0

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestReaderGenericDecode(t *testing.T) {
	reader := NewReader(&reluctantReader{src: getEncoded()})

	var second AmfType

	reader.Skip(1)
	assert.Nil(t, reader.Error())
	reader.Decode(&second)
	assert.Nil(t, reader.Error())
	assert.Equal(t, float64(42), second.(*Number).GetNumber())
	reader.Skip(1)
	assert.Equal(t, io.EOF, reader.Error())
}

func TestReaderAssertions(t *testing.T) {
	reader := NewReader(&reluctantReader{src: getEncoded()})

	var first *String
	var second *String

	reader.String(&first)
	assert.Nil(t, reader.Error())
	reader.String(&second)
	assert.Equal(t, "Amf0: expected String (marker 0x2) but got 0x0", reader.Error().Error())
}
