package amf0

import (
	"errors"
	"fmt"
	"io"
)

// Decodes a single packet from an io.Reader
func Decode(r io.Reader) (AmfType, error) {
	token, err := readBytes(r, 1)
	if err != nil {
		return nil, err
	}

	packet, err := Identify(token[0])
	if err != nil {
		return nil, err
	}

	return packet, packet.Decode(r)
}

// Returns a packet decoded from the position in the byte string.
// Returns the decoded type and the number of bytes consumed.
func DecodeFrom(b []byte, pos int) (AmfType, int, error) {
	l := len(b)
	if pos >= l {
		return nil, 0, io.EOF
	}

	packet, err := Identify(b[pos])
	if err != nil {
		return nil, 0, err
	}

	n, err := packet.DecodeFrom(b, pos+1)
	return packet, n + 1, err
}

// Identifies the packet type given by `b`, returning a new
// AmfType if recognized, or error otherwise.
func Identify(b byte) (AmfType, error) {
	if int(b) >= len(table) {
		return nil, genUnknownErrorFor(b)
	}
	fn := table[b]
	if fn == nil {
		return nil, genUnknownErrorFor(b)
	}

	return fn(), nil
}

// Generates an error for an invalid packet market b.
func genUnknownErrorFor(b byte) error {
	return errors.New(fmt.Sprintf("Unknown packet identified by %d", b))
}
