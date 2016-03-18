package amf0

import (
	"bytes"
	"errors"
	"io"
)

var (
	// NotFoundError is returned when trying to get a key that doesn't
	// exist.
	NotFoundError = errors.New("Item not found in the object.")

	// WrongTypeError is returned when a key exists, but its type is not the
	// one requested.
	WrongTypeError = errors.New("Item not found in the object.")
)

type Paired struct {
	tuples []*tuple
}

func NewPaired() *Paired {
	return &Paired{tuples: make([]*tuple, 0)}
}

// Returns the number of kv pairs in the object.
func (p *Paired) Len() int { return len(p.tuples) }

// Adds a new pair to the object.
func (p *Paired) Add(key string, value AmfType) {
	p.tuples = append(p.tuples, &tuple{
		Key:   []byte(key),
		Value: value,
	})
}

// Returns a string type AMF specified by the key. If the
// key isn't found it returns a NotFoundError. If it is found
// but is of the wrong type, this returns a WrongTypeError.
func (p *Paired) String(key string) (*String, error) {
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
func (p *Paired) Bool(key string) (*Bool, error) {
	val, err := p.Get(key)
	if err != nil {
		return nil, err
	}

	if cast, ok := val.(*Bool); ok {
		return cast, nil
	}

	return nil, WrongTypeError
}

// Returns an item specified by the key, or returns a NotFoundError.
func (p *Paired) Get(key string) (AmfType, error) {
	kb := []byte(key)
	for _, tuple := range p.tuples {
		if bytes.Compare(tuple.Key, kb) == 0 {
			return tuple.Value, nil
		}
	}

	return nil, NotFoundError
}

// Decodes a single kv pair from the array
func (p *Paired) decodePair(r io.Reader) error {
	str := new(String)
	if err := str.Decode(r); err != nil {
		return err
	}

	value, err := Decode(r)
	if err != nil {
		return err
	}

	p.tuples = append(p.tuples, &tuple{
		Key:   []byte(string(*str)),
		Value: value,
	})
	return nil
}

// Writes out the key pairs.
func (p *Paired) writePairs(w io.Writer) (int, error) {
	buf := new(bytes.Buffer)

	for _, tuple := range p.tuples {
		tuple.Encode(buf)
	}

	n, err := io.Copy(w, buf)
	return int(n), err
}

type tuple struct {
	Key   []byte
	Value AmfType
}

func (t *tuple) Encode(w io.Writer) (int, error) {
	buf := new(bytes.Buffer)

	if _, err := NewString(string(t.Key)).Encode(buf); err != nil {
		return 0, err
	}
	if _, err := buf.Write([]byte{t.Value.Marker()}); err != nil {
		return 0, err
	}
	if _, err := t.Value.Encode(buf); err != nil {
		return 0, err
	}

	n, err := io.Copy(w, buf)
	return int(n), err
}
