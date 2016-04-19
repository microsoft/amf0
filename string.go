package amf0

import (
	"bytes"
	"io"
	"reflect"
)

type (
	String     string
	LongString string
)

var (
	_ AmfType = new(String)
	_ AmfType = new(LongString)
)

func NewString(str string) *String {
	s := new(String)
	*s = String(str)

	return s
}

func NewLongString(str string) *LongString {
	s := new(LongString)
	*s = LongString(str)

	return s
}

func (s *String) Marker() byte         { return 0x02 }
func (s *String) Native() reflect.Type { return reflect.TypeOf("") }
func (s *String) Decode(r io.Reader) error {
	if v, err := strDecode(r, 2); err != nil {
		return err
	} else {
		*s = String(v)
	}

	return nil
}
func (s *String) Encode(w io.Writer) (int, error) {
	return strEncode(string(*s), w, 2)
}

func (l *LongString) Marker() byte         { return 0x0c }
func (l *LongString) Native() reflect.Type { return reflect.TypeOf(l) }
func (l *LongString) Decode(r io.Reader) error {
	if v, err := strDecode(r, 4); err != nil {
		return err
	} else {
		*l = LongString(v)
	}

	return nil
}
func (l *LongString) Encode(w io.Writer) (int, error) {
	return strEncode(string(*l), w, 4)
}

func strDecode(r io.Reader, sizeLen int) (string, error) {
	sizeBytes := make([]byte, sizeLen)
	if _, err := io.ReadFull(r, sizeBytes); err != nil {
		return "", err
	}

	var slen uint64
	for i := 0; i < sizeLen; i++ {
		slen |= uint64(sizeBytes[i]) << (uint(sizeLen-i-1) << 3)
	}

	str := make([]byte, int(slen))
	if _, err := io.ReadFull(r, str); err != nil {
		return "", err
	}

	return string(str), nil
}

func strEncode(str string, w io.Writer, sizeLen int) (int, error) {
	b := new(bytes.Buffer)
	strLen := len(str)
	for i := 0; i < sizeLen; i++ {
		b.WriteByte(byte(strLen >> (uint(sizeLen-i-1) << 3)))
	}
	b.WriteString(str)

	if n, err := io.Copy(w, b); err != nil {
		return int(n), err
	}

	return strLen + sizeLen, nil
}
