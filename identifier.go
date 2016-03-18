package amf0

import (
	"fmt"
	"reflect"
)

type Identifier struct {
	typs map[byte]reflect.Type
}

var (
	DefaultIdentifier *Identifier = NewIdentifier([]AmfType{
		new(Array), new(Null), new(Undefined), new(Bool), new(Number),
		new(Object), new(String), new(LongString),
	})
)

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
