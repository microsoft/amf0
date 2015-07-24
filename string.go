package amf0

import (
	"io"
)

// Base string is a variable-length string container used for
// both String and Long String types
type baseString struct {
	bytes   []byte
	body    string
	sizeLen int
}

// Implements AmfType.Decode
func (s *baseString) Decode(r io.Reader) error {
	sizeBytes, err := readBytes(r, s.sizeLen)
	if err != nil {
		return err
	}

	size := int(getVarUint(sizeBytes, 0, s.sizeLen))
	s.bytes, err = readBytes(r, size)
	if err != nil {
		return nil
	}

	return nil
}

// Implements AmfType.DecodeFrom
func (s *baseString) DecodeFrom(slice []byte, pos int) (int, error) {
	if len(slice)-pos < s.sizeLen {
		return 0, io.EOF
	}

	size := int(getVarUint(slice, pos, s.sizeLen))
	total := size + s.sizeLen
	if len(slice)-pos < total {
		return 0, io.EOF
	}

	s.bytes = slice[pos+s.sizeLen : pos+total]

	return total, nil
}

// Gets the contents of this message as a byte slice.
func (s *baseString) GetBytes() []byte {
	return s.bytes
}

// Returns the string content of this type.
func (s *baseString) GetBody() string {
	// The body is decoded lazily, since utf
	// decoding is relatively expensive.
	if s.body == "" {
		s.body = string(s.bytes)
	}

	return s.body
}

// Sets the trying content of this type.
func (s *baseString) SetBody(str string) {
	s.body = str
	s.bytes = []byte(str)
}

// Implements AmfType.Encode
func (s *baseString) Encode(w io.Writer) (int, error) {
	return w.Write(s.EncodeBytes())
}

// Implements AmfType.EncodeTo
func (s *baseString) EncodeTo(slice []byte, pos int) {
	copy(slice[pos:], s.EncodeBytes())
}

// Implements AmfType.EncodeBytes
func (s *baseString) EncodeBytes() []byte {
	bytes := make([]byte, 1+s.sizeLen+len(s.bytes))
	bytes[0] = MARKER_STRING
	putVarUint(bytes, 1, uint64(len(s.bytes)), s.sizeLen)
	copy(bytes[1+s.sizeLen:], s.bytes)
	return bytes
}

type String struct{ *baseString }
type LongString struct{ *baseString }

var _ AmfType = &String{}
var _ AmfType = &LongString{}

// Creates a new string type, with an optional initial content.
func NewString(str ...string) *String {
	s := &String{&baseString{sizeLen: 2}}
	if len(str) == 1 {
		s.SetBody(str[0])
	}

	return s
}

// Creates a new long string type, with an optional initial content.
func NewLongString(str ...string) *LongString {
	s := &LongString{&baseString{sizeLen: 4}}
	if len(str) == 1 {
		s.SetBody(str[0])
	}

	return s
}

// Implements AmfType.Marker
func (s *String) Marker() byte {
	return MARKER_STRING
}

// Implements AmfType.Marker
func (s *LongString) Marker() byte {
	return MARKER_LONG_STRING
}
