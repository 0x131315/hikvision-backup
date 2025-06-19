package http

import (
	"bytes"
	"fmt"
	"github.com/0x131315/hikvision-backup/internal/app/config"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	code   int
	Stream io.ReadCloser
	Size   int
	digest string
}

const contentType = "application/xml"
const retryCnt = 3

var client *http.Client
var authHeader string

func init() {
	if config.Get().NoProxy {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: nil, // Отключаем использование прокси
			},
		}
	} else {
		client = &http.Client{}
	}
}

func Send(method, uri, body string) *Response {
	slog.Debug("http request",
		"method", method,
		"uri", uri,
		"body", strings.Replace(strings.Replace(body, "\n", "", -1), " ", "", -1),
	)
	resp := doRequest(method, uri, body)
	slog.Debug("http response",
		"code", resp.code,
		"size", resp.Size,
		"digest", resp.digest,
	)

	var badCnt int

	for resp.code == http.StatusUnauthorized {
		resp.Stream.Close()
		badCnt++
		if badCnt > retryCnt {
			slog.Debug("http retry failed, skipped")
			return nil
		}
		slog.Debug("send auth", "count", fmt.Sprintf("%d/%d", badCnt, retryCnt))
		updateDigest(resp.digest)
		authHeader = getNextAuthHeader(method, uri)
		resp = doRequest(method, uri, body)
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
		resp = doRequest(method, uri, body)
	}

	if resp.code != http.StatusOK {
		slog.Error("http error", "code", resp.code)
		os.Exit(1)
	}

	return resp
}

func doRequest(method, uri, body string) *Response {
	req, err := http.NewRequest(method, buildUrl(uri), bytes.NewBufferString(body))
	if err != nil {
		slog.Error("Failed to build request", "error", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
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
