package amf0

import (
	"bytes"
	"io"
)

type Object struct {
	*paired
}

var _ AmfType = &Object{}
var _ Paired = &Object{}

func NewObject() *Object {
	return &Object{newPaired()}
}

var objectEndSeq = []byte{0x00, 0x00, MARKER_OBJECT_END}

// Implements AmfType.Decode
func (o *Object) Decode(r io.Reader) error {
	rwr := newRwReader(r)
	for {
		cont, err := objectShouldContinue(rwr)
		if err != nil {
			return err
		}
		if !cont {
			break
		}

		if err := o.decodePair(rwr); err != nil {
			return err
		}
	}

	return nil
}

// Implements AmfType.Encode
func (o *Object) Encode(w io.Writer) (int, error) {
	wc := newWriteCollector(w)
	wc.Write([]byte{MARKER_OBJECT})
	o.writePairs(wc)
	wc.Write(objectEndSeq)
	return wc.Totals()
}

// Implements AmfType.EncodeBytes
func (o *Object) EncodeBytes() []byte {
	buf := new(bytes.Buffer)
	o.Encode(buf)

	return buf.Bytes()
}

// Implements AmfType.Marker
func (o *Object) Marker() byte {
	return MARKER_OBJECT
}

func objectShouldContinue(r rewindingReader) (bool, error) {
	b, err := readBytes(r, 3)
	if err != nil {
		return false, err
	}

	if bytes.Compare(objectEndSeq, b) == 0 {
		return false, nil
	}

	r.Rewind(b)
	return true, nil
}
