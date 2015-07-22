package amf0

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestScanner(t *testing.T) {
	r := &reluctantReader{src: getEncoded()}
	scanner := NewScanner(r)

	assert.True(t, scanner.Scan())
	assert.Equal(t, "こんにちは", scanner.Type().(*String).GetBody())
	assert.True(t, scanner.Scan())
	assert.Equal(t, float64(42), scanner.Type().(*Number).GetNumber())
	assert.False(t, scanner.Scan())
	assert.Equal(t, io.EOF, scanner.Err())
}
