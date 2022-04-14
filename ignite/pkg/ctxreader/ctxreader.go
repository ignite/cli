// Package ctxreader brings context.Context to io.Reader
package ctxreader

import (
	"context"
	"io"
	"sync"
)

type cancelableReader struct {
	io.Reader
	ctx context.Context
	m   sync.Mutex
	err error
}

// New returns a new reader that emits a context error through its r.Read() method
// when ctx canceled.
func New(ctx context.Context, r io.Reader) io.Reader {
	return &cancelableReader{Reader: r, ctx: ctx}
}

// Read implements io.Reader and it stops blocking when reading is completed
// or context is cancelled.
func (r *cancelableReader) Read(data []byte) (n int, err error) {
	r.m.Lock()
	defer r.m.Unlock()

	if r.err != nil {
		return 0, r.err
	}

	var (
		readerN   int
		readerErr error
	)
	isRead := make(chan struct{})
	go func() {
		readerN, readerErr = r.Reader.Read(data)
		close(isRead)
	}()

	select {
	case <-r.ctx.Done():
		r.err = r.ctx.Err()
		return 0, r.ctx.Err()
	case <-isRead:
		r.err = readerErr
		return readerN, readerErr
	}
}
