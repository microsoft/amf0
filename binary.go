package amf0

import (
	"io"
)

// Writes a uint32 at a specific position in the slice.
func putUint32(b []byte, pos int, v uint32) {
	b[pos+0] = byte(v >> 24)
	b[pos+1] = byte(v >> 16)
	b[pos+2] = byte(v >> 8)
	b[pos+3] = byte(v)
}

// Gets a uint32 from a specific position in the slice
func getUint32(b []byte, pos int) uint32 {
	return uint32(b[pos+3]) |
		uint32(b[pos+2])<<8 |
		uint32(b[pos+1])<<16 |
		uint32(b[pos])<<24
}

// Writes a uint32 at a specific position in the slice.
func putUint64(b []byte, pos int, v uint64) {
	b[pos+0] = byte(v >> 56)
	b[pos+1] = byte(v >> 48)
	b[pos+2] = byte(v >> 40)
	b[pos+3] = byte(v >> 32)
	b[pos+4] = byte(v >> 24)
	b[pos+5] = byte(v >> 16)
	b[pos+6] = byte(v >> 8)
	b[pos+7] = byte(v)
}

// Gets a uint64 from a specific position in the slice
func getUint64(b []byte, pos int) uint64 {
	return uint64(b[pos+7]) |
		uint64(b[pos+6])<<8 |
		uint64(b[pos+5])<<16 |
		uint64(b[pos+4])<<24 |
		uint64(b[pos+3])<<32 |
		uint64(b[pos+2])<<40 |
		uint64(b[pos+1])<<48 |
		uint64(b[pos])<<56
}

// Writes a uint32 at a specific position in the slice.
func putUint16(b []byte, pos int, v uint16) {
	b[pos] = byte(v >> 8)
	b[pos+1] = byte(v)
}

// Gets a uint16 from a specific position in the slice
func getUint16(b []byte, pos int) uint16 {
	return uint16(b[pos+1]) | uint16(b[pos])<<8
}

// Reads a variable-length unsigned integer from the byte
// slice, returning it as a uint64.
func getVarUint(b []byte, pos int, size int) uint64 {
	var out uint64

	for i := 0; i < size; i++ {
		out |= uint64(b[pos+i]) << (uint(size-i-1) << 3)
	}

	return out
}

// Writes a variable-length unsigned integer from to byte
// slice starting at the specified position.
func putVarUint(b []byte, pos int, v uint64, size int) {
	for i := 0; i < size; i++ {
		b[pos+i] = byte(v >> (uint(size-i-1) << 3))
	}
}

// Reads `n` bytes from the reader, writing them to the `to` byte
// slice starting at position `start`. It's expected that `to`
// will have the necessary capacity to contain them.
func readBytesTo(r io.Reader, n int, to []byte, start int) error {
	read, err := readBytes(r, n)
	if err != nil {
		return err
	}

	copy(to[start:], read)
	return nil
}

// Reads and returns `n` bytes from the passed reader.
func readBytes(r io.Reader, n int) ([]byte, error) {
	output := make([]byte, 0, n)
	for n > 0 {
		buf := make([]byte, n)
		read, err := r.Read(buf)
		if err != nil {
			return nil, err
		}

		output = append(output, buf[:read]...)
		n -= read
	}

	return output, nil
}
