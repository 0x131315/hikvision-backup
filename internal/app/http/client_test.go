package http

import (
	"context"
	"io"
	stdhttp "net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/0x131315/hikvision-backup/internal/app/config"
)

func newTestClient(baseURL string, retry int) *Client {
	conf := config.Config{
		BaseURL:     baseURL,
		RetryCnt:    retry,
		HttpTimeout: 3,
	}
	return NewHttpClient(context.Background(), conf)
}

func TestSend_RetryAndSuccess(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n < 3 {
			w.WriteHeader(stdhttp.StatusInternalServerError)
			_, _ = w.Write([]byte("retry"))
			return
		}
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL, 2)
	resp := c.Send(stdhttp.MethodPost, "/test", "<x/>")
	if resp == nil {
		t.Fatalf("expected non-nil response")
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("expected 3 attempts, got %d", got)
	}
	if resp.Value() != "ok" {
		t.Fatalf("expected body 'ok', got %q", resp.Value())
	}
}

func TestSend_RetryExhaustedReturnsNil(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(stdhttp.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL, 1)
	resp := c.Send(stdhttp.MethodPost, "/test", "<x/>")
	if resp != nil {
		t.Fatalf("expected nil response on retry exhaustion")
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Fatalf("expected 2 attempts, got %d", got)
	}
}

func TestGetStream_RetryAndSuccess(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n == 1 {
			w.WriteHeader(stdhttp.StatusUnauthorized)
			_, _ = w.Write([]byte("unauthorized"))
			return
		}
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte("stream-ok"))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL, 1)
	resp := c.GetStream(stdhttp.MethodGet, "/video", "<req/>")
	if resp == nil {
		t.Fatalf("expected non-nil stream response")
	}
	defer resp.Stream().Close()

	body, err := io.ReadAll(resp.Stream())
	if err != nil {
		t.Fatalf("read stream: %v", err)
	}
	if string(body) != "stream-ok" {
		t.Fatalf("expected stream body 'stream-ok', got %q", string(body))
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Fatalf("expected 2 attempts, got %d", got)
	}
}
