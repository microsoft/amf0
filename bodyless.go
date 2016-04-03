package amf0

import (
	"io"
	"reflect"
)

type (
	Null      struct{}
	Undefined struct{}

	Bodyless interface {
		IsBodyless() bool
	}
)

var (
	_ AmfType = new(Null)
	_ AmfType = new(Undefined)
)

func (n *Null) Decode(r io.Reader) error        { return nil }
func (n *Null) Encode(w io.Writer) (int, error) { return 0, nil }
func (n *Null) Marker() byte                    { return 0x05 }
func (n *Null) Native() reflect.Type            { return reflect.TypeOf(n).Elem() }
func (n *Null) IsBodyless() bool                { return true }

func (u *Undefined) Decode(r io.Reader) error        { return nil }
func (u *Undefined) Encode(w io.Writer) (int, error) { return 0, nil }
func (u *Undefined) Marker() byte                    { return 0x06 }
func (u *Undefined) Native() reflect.Type            { return reflect.TypeOf(u).Elem() }
func (u *Undefined) IsBodyless() bool                { return true }
