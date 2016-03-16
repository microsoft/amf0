package amf0

import (
	"fmt"
	"io"
)

// Provides a high-level reader structure for managing AMF streams.
// Example usage:
//
//      var myNum *Number
//      var myObj *Object
//      err := NewReader(myReader).
//          Skip(2).
//          Number(&myNum).
//          Object(&myObj).
//          Error()
//
//      if err != nil {
//          return err
//      }
//
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
