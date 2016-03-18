package amf0

import "io"

type Decoder func(r io.Reader) (AmfType, error)

var (
	Decode Decoder = func(r io.Reader) (AmfType, error) {
		var typeId [1]byte
		if _, err := r.Read(typeId[:]); err != nil {
			return nil, err
		}

		packet, err := DefaultIdentifier.Identify(typeId[0])
		if err != nil {
			return nil, err
		}

		if err = packet.Decode(r); err != nil {
			return nil, err
		}

		return packet, nil
	}
)
