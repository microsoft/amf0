package amf0_test

import (
	"testing"

	"github.com/WatchBeam/amf0"
	"github.com/stretchr/testify/assert"
)

var (
	BoolFactory = func() amf0.AmfType { return new(amf0.Bool) }
)

func TestIdentifierConstruction(t *testing.T) {
	i := amf0.NewIdentifier()

	assert.IsType(t, &amf0.Identifier{}, i)
}

func TestFetchingKnownPacketTypes(t *testing.T) {
	i := amf0.NewIdentifier(BoolFactory)

	v := i.TypeOf(0x01)

	assert.IsType(t, new(amf0.Bool), v)
}

func TestFetchingUnknownPacketTypes(t *testing.T) {
	i := amf0.NewIdentifier()

	typ := i.TypeOf(0x01)

	assert.Nil(t, typ)
}
