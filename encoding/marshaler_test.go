package encoding_test

import (
	"testing"

	"github.com/WatchBeam/amf0"
	"github.com/WatchBeam/amf0/encoding"
	"github.com/stretchr/testify/assert"
)

func TestMarshalerConstruction(t *testing.T) {
	m := encoding.NewMarshaler()

	assert.IsType(t, &encoding.Marshaler{}, m)
}

func TestMarshalingNativeMembers(t *testing.T) {
	buf, err := encoding.Marshal(&struct {
		Foo bool
	}{false})

	assert.Nil(t, err)

	assert.Equal(t, byte(0x01), buf[0],
		"amf0/encoding: did not marshal	marker")
	assert.Equal(t, byte(0x00), buf[1],
		"amf0/encoding: did not marshal value")
}

func TestMarshallingNonNativeMembers(t *testing.T) {
	buf, err := encoding.Marshal(&struct {
		Foo amf0.Object
	}{*amf0.NewObject()})

	assert.Nil(t, err)
	assert.Equal(t, []byte{0x3, 0x0, 0x0, 0x9}, buf)
}
