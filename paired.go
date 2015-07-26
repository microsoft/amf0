package amf0

import (
	"bytes"
	"errors"
	"io"
)

// Decode will be called each time a key is read. It should return
// whether to keep going. If it reads data from the reader, then
// it should rewind it back to its original position.
type pairedDecode func(r rewindingReader) (bool, error)

// The Paired structure is used to parse types which have key/value
// pairs, such as objects and arrays.
type paired struct {
	decode pairedDecode
	pairs  []*pair
}

type pair struct {
	Key   []byte
	Value AmfType
}

func newPaired(fn pairedDecode) *paired {
	return &paired{
		decode: fn,
		pairs:  make([]*pair, 0),
	}
}

func newPair(key []byte, value AmfType) *pair {
	return &pair{
		Key:   key,
		Value: value,
	}
}

// This error is returned when trying to get a key that doesn't exist.
var NotFoundError = errors.New("Item not found in the object.")

// This error is returned when a key exists, but its type is not
// the one requested.
var WrongTypeError = errors.New("Item not found in the object.")

// Implements AmfType.Decode
func (p *paired) Decode(r io.Reader) error {
	var pairs []*pair
	str := NewString()
	rwr := newRwReader(r)
	for {
		cont, err := p.decode(rwr)
		if err != nil {
			return err
		}
		if !cont {
			break
		}

		if err := str.Decode(rwr); err != nil {
			return err
		}

		value, err := Decode(rwr)
		if err != nil {
			return err
		}

		pairs = append(pairs, &pair{Key: str.GetBytes(), Value: value})
	}

	p.pairs = pairs
	return nil
}

// Returns a string type AMF specified by the key. If the
// key isn't found it returns a NotFoundError. If it is found
// but is of the wrong type, this returns a WrongTypeError.
func (p *paired) String(key string) (*String, error) {
	val, err := p.Get(key)
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
func (p *paired) Boolean(key string) (*Boolean, error) {
	val, err := p.Get(key)
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
func (p *paired) Number(key string) (*Number, error) {
	val, err := p.Get(key)
	if err != nil {
		return nil, err
	}

	if cast, ok := val.(*Number); ok {
		return cast, nil
	}

	return nil, WrongTypeError
}

// Returns an item specified by the key, or returns a NotFoundError.
func (p *paired) Get(key string) (AmfType, error) {
	kb := []byte(key)
	for _, pair := range p.pairs {
		if bytes.Compare(pair.Key, kb) == 0 {
			return pair.Value, nil
		}
	}

	return nil, NotFoundError
}

// Adds a new pair to the object.
func (p *paired) Add(key string, value AmfType) {
	p.pairs = append(p.pairs, newPair([]byte(key), value))
}

// Returns the number of kv pairs in the object.
func (p *paired) Size() int {
	return len(p.pairs)
}

// Writes out the key pairs.
func (p *paired) writePairs(wc *writeCollector) {
	keylen := make([]byte, 2)
	for _, pair := range p.pairs {
		putUint16(keylen, 0, uint16(len(pair.Key)))
		wc.Write(keylen)
		wc.Write(pair.Key)
		pair.Value.Encode(wc)
	}
}
