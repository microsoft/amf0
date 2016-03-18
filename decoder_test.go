package amf0_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/WatchBeam/amf0"
	"github.com/stretchr/testify/assert"
)

func TestEmptyDecode(t *testing.T) {
	b := []byte{}

	typ, err := amf0.Decode(bytes.NewReader(b))

	assert.Nil(t, typ)
	assert.Equal(t, io.EOF, err)
}

func TestSuccesfulDecode(t *testing.T) {
	b := []byte{0x01, 0xf}

	typ, err := amf0.Decode(bytes.NewReader(b))

	assert.Nil(t, err)
	assert.IsType(t, new(amf0.Bool), typ)
	assert.True(t, bool(*(typ.(*amf0.Bool))))
}

func TestUnidentifiableDecode(t *testing.T) {
	b := []byte{0xff}

	typ, err := amf0.Decode(bytes.NewReader(b))

	assert.Nil(t, typ)
	assert.Equal(t, "amf0: unknown packet identifier for 0xff", err.Error())
}

func TestErrorfulDecode(t *testing.T) {
	b := []byte{0x02}

	typ, err := amf0.Decode(bytes.NewReader(b))

	assert.Nil(t, typ)
	assert.Equal(t, io.EOF, err)
}
