package amf0

import (
	"bytes"
	"errors"
	"io"
)

// The Paired interface can be used for accessing array and object types.
type Paired interface {
	AmfType
	Add(key string, value AmfType)
	String(key string) (*String, error)
	Bool(key string) (*Bool, error)
	Number(key string) (*Number, error)
	Get(key string) (AmfType, error)
	Size() int
}

// The Paired structure is used to parse types which have key/value
// pairs, such as objects and arrays.
type paired struct {
	pairs []*pair
}

type pair struct {
	Key   []byte
	Value AmfType
}

func newPaired() *paired {
	return &paired{pairs: make([]*pair, 0)}
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

// Decodes a single kv pair from the array
func (p *paired) decodePair(r io.Reader) error {
	str := NewString()

	if err := str.Decode(r); err != nil {
		return err
	}

	value, err := Decode(r)
	if err != nil {
		return err
	}

	p.pairs = append(p.pairs, &pair{
		Key:   str.GetBytes(),
		Value: value,
	})

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
func (p *paired) Bool(key string) (*Bool, error) {
	val, err := p.Get(key)
	if err != nil {
		return nil, err
	}

	if cast, ok := val.(*Bool); ok {
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
