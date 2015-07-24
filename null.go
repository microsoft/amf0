package amf0

type NullType struct {
	bodylessType
}

var _ AmfType = &NullType{}

func NewNull() *NullType {
	return &NullType{bodylessType{marker: MARKER_NULL}}
}
