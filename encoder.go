package amf0

import (
	"bytes"
	"io"
)

// Encoder represents a func capable of writing a representation of the AmfType
// "t" to the io.Writer "w". By contract, this type must obey the rules of amf0
// encoding, which means there must be a marker, and a payload.
type Encoder func(t AmfType, w io.Writer) (int, error)

var (
	// Encode serves as a default implementation of the Encoder type. It is
	// compliant with the amf0 specification. It returns the number of bytes
	// that it wrote (usually 1+len(payload)), and any errors that it
	// encountered along the way.
	Encode Encoder = func(t AmfType, w io.Writer) (int, error) {
		buf := new(bytes.Buffer)

		buf.Write([]byte{t.Marker()})
		if _, err := t.Encode(buf); err != nil {
			return 0, err
		}

		n, err := io.Copy(w, buf)
		return int(n), err
	}
)

// EncodeToBytes uses the default Encoder to marshal the given AmfType `t` to a
// []byte, instead of writing to an io.Writer. Any error returned above will be
// returned here.
func EncodeToBytes(t AmfType) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := Encode(t, buf)

	return buf.Bytes(), err
}
