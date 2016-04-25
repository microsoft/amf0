package amf0

import (
	"fmt"
	"io"
)

type UnknownPacketError byte

func (e UnknownPacketError) Error() string {
	return fmt.Sprintf("amf0: unknown packet identifier for 0x%x", byte(e))
}

type Decoder func(r io.Reader) (AmfType, error)

var (
	Decode Decoder = func(r io.Reader) (AmfType, error) {
		var typeId [1]byte
		if _, err := io.ReadFull(r, typeId[:]); err != nil {
			return nil, err
		}

		typ := DefaultIdentifier.TypeOf(typeId[0])
		if typ == nil {
			return nil, UnknownPacketError(typeId[0])
		}

		if err := typ.Decode(r); err != nil {
			return nil, err
		}

		return typ, nil
	}
)
