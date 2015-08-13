package amf0

import (
	"bytes"
	"errors"
	"io"
)

type Array struct {
	*paired
}

var _ AmfType = &Array{}
var _ Paired = &Array{}

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

	ender, err := readBytes(r, 3)
	if err != nil {
		return err
	}

	if bytes.Compare(ender, objectEndSeq) != 0 {
		return errors.New("amf0: missing end sequence for array")
	}

	return nil
}

// Implements AmfType.Encode
func (a *Array) Encode(w io.Writer) (int, error) {
	wc := newWriteCollector(w)
	wc.Write([]byte{MARKER_ECMA_ARRAY})

	num := make([]byte, 4)
	putUint32(num, 0, uint32(len(a.pairs)))
	wc.Write(num)
	a.writePairs(wc)
	wc.Write(objectEndSeq)

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
