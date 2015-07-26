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
