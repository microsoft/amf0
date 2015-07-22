package amf0

import (
	"io"
)

// AMF Types (strings, numbers, etc) implement the AmfType
// interface, which specifies that they have methods available
// to decode and encode data.
type AmfType interface {
	// Decodes information for the type from the reader. This expects
	// the reader to return starting from the first byte _after_ the
	// type marker.
	Decode(io.Reader) error
	// Decodes information for the type from the byte slice starting
	// at position `pos`. The position should point to the first
	// byte _after_ the type marker.
	DecodeFrom(slice []byte, pos int) error
	// Encodes and writes the type to the reader.
	Encode(io.Writer)
	// Encodes and writes the type to the position in the slice.
	EncodeTo(slice []byte, pos int)
	// Encodes and returns a byte slice representing the type.
	EncodeBytes() []byte
}
