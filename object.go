package amf0

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
)

var (
	ObjectEndSeq = []byte{0x00, 0x00, 0x09}
)

type Object struct {
	*Paired
}

var _ AmfType = new(Object)

func NewObject() *Object {
	return &Object{NewPaired()}
}

// Implements AmfType.Marker
func (o *Object) Marker() byte { return 0x03 }

// Implements AmfType.Native
func (o *Object) Native() reflect.Type { return reflect.TypeOf(o) }

// Implements AmfType.Decode
func (o *Object) Decode(r io.Reader) error {
	br := bufio.NewReader(r)

	for {
		cont, err := objectShouldContinue(br)
		if err != nil {
			return err
		}
		if !cont {
			if _, err = br.Discard(3); err != nil {
				return err
			}
			break
		}

		if err := o.decodePair(br); err != nil {
			return err
		}
	}

	return nil
}

// Implements AmfType.Encode
func (o *Object) Encode(w io.Writer) (int, error) {
	buf := new(bytes.Buffer)

	o.writePairs(buf)
	buf.Write(ObjectEndSeq)

	n, err := io.Copy(w, buf)
	return int(n), err
}
func objectShouldContinue(r *bufio.Reader) (bool, error) {
	b, err := r.Peek(len(ObjectEndSeq))
	if err != nil {
		return false, err
	}

	return !bytes.Equal(ObjectEndSeq, b), nil
}
