package amf0

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberEncoding(t *testing.T) {
	s := Number(float64(0x4242))

	buf := new(bytes.Buffer)
	n, err := s.Encode(buf)

	assert.Equal(t, 8, n)
	assert.Nil(t, err)

	assert.Equal(t, []byte{
		64, 208, 144, 128, 0, 0, 0, 0,
	}, buf.Bytes())
}

func TestNumberDecodes(t *testing.T) {
	o := new(Number)
	err := o.Decode(bytes.NewReader([]byte{
		64, 208, 144, 128, 0, 0, 0, 0,
	}))

	assert.Nil(t, err)
	assert.Equal(t, float64(0x4242), float64(*o))
}
