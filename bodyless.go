package amf0

import "io"

type (
	Null      struct{}
	Undefined struct{}
)

var (
	_ AmfType = new(Null)
	_ AmfType = new(Undefined)
)

func (n *Null) Decode(r io.Reader) error        { return nil }
func (n *Null) Encode(w io.Writer) (int, error) { return 0, nil }
func (n *Null) Marker() byte                    { return 0x05 }

func (u *Undefined) Decode(r io.Reader) error        { return nil }
func (u *Undefined) Encode(w io.Writer) (int, error) { return 0, nil }
func (u *Undefined) Marker() byte                    { return 0x06 }
