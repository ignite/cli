package lineprefixer

import "io"

type Writer struct{}

func NewWriter(prefix string, w io.Writer) *Writer {
	return nil
}

func (w *Writer) Write(p []byte) (n int, err error) { return 0, nil }
