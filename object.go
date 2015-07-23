package amf0

import (
	"bytes"
	"errors"
	"io"
)

type Object struct {
	// object pairs are ordered, so we're not using a map.
	pairs   []*ObjectPair
	encoded []byte
}

type ObjectPair struct {
	Key   []byte
	Value AmfType
}

var _ AmfType = &Object{}

func NewObject() *Object {
	return &Object{}
}

// This error is returned when trying to get a key that doesn't exist.
var NotFoundError = errors.New("Item not found in the object.")

// This error is returned when a key exists, but its type is not
// the one requested.
var WrongTypeError = errors.New("Item not found in the object.")

// Implements AmfType.Decode
func (o *Object) Decode(r io.Reader) error {
	var pairs []*ObjectPair
	buf := make([]byte, 0)
	str := NewString()
	for {
		if err := str.Decode(r); err != nil {
			return err
		}

		key := str.GetBytes()
		if len(key) == 0 {
			readBytes(r, 1)
			break
		}

		value, err := Decode(r)
		if err != nil {
			return err
		}

		pairs = append(pairs, &ObjectPair{Key: key, Value: value})

		buf = append(buf, key...)
		buf = append(buf, value.EncodeBytes()...)
	}

	o.pairs = pairs
	o.encoded = buf
	return nil
}

// Implements AmfType.DecodeFrom
func (o *Object) DecodeFrom(slice []byte, pos int) (int, error) {
	var pairs []*ObjectPair
	str := NewString()
	start := pos
	for {
		n, err := str.DecodeFrom(slice, pos)
		if err != nil {
			return 0, err
		}

		key := str.GetBytes()
		pos += n
		if len(key) == 0 {
			pos++
			break
		}

		value, n, err := DecodeFrom(slice, pos)
		if err != nil {
			return 0, err
		}

		pos += n
		pairs = append(pairs, &ObjectPair{Key: key, Value: value})
	}

	o.pairs = pairs
	o.encoded = slice[start:pos]
	return pos - start, nil
}

// Returns a string type AMF specified by the key. If the
// key isn't found it returns a NotFoundError. If it is found
// but is of the wrong type, this returns a WrongTypeError.
func (o *Object) String(key string) (*String, error) {
	val, err := o.Get(key)
	if err != nil {
		return nil, err
	}

	if cast, ok := val.(*String); ok {
		return cast, nil
	}

	return nil, WrongTypeError
}

// Returns a boolean type AMF specified by the key. If the
// key isn't found it returns a NotFoundError. If it is found
// but is of the wrong type, this returns a WrongTypeError.
func (o *Object) Boolean(key string) (*Boolean, error) {
	val, err := o.Get(key)
	if err != nil {
		return nil, err
	}

	if cast, ok := val.(*Boolean); ok {
		return cast, nil
	}

	return nil, WrongTypeError
}

// Returns a number type AMF specified by the key. If the
// key isn't found it returns a NotFoundError. If it is found
// but is of the wrong type, this returns a WrongTypeError.
func (o *Object) Number(key string) (*Number, error) {
	val, err := o.Get(key)
	if err != nil {
		return nil, err
	}

	if cast, ok := val.(*Number); ok {
		return cast, nil
	}

	return nil, WrongTypeError
}

// Returns an item specified by the key, or returns a NotFoundError.
func (o *Object) Get(key string) (AmfType, error) {
	kb := []byte(key)
	for _, pair := range o.pairs {
		if bytes.Compare(pair.Key, kb) == 0 {
			return pair.Value, nil
		}
	}

	return nil, NotFoundError
}

// Returns the number of kv pairs in the object.
func (o *Object) Size() int {
	return len(o.pairs)
}

// Implements AmfType.Encode
func (o *Object) Encode(w io.Writer) {
	w.Write(o.encoded)
}

// Implements AmfType.EncodeTo
func (o *Object) EncodeTo(slice []byte, pos int) {
	copy(slice[pos:], o.encoded)
}

// Implements AmfType.EncodeBytes
func (o *Object) EncodeBytes() []byte {
	return o.encoded
}
