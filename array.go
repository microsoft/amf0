package amf0

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

type Array struct {
	*Paired
}

var _ AmfType = new(Array)

func NewArray() *Array {
	return &Array{NewPaired()}
}

// Implements AmfType.Marker
func (a *Array) Marker() byte { return 0x08 }

func (a *Array) Native() reflect.Type { return reflect.TypeOf(a) }

// Implements AmfType.Decode
func (a *Array) Decode(r io.Reader) error {
	var n [4]byte
	if _, err := io.ReadFull(r, n[:]); err != nil {
		return err
	}

	a.tuples = make([]*tuple, 0, binary.BigEndian.Uint32(n[:]))
	for i := 0; i < cap(a.tuples); i++ {
		if err := a.decodePair(r); err != nil {
			return err
		}
	}

	var endSeq [3]byte
	if _, err := io.ReadFull(r, endSeq[:]); err != nil {
		return err
	}

	if !bytes.Equal(ObjectEndSeq, endSeq[:]) {
		return errors.New("amf0: missing end sequence for array")
	}

	return nil
}

// Implements AmfType.Encode
func (a *Array) Encode(w io.Writer) (int, error) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, uint32(len(a.tuples)))
	a.writePairs(buf)
	buf.Write(ObjectEndSeq)

	n, err := io.Copy(w, buf)
	return int(n), err
}

// Implements AmfType.EncodeBytes
func (a *Array) EncodeBytes() []byte {
	buf := new(bytes.Buffer)
	a.Encode(buf)

	return buf.Bytes()
}
