// Package ctxreader brings context.Context to io.Reader
package ctxreader

import (
	"context"
	"io"
)

type cancelableReader struct {
	io.Reader
	ctx context.Context
}

// New returns a new reader that emits a context error through its r.Read() method
// when ctx canceled.
func New(ctx context.Context, r io.Reader) io.Reader {
	return &cancelableReader{Reader: r, ctx: ctx}
}

func (r *cancelableReader) Read(data []byte) (n int, err error) {
	isRead := make(chan struct{})

	go func() {
		n, err = r.Reader.Read(data)
		close(isRead)
	}()

	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()

	case <-isRead:
		return
	}
}
