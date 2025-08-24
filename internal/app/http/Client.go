package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/0x131315/hikvision-backup/internal/app/config"
)

type Response struct {
	code   int
	Stream io.ReadCloser
	Size   int
	digest string
}

type Client struct {
	ctx        context.Context
	client     *http.Client
	conf       config.Config
	digest     *Digest
	authHeader string
}

func NewHttpClient(ctx context.Context, conf config.Config) *Client {
	client := &http.Client{}

	//ignore proxy in env variables (for debug sessions)
	if conf.NoProxy {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: nil, // Отключаем использование прокси
			},
		}
	}

	return &Client{ctx: ctx, client: client, conf: conf, digest: NewDigest(conf)}
}

func (c *Client) Send(method, uri, body string) *Response {
	slog.Debug("http request",
		"method", method,
		"uri", uri,
		"body", strings.Replace(strings.Replace(body, "\n", "", -1), " ", "", -1),
	)
	resp := c.doRequest(method, uri, body)
	slog.Debug("http response",
		"code", resp.code,
		"size", resp.Size,
		"digest", resp.digest,
	)

	var badCnt int
	retryCnt := c.conf.RetryCnt

	for resp.code == http.StatusUnauthorized {
		resp.Stream.Close()
		badCnt++
		if badCnt > retryCnt {
			slog.Debug("http retry failed, skipped")
			return nil
		}
		slog.Debug("send auth", "count", fmt.Sprintf("%d/%d", badCnt, retryCnt))
		c.digest.updateDigest(resp.digest)
		c.authHeader = c.digest.getNextAuthHeader(method, uri)
		resp = c.doRequest(method, uri, body)
	}

	badCnt = 0
	for resp.code == http.StatusInternalServerError {
		resp.Stream.Close()
		badCnt++
		if badCnt > retryCnt {
			slog.Debug("http retry failed, skipped")
			return nil
		}
		slog.Debug("resend request", "count", fmt.Sprintf("%d/%d", badCnt, retryCnt))
		resp = c.doRequest(method, uri, body)
	}

	if resp.code != http.StatusOK {
		slog.Error("http error", "code", resp.code)
		os.Exit(1)
	}

	return resp
}

func (c *Client) doRequest(method, path, body string) *Response {
	req, err := http.NewRequest(method, buildUrl(path), bytes.NewBufferString(body))
	if err != nil {
		slog.Error("Failed to build request", "error", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Content-Type", "application/xml")

	resp, err := c.client.Do(req)
	if err != nil {
		slog.Error("Failed to send request", "error", err)
		os.Exit(1)
	}

	return &Response{
		code:   resp.StatusCode,
		Stream: resp.Body,
		Size:   int(resp.ContentLength),
		digest: resp.Header.Get("WWW-Authenticate"),
	}
}

func buildUrl(path string) string {
	host := config.Get().Host
	return fmt.Sprintf("http://%s%s", host, path)
}
