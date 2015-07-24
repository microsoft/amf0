package amf0

import (
	"io"
)

// Base string is a variable-length string container used for
// both String and Long String types
type baseString struct {
	encoded []byte
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
	body, err := readBytes(r, size)
	if err != nil {
		return nil
	}

	s.encoded = append(sizeBytes, body...)

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

	s.encoded = slice[pos : pos+total]

	return total, nil
}

// Gets the contents of this message as a byte slice.
func (s *baseString) GetBytes() []byte {
	return s.encoded[s.sizeLen:]
}

// Returns the string content of this type.
func (s *baseString) GetBody() string {
	// The body is decoded lazily, since utf
	// decoding is relatively expensive.
	if s.body == "" {
		s.body = string(s.GetBytes())
	}

	return s.body
}

// Sets the trying content of this type.
func (s *baseString) SetBody(str string) {
	s.body = str

	bytes := []byte(str)
	l := len(bytes)

	s.encoded = make([]byte, l+s.sizeLen)
	putVarUint(s.encoded, 0, uint64(l), s.sizeLen)
	copy(s.encoded[s.sizeLen:], bytes)
}

// Implements AmfType.Encode
func (s *baseString) Encode(w io.Writer) (int, error) {
	return w.Write(s.encoded)
}

// Implements AmfType.EncodeTo
func (s *baseString) EncodeTo(slice []byte, pos int) {
	copy(slice[pos:], s.encoded)
}

// Implements AmfType.EncodeBytes
func (s *baseString) EncodeBytes() []byte {
	return s.encoded
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
