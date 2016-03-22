package amf0

import "reflect"

// Identifier is a type capable of preforming bi-direcitonal lookups with
// respect to the various AmfTypes implemented here. It has two discrete
// responsibilities:
//	- Fetch the AmfType assosciated with a given marker ID.
//	```
//	typ := DefaultIdentifier.TypeOf(0x01)
//	  #=> (reflect.TypeOf(new(amf0.Bool)).Elem())
//	```
//
//	- Fetch the AmfType assosciated with a given native type. (i.e., going
//	from a `bool` to an `amf0.Bool`.)
//	```
//	typ := DefaultIdentifier.AmfType(true)
//	  #=> (reflect.TypeOf(new(amf0.Bool)).Elem())
//	```
type Identifier struct {
	ids  map[byte]reflect.Type
	typs map[reflect.Type]reflect.Type
}

var (
	// DefaultIdentifier is the default implementation of the Identifier
	// type. It holds knowledge of all implemented amf0 types in this
	// package.
	DefaultIdentifier *Identifier = NewIdentifier(
		new(Array), new(Null), new(Undefined), new(Bool), new(Number),
		new(Object), new(String), new(LongString),
	)
)

// NewIdentifier returns a pointer to a new instance of the Identifier type. By
// calling this method, all of the TypeOf and AmfType permutations are
// precomputed, saving tiem in the future.
func NewIdentifier(types ...AmfType) *Identifier {
	i := &Identifier{
		ids:  make(map[byte]reflect.Type),
		typs: make(map[reflect.Type]reflect.Type),
	}

	for _, t := range types {
		typ := reflect.TypeOf(t).Elem()

		i.ids[t.Marker()] = typ
		i.typs[t.Native()] = typ
	}

	return i
}

// TypeOf returns the AmfType assosciated with a given marker ID.
func (i *Identifier) TypeOf(id byte) reflect.Type {
	return i.ids[id]
}

// AmfType returns the AmfType assosciated with a given primitive.
func (i *Identifier) AmfType(v interface{}) reflect.Type {
	return i.typs[reflect.TypeOf(v)]
}
