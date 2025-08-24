package api

import (
	"context"
	"io"
)

type ctxReadCloser struct {
	ctx    context.Context
	reader io.ReadCloser
}

func (cr *ctxReadCloser) Read(p []byte) (n int, err error) {
	select {
	case <-cr.ctx.Done():
		return 0, cr.ctx.Err()
	default:
		return cr.reader.Read(p)
	}
}

func (cr *ctxReadCloser) Close() error {
	return cr.reader.Close()
}
