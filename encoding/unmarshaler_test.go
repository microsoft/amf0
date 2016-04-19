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

func TestUnmarshallingANonNilPointerValue(t *testing.T) {
	target := struct {
		Foo *amf0.Bool
	}{}

	err := encoding.Unmarshal(bytes.NewReader([]byte{
		0x01, 0x01, // <bool>, <true>
	}), &target)

	assert.Nil(t, err)
	assert.EqualValues(t, true, *target.Foo)
}

func TestUnmarshallingANonNilValue(t *testing.T) {
	target := struct {
		Foo amf0.Bool
	}{}

	err := encoding.Unmarshal(bytes.NewReader([]byte{
		0x01, 0x01, // <bool>, <true>
	}), &target)

	assert.Nil(t, err)
	assert.EqualValues(t, true, target.Foo)
}

func TestUnmarshallingANilPointerValue(t *testing.T) {
	target := struct {
		Foo *amf0.Bool
	}{}

	err := encoding.Unmarshal(bytes.NewReader([]byte{
		0x05, // <null>
	}), &target)

	assert.Nil(t, err)
	assert.Nil(t, target.Foo)
}

func TestUnmarshallingANilValue(t *testing.T) {
	target := struct {
		Foo amf0.Bool
	}{}

	err := encoding.Unmarshal(bytes.NewReader([]byte{
		0x05, // <null>
	}), &target)

	assert.Equal(t, "amf0: unable to assign nil to type amf0.Bool",
		err.Error())
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

func TestStrucScanWithNil(t *testing.T) {
	target := struct {
		Foo *amf0.Object
	}{}

	err := encoding.Unmarshal(bytes.NewBuffer([]byte{
		0x05, // <null>
	}), &target)

	assert.Nil(t, err)
	assert.Nil(t, target.Foo)
}
