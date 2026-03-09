package http

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestCtxReadCloser_ReadCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	r := &ctxReadCloser{
		ctx:    ctx,
		reader: io.NopCloser(strings.NewReader("payload")),
	}

	buf := make([]byte, 8)
	_, err := r.Read(buf)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestCtxReadCloser_ReadDelegatesWhenActive(t *testing.T) {
	r := &ctxReadCloser{
		ctx:    context.Background(),
		reader: io.NopCloser(strings.NewReader("ok")),
	}

	buf := make([]byte, 2)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if n != 2 || string(buf) != "ok" {
		t.Fatalf("unexpected read result: n=%d buf=%q", n, string(buf))
	}
}
