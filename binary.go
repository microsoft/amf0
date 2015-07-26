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

// The write collector wraps around a writer. It allowed for
// multiple writes without having to check errors and gather
// totals on every write.
type writeCollector struct {
	w     io.Writer
	err   error
	total int
}

func newWriteCollector(w io.Writer) *writeCollector {
	return &writeCollector{w: w}
}

// Attempts to write bytes to the underlying writer.
func (w *writeCollector) Write(b []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	n, err := w.w.Write(b)
	w.err = err
	w.total += n
	return n, err
}

// Returns the total bytes written, and any error that occurred.
func (w *writeCollector) Totals() (int, error) {
	return w.total, w.err
}

var _ io.Writer = &writeCollector{}

// Rewinding reader is used to wrap an io.Reader and allow bytes
// to be pushed back on to read again later.
type rewindingReader interface {
	io.Reader
	// Push bytes back on to the reader. This is a FILO stack.
	Rewind(b []byte)
}

// Implementation of the rewindingReader.
type rwReader struct {
	reader  io.Reader
	rewound [][]byte
	// idx points to the position where the next rewound byte
	// would be stored
	idx int
}

func newRwReader(r io.Reader) *rwReader {
	return &rwReader{
		reader:  r,
		rewound: make([][]byte, 8),
		idx:     0,
	}
}

// Reads bytes into p, like an io.Reader. If there are rewound
// bytes, then they will be read first.
func (r *rwReader) Read(p []byte) (n int, err error) {
	if r.idx == 0 {
		return r.reader.Read(p)
	}

	b := r.rewound[r.idx-1]
	lb, lp := len(b), len(p)

	copy(p, b)

	if lp >= lb {
		r.idx--
		return lb, nil
	} else {
		r.rewound[r.idx-1] = b[lp:]
		return lp, nil
	}
}

// Pushes byes back on to the reader, creating a FILO stack.
// After the stack is empty (no more rewound bytes), it will
// resume reading from the underlying reader.
func (r *rwReader) Rewind(p []byte) {
	r.rewound[r.idx] = p
	r.idx++
}

var _ rewindingReader = &rwReader{}
