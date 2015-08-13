package amf0

import (
	"fmt"
	"io"
)

// Provides a high-level reader structure for managing AMF streams.
type Reader struct {
	reader  io.Reader
	lastErr error
}

// Creates a Reader that wraps around an io.Reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{reader: r}
}

// Decodes an AmfType from the reader into the pointer.
func (r *Reader) Decode(into *AmfType) *Reader {
	if r.lastErr != nil {
		return r
	}

	amf, err := Decode(r.reader)
	*into = amf
	r.lastErr = err

	return r
}

func (r *Reader) decode(marker byte, name string, into AmfType) *Reader {
	if r.lastErr != nil {
		return r
	}

	code, err := readBytes(r.reader, 1)
	if err != nil {
		r.lastErr = err
		return r
	}

	if code[0] != marker {
		r.lastErr = fmt.Errorf(
			"Amf0: expected %s (marker 0x%x) but got 0x%x",
			name, marker, code[0],
		)
		return r
	}

	r.lastErr = into.Decode(r.reader)
	return r
}

// Decodes a string from the stream.
func (r *Reader) String(into **String) *Reader {
	target := NewString()
	*into = target

	return r.decode(MARKER_STRING, "String", target)
}

// Decodes a string from the stream.
func (r *Reader) LongString(into **LongString) *Reader {
	target := NewLongString()
	*into = target

	return r.decode(MARKER_LONG_STRING, "LongString", target)
}

// Decodes a string from the stream.
func (r *Reader) Boolean(into **Boolean) *Reader {
	target := NewBoolean()
	*into = target

	return r.decode(MARKER_BOOLEAN, "Boolean", target)
}

// Decodes a string from the stream.
func (r *Reader) Number(into **Number) *Reader {
	target := NewNumber()
	*into = target

	return r.decode(MARKER_NUMBER, "Number", target)
}

// Decodes a string from the stream.
func (r *Reader) Array(into **Array) *Reader {
	target := NewArray()
	*into = target

	return r.decode(MARKER_ECMA_ARRAY, "Array", target)
}

// Decodes a string from the stream.
func (r *Reader) Object(into **Object) *Reader {
	target := NewObject()
	*into = target

	return r.decode(MARKER_OBJECT, "Object", target)
}

// Decodes a string from the stream.
func (r *Reader) Null(into **NullType) *Reader {
	target := NewNull()
	*into = target

	return r.decode(MARKER_NULL, "NullType", target)
}

// Skips the next n atoms in the reader.
func (r *Reader) Skip(n int) *Reader {
	var target AmfType
	for i := 0; i < n; i++ {
		r.Decode(&target)
	}

	return r
}

// Returns the last encountered error in the reader (or nil)
func (r *Reader) Error() error {
	return r.lastErr
}
