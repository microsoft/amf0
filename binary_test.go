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
