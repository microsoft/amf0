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
	s.encoded = make([]byte, s.sizeLen+size)
	if err := readBytesTo(r, size, s.encoded, s.sizeLen); err != nil {
		return err
	}

	copy(s.encoded, sizeBytes)
	s.body = string(s.encoded[s.sizeLen:])

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
	s.body = string(slice[pos+s.sizeLen : pos+total])

	return total, nil
}

// Returns the string content of this type.
func (s *baseString) GetBody() string {
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
func (s *baseString) Encode(w io.Writer) {
	w.Write(s.encoded)
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

func NewString() *String {
	return &String{&baseString{sizeLen: 2}}
}

func NewLongString() *LongString {
	return &LongString{&baseString{sizeLen: 4}}
}
