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
	ids  map[byte]TypeFactory
	typs map[reflect.Type]TypeFactory
}

// TypeFactory is a factory type that returns new instances of a specific
// AmfType.
type TypeFactory func() AmfType

var (
	// DefaultIdentifier is the default implementation of the Identifier
	// type. It holds knowledge of all implemented amf0 types in this
	// package.
	DefaultIdentifier = NewIdentifier(
		func() AmfType { return &Array{NewPaired()} },
		func() AmfType { return new(Null) },
		func() AmfType { return new(Undefined) },
		func() AmfType { return new(Bool) },
		func() AmfType { return new(Number) },
		func() AmfType { return &Object{NewPaired()} },
		func() AmfType { return new(String) },
		func() AmfType { return new(LongString) },
	)
)

// NewIdentifier returns a pointer to a new instance of the Identifier type. By
// calling this method, all of the TypeOf and AmfType permutations are
// precomputed, saving tiem in the future.
func NewIdentifier(types ...TypeFactory) *Identifier {
	i := &Identifier{
		ids:  make(map[byte]TypeFactory),
		typs: make(map[reflect.Type]TypeFactory),
	}

	for _, f := range types {
		t := f()

		i.ids[t.Marker()] = f
		i.typs[t.Native()] = f
	}

	return i
}

// TypeOf returns the AmfType assosciated with a given marker ID.
func (i *Identifier) TypeOf(id byte) AmfType {
	f := i.ids[id]
	if f == nil {
		return nil
	}

	return f()
}

// NewMatchingTypeFromValue returns a new instance of an AmfType in the
// same kind as given by v. If no matching type is found,
// nil is returned instead.
func (i *Identifier) NewMatchingTypeFromValue(val reflect.Value) AmfType {
	if f := i.typs[val.Type()]; f != nil {
		return f()
	}

	if val.Kind() != reflect.Ptr {
		return i.NewMatchingTypeFromValue(val.Addr())
	}

	return nil
}

// NewMatchingType returns a new instance of an AmfType in the same kind as
// given by v. If no matching type is found, nil is returned instead.
func (i *Identifier) NewMatchingType(v interface{}) AmfType {
	return i.NewMatchingType(reflect.ValueOf(v))
}
