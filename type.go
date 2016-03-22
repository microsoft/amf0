package amf0

import (
	"io"
	"reflect"
)

// AMF Types (strings, numbers, etc) implement the AmfType
// interface, which specifies that they have methods available
// to decode and encode data.
type AmfType interface {
	// Decodes information for the type from the reader. This expects
	// the reader to return starting from the first byte _after_ the
	// type marker.
	Decode(io.Reader) error

	// Encodes and writes the type to the reader. Returns an error
	// if one occurred on the reader.
	Encode(io.Writer) (int, error)

	// Gets the associated marker byte for the type.
	Marker() byte

	// Native returns the native Golang type assosciated with this AmfType.
	Native() reflect.Type
}
