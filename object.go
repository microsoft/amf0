package amf0

import (
	"bytes"
	"errors"
	"io"
)

type Object struct {
	// object pairs are ordered, so we're not using a map.
	pairs    []*objectPair
	encoded  []byte
	uptodate bool
}

type objectPair struct {
	Key   []byte
	Value AmfType
}

var _ AmfType = &Object{}

func NewObject() *Object {
	return &Object{}
}

func newPair(key string, value AmfType) *objectPair {
	return &objectPair{
		Key:   []byte(key),
		Value: value,
	}
}

// This error is returned when trying to get a key that doesn't exist.
var NotFoundError = errors.New("Item not found in the object.")

// This error is returned when a key exists, but its type is not
// the one requested.
var WrongTypeError = errors.New("Item not found in the object.")

// Implements AmfType.Decode
func (o *Object) Decode(r io.Reader) error {
	var pairs []*objectPair
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

		pairs = append(pairs, &objectPair{Key: key, Value: value})

		buf = append(buf, key...)
		buf = append(buf, value.EncodeBytes()...)
	}

	o.pairs = pairs
	o.encoded = buf
	o.uptodate = true
	return nil
}

// Implements AmfType.DecodeFrom
func (o *Object) DecodeFrom(slice []byte, pos int) (int, error) {
	var pairs []*objectPair
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
		pairs = append(pairs, &objectPair{Key: key, Value: value})
	}

	o.pairs = pairs
	o.encoded = slice[start:pos]
	o.uptodate = true
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

func (o *Object) makeEncoded() {
	// todo: see if calculating, allocating, and writing to a
	// single slice is more efficient.

	buf := new(bytes.Buffer)
	keylen := make([]byte, 2)
	for _, pair := range o.pairs {
		putUint16(keylen, 0, uint16(len(pair.Key)))
		buf.Write(keylen)
		buf.Write(pair.Key)
		buf.Write([]byte{pair.Value.Marker()})
		pair.Value.Encode(buf)
	}

	buf.Write([]byte{0x00, 0x00, MARKER_OBJECT_END})
	o.encoded = buf.Bytes()
	o.uptodate = true
}

// Adds a new pair to the object.
func (o *Object) Add(key string, value AmfType) *Object {
	o.pairs = append(o.pairs, newPair(key, value))
	o.uptodate = false
	return o
}

// Returns the number of kv pairs in the object.
func (o *Object) Size() int {
	return len(o.pairs)
}

// Implements AmfType.Encode
func (o *Object) Encode(w io.Writer) (int, error) {
	return w.Write(o.EncodeBytes())
}

// Implements AmfType.EncodeTo
func (o *Object) EncodeTo(slice []byte, pos int) {
	copy(slice[pos:], o.EncodeBytes())
}

// Implements AmfType.EncodeBytes
func (o *Object) EncodeBytes() []byte {
	if !o.uptodate {
		o.makeEncoded()
	}

	return o.encoded
}

// Implements AmfType.Marker
func (o *Object) Marker() byte {
	return MARKER_OBJECT
}
