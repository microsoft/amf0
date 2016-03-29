package encoding_test

import (
	"bytes"
	"testing"

	"github.com/WatchBeam/amf0"
	"github.com/WatchBeam/amf0/encoding"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshallerConstruction(t *testing.T) {
	s := encoding.NewUnmarshaler(new(bytes.Buffer))

	assert.IsType(t, &encoding.Unmarshaler{}, s)
}

func TestSimpleUnmarshal(t *testing.T) {
	target := struct {
		Foo bool
	}{}

	err := encoding.Unmarshal(bytes.NewBuffer([]byte{
		0x01, 0x01, // <bool>, <true>
	}), &target)

	assert.Nil(t, err)
	assert.Equal(t, true, target.Foo)
}

func TestStructScanWithoutConverstion(t *testing.T) {
	target := struct {
		Foo amf0.Bool
	}{}

	err := encoding.Unmarshal(bytes.NewBuffer([]byte{
		0x01, 0x01, // <bool>, <true>
	}), &target)

	assert.Nil(t, err)
	assert.Equal(t, amf0.Bool(true), target.Foo)
}
