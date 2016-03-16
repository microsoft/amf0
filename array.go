package amf0

import (
	"bytes"
	"encoding/binary"
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

// Implements AmfType.Marker
func (a *Array) Marker() byte { return MARKER_ECMA_ARRAY }

// Implements AmfType.Decode
func (a *Array) Decode(r io.Reader) error {
	var n [4]byte
	if _, err := r.Read(n[:]); err != nil {
		return err
	}

	a.pairs = make([]*pair, 0, binary.BigEndian.Uint32(n[:]))
	for i := 0; i < cap(a.pairs); i++ {
		if err := a.decodePair(r); err != nil {
			return err
		}
	}

	var endSeq [3]byte
	if _, err := r.Read(endSeq[:]); err != nil {
		return err
	}

	if !bytes.Equal(objectEndSeq, endSeq[:]) {
		return errors.New("amf0: missing end sequence for array")
	}

	return nil
}

// Implements AmfType.Encode
func (a *Array) Encode(w io.Writer) (int, error) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, uint32(len(a.pairs)))
	a.writePairs(buf)
	buf.Write(objectEndSeq)

	n, err := io.Copy(w, buf)
	return int(n), err
}

// Implements AmfType.EncodeBytes
func (a *Array) EncodeBytes() []byte {
	buf := new(bytes.Buffer)
	a.Encode(buf)

	return buf.Bytes()
}
