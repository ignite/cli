// Package ctxreader brings context.Context to io.Reader
package ctxreader

import (
	"context"
	"io"
)

type (
	cancelableReader struct {
		io.Reader
		ctx context.Context
	}
	ioReadResult struct {
		n   int
		err error
	}
)

// New returns a new reader that emits a context error through its r.Read() method
// when ctx canceled.
func New(ctx context.Context, r io.Reader) io.Reader {
	return &cancelableReader{Reader: r, ctx: ctx}
}

func (r *cancelableReader) Read(data []byte) (n int, err error) {
	chData := make([]byte, len(data))
	chResult := make(chan ioReadResult, 1)

	go func() {
		n, err := r.Reader.Read(chData)
		chResult <- ioReadResult{n, err}
		close(chResult)
	}()

	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()

	case ret := <-chResult:
		copy(data, chData)
		return ret.n, ret.err
	}
}
