package amf0

import (
	"fmt"
	"io"
	"reflect"
)

// Decodes a single packet from an io.Reader
func Decode(r io.Reader) (AmfType, error) {
	token, err := readBytes(r, 1)
	if err != nil {
		return nil, err
	}

	packet, err := DefaultIdentifier.Identify(token[0])
	if err != nil {
		return nil, err
	}

	return packet, packet.Decode(r)
}

var DefaultIdentifier *Identifier = NewIdentifier([]AmfType{
	new(Bool), new(Number), new(String), new(Number),
})

type Identifier struct {
	typs map[byte]reflect.Type
}

func NewIdentifier(types []AmfType) *Identifier {
	i := &Identifier{
		typs: make(map[byte]reflect.Type),
	}

	for _, t := range types {
		i.typs[t.Marker()] = reflect.TypeOf(t).Elem()
	}

	return i
}

func (i *Identifier) Identify(id byte) (AmfType, error) {
	typ := i.typs[id]
	if typ == nil {
		return nil, UnknownPacketError(id)
	}

	v := reflect.New(typ).Interface().(AmfType)
	return v, nil
}

type UnknownPacketError byte

func (e UnknownPacketError) Error() string {
	return fmt.Sprintf("Unknown packet identifier for %d", byte(e))
}
