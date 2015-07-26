package amf0

import (
	"encoding/binary"
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func aBunchOfBytes(amount int) []byte {
	b := make([]byte, amount)
	for i := 0; i < amount; i++ {
		b[i] = byte(rand.Intn(0xFF))
	}

	return b
}

type reluctantReader struct {
	src   []byte
	start int
}

func (r *reluctantReader) Read(b []byte) (int, error) {
	if r.start == len(r.src) {
		return 0, io.EOF
	}
	end := r.start + rand.Intn(3) + 2
	if end-r.start > len(b) {
		end = r.start + len(b)
	}
	if end > len(r.src) {
		end = len(r.src)
	}
	copy(b, r.src[r.start:end])
	delta := end - r.start
	r.start = end
	return delta, nil
}

func TestReadBytes(t *testing.T) {
	data := aBunchOfBytes(128)
	r := &reluctantReader{src: data}

	out, err := readBytes(r, 128)

	assert.Nil(t, err)
	assert.Equal(t, data, out)
}

func TestReadBytesTo(t *testing.T) {
	data := aBunchOfBytes(64)
	r := &reluctantReader{src: data}

	container1 := aBunchOfBytes(128)
	container2 := make([]byte, 128)
	copy(container2, container1)
	err := readBytesTo(r, 64, container1, 32)

	assert.Nil(t, err)
	assert.Equal(t, container2[0:31], container1[0:31])
	assert.Equal(t, data, container1[32:96])
	assert.Equal(t, container2[96:], container1[96:])
}

func TestGetVarUint(t *testing.T) {
	b := []byte{0x42}
	assert.Equal(t, uint64(0x42), getVarUint(b, 0, 1))
	b = []byte{0x42, 0x42}
	assert.Equal(t, uint64(binary.BigEndian.Uint16(b)), getVarUint(b, 0, 2))
	b = []byte{0x42, 0x42, 0x42, 0x42}
	assert.Equal(t, uint64(binary.BigEndian.Uint32(b)), getVarUint(b, 0, 4))
	b = []byte{0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42}
	assert.Equal(t, binary.BigEndian.Uint64(b), getVarUint(b, 0, 8))
}

func TestRwReader(t *testing.T) {
	b := make([]byte, 0xFF)
	for i := 0x00; i < 0xFF; i++ {
		b[i] = byte(i)
	}

	r := newRwReader(&reluctantReader{src: b})
	out, _ := readBytes(r, 5)
	assert.Equal(t, []byte{0x00, 0x01, 0x02, 0x03, 0x04}, out)
	r.Rewind([]byte{0x03, 0x04})
	out, _ = readBytes(r, 5)
	assert.Equal(t, []byte{0x03, 0x04, 0x05, 0x06, 0x07}, out)
	out, _ = readBytes(r, 5)
	assert.Equal(t, []byte{0x08, 0x09, 0x0a, 0x0b, 0x0c}, out)

	out = make([]byte, 1)
	r.Rewind([]byte{0x0b, 0x0c})
	n, _ := r.Read(out)
	assert.Equal(t, 1, n)
	assert.Equal(t, byte(0x0b), out[0])
	n, _ = r.Read(out)
	assert.Equal(t, 1, n)
	assert.Equal(t, byte(0x0c), out[0])
	n, _ = r.Read(out)
	assert.Equal(t, 1, n)
	assert.Equal(t, byte(0x0d), out[0])
}
