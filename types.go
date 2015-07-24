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
	// byte _after_ the type marker. Returns the number of bytes
	// consumed, or an error.
	DecodeFrom(slice []byte, pos int) (int, error)
	// Encodes and writes the type to the reader.
	Encode(io.Writer)
	// Encodes and writes the type to the position in the slice.
	EncodeTo(slice []byte, pos int)
	// Encodes and returns a byte slice representing the type.
	EncodeBytes() []byte
	// Gets the associated marker byte for the type.
	Marker() byte
}

const (
	MARKER_NUMBER       byte = 0x00
	MARKER_BOOLEAN           = 0x01
	MARKER_STRING            = 0x02
	MARKER_OBJECT            = 0x03
	MARKER_MOVIE_CLIP        = 0x03
	MARKER_NULL              = 0x05
	MARKER_UNDEFINED         = 0x06
	MARKER_REFERENCE         = 0x07
	MARKER_ECMA_ARRAY        = 0x08
	MARKER_OBJECT_END        = 0x09
	MARKER_STRICT_ARRAY      = 0x0A
	MARKER_DATE              = 0x0B
	MARKER_LONG_STRING       = 0x0C
	MARKER_UNSUPPORTED       = 0x0D
	MARKER_RECORD_SET        = 0x0E
	MARKER_XML_DOCUMENT      = 0x0F
	MARKER_TYPED_OBJECT      = 0x10
)

type factory func() AmfType

// Create basically a lookup table for markers. Past four or so
// types, this becomes much more efficient that a switch statement.
var table = [...]factory{
	MARKER_NUMBER: func() AmfType {
		return NewNumber()
	},
	MARKER_BOOLEAN: func() AmfType {
		return NewBoolean()
	},
	MARKER_STRING: func() AmfType {
		return NewString()
	},
	MARKER_LONG_STRING: func() AmfType {
		return NewLongString()
	},
	MARKER_OBJECT: func() AmfType {
		return NewObject()
	},
	MARKER_NULL: func() AmfType {
		return NewNull()
	},
}
