package truncatedbuffer

import (
	"bytes"
)

// TruncatedBuffer contains a bytes buffer that has a limited capacity.
// The buffer is truncated on Write if the length reaches the maximum capacity.
// Only the first bytes are preserved.
type TruncatedBuffer struct {
	buf *bytes.Buffer
	cap int
}

// NewTruncatedBuffer returns a new TruncatedBuffer.
// If the provided cap is 0, the truncated buffer has no limit for truncating.
func NewTruncatedBuffer(cap int) *TruncatedBuffer {
	return &TruncatedBuffer{
		buf: &bytes.Buffer{},
		cap: cap,
	}
}

// GetBuffer returns the buffer.
func (b TruncatedBuffer) GetBuffer() *bytes.Buffer {
	return b.buf
}

// GetCap returns the maximum capacity of the buffer.
func (b TruncatedBuffer) GetCap() int {
	return b.cap
}

// Write implements io.Writer.
func (b *TruncatedBuffer) Write(p []byte) (n int, err error) {
	n, err = b.buf.Write(p)
	if err != nil {
		return n, err
	}

	// Check surplus bytes
	surplus := b.buf.Len() - b.cap

	if b.cap > 0 && surplus > 0 {
		b.buf.Truncate(b.cap)
	}

	return n, nil
}
