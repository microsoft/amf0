package amf0

import (
	"bytes"
	"io"
)

type Array struct {
	*paired
}

var _ AmfType = &Array{}

func NewArray() *Array {
	return &Array{newPaired()}
}

// Implements AmfType.Decode
func (a *Array) Decode(r io.Reader) error {
	numBytes, err := readBytes(r, 4)
	if err != nil {
		return err
	}

	num := int(getUint32(numBytes, 0))
	a.pairs = make([]*pair, 0, num)

	for i := 0; i < num; i++ {
		if err := a.decodePair(r); err != nil {
			return err
		}
	}

	return nil
}

// Adds a new pair to the object.
func (a *Array) Add(key string, value AmfType) *Array {
	a.paired.Add(key, value)
	return a
}

// Implements AmfType.Encode
func (a *Array) Encode(w io.Writer) (int, error) {
	wc := newWriteCollector(w)
	wc.Write([]byte{MARKER_ECMA_ARRAY})

	num := make([]byte, 4)
	putUint32(num, 0, uint32(len(a.pairs)))
	wc.Write(num)

	a.writePairs(wc)
	return wc.Totals()
}

// Implements AmfType.EncodeBytes
func (a *Array) EncodeBytes() []byte {
	buf := new(bytes.Buffer)
	a.Encode(buf)

	return buf.Bytes()
}

// Implements AmfType.Marker
func (a *Array) Marker() byte {
	return MARKER_ECMA_ARRAY
}
