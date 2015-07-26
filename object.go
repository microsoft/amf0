package amf0

import (
	"bytes"
	"io"
)

type Object struct {
	*paired
}

type objectPair struct {
	Key   []byte
	Value AmfType
}

var _ AmfType = &Object{}

func NewObject() *Object {
	return &Object{newPaired(objectKeyDecode)}
}

// Adds a new pair to the object.
func (o *Object) Add(key string, value AmfType) *Object {
	o.paired.Add(key, value)
	return o
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

var objectEndSeq = []byte{0x00, 0x00, MARKER_OBJECT_END}

func objectKeyDecode(r rewindingReader) (bool, error) {
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
