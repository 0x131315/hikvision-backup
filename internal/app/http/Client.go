package http

import (
	"context"
	"crypto/tls"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/0x131315/hikvision-backup/internal/app/config"
	"github.com/go-resty/resty/v2"
	"github.com/icholy/digest"
)

type Response struct {
	code int
	size int
}

func (r *Response) Size() int {
	return r.size
}

type StringResponse struct {
	Response
	value string
}

func (r *StringResponse) Value() string {
	return r.value
}

type BinaryResponse struct {
	Response
	stream io.ReadCloser
}

func (r *BinaryResponse) Stream() io.ReadCloser {
	return r.stream
}

func (c *Client) Send(method, uri, body string) *StringResponse {
	resp := c.send(method, uri, body, false)

	return &StringResponse{
		Response: Response{
			code: resp.StatusCode(),
			size: int(resp.Size()),
		},
		value: resp.String(),
	}
}

func (c *Client) GetStream(method, uri, body string) *BinaryResponse {
	resp := c.send(method, uri, body, true)

	return &BinaryResponse{
		Response: Response{
			code: resp.StatusCode(),
			size: int(resp.RawResponse.ContentLength),
		},
		stream: &ctxReadCloser{
			ctx:    c.ctx,
			reader: resp.RawBody(),
		},
	}
}

type Client struct {
	ctx    context.Context
	client *resty.Client
	conf   config.Config
}

func NewHttpClient(ctx context.Context, conf config.Config) *Client {
	//ignore proxy in env variables (for debug sessions)
	if conf.NoProxy {
		os.Unsetenv("http_proxy")
		os.Unsetenv("https_proxy")
	}

	client := resty.New()
	client.SetRetryWaitTime(500 * time.Millisecond)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	client.SetAllowGetMethodPayload(true)
	client.SetHeader("Accept", "application/xml")
	client.SetHeader("Content-Type", "application/xml")

	client.SetRetryCount(conf.RetryCnt)
	client.SetBaseURL("http://" + conf.Host)
	client.SetTimeout(time.Duration(conf.HttpTimeout) * time.Second)

	if conf.LogLvl == slog.LevelDebug {
		client.SetDebug(true)
	}

	if conf.User != "" {
		digestTransport := &digest.Transport{
			Username: conf.User,
			Password: conf.Pass,
			NoReuse:  true, // Forces new challenge (nonce) for each request
		}
		client.SetTransport(digestTransport)
	}

	return &Client{ctx: ctx, client: client, conf: conf}
}

func (c *Client) send(method, uri, body string, noParse bool) *resty.Response {
	slog.Debug("http request",
		"method", method,
		"uri", uri,
		"body", strings.Replace(strings.Replace(body, "\n", "", -1), " ", "", -1),
	)

	req := c.client.R().SetBody(body).SetDoNotParseResponse(noParse)
	resp, err := req.Execute(method, uri)
	if err != nil {
		slog.Error("Failed to send request", "error", err)
		os.Exit(1)
	}

	slog.Debug("http response",
		"code", resp.StatusCode(),
		"size", resp.RawResponse.ContentLength,
	)

	if resp.StatusCode() >= 500 {
		retryCnt := config.Get().RetryCnt
		for retryCnt > 0 && resp.StatusCode() >= 500 {
			slog.Debug("send retry request", "retryCnt", retryCnt)
			resp = c.send(method, uri, body, noParse)
		}
	}

	if resp.StatusCode() != http.StatusOK {
		slog.Error("http error", "code", resp.StatusCode())
		os.Exit(1)
	}

	return resp
}
