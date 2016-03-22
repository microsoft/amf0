package amf0_test

import (
	"reflect"
	"testing"

	"github.com/WatchBeam/amf0"
	"github.com/stretchr/testify/assert"
)

func TestIdentifierConstruction(t *testing.T) {
	i := amf0.NewIdentifier()

	assert.IsType(t, &amf0.Identifier{}, i)
}

func TestFetchingKnownPacketTypes(t *testing.T) {
	i := amf0.NewIdentifier(new(amf0.Bool))

	typ := i.TypeOf(0x01)

	assert.Equal(t, reflect.TypeOf(new(amf0.Bool)).Elem(), typ)
}

func TestFetchingUnknownPacketTypes(t *testing.T) {
	i := amf0.NewIdentifier()

	typ := i.TypeOf(0x01)

	assert.Nil(t, typ)
}

func TestFetchingKnownPacketNatives(t *testing.T) {
	i := amf0.NewIdentifier(new(amf0.Bool))

	typ := i.AmfType(false)

	assert.Equal(t, reflect.TypeOf(new(amf0.Bool)).Elem(), typ)
}

func TestFetchingUnknownPacketNatives(t *testing.T) {
	i := amf0.NewIdentifier()

	typ := i.AmfType(false)

	assert.Nil(t, typ)
}
