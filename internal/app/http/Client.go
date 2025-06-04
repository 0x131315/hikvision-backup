package http

import (
	"bytes"
	"fmt"
	"github.com/0x131315/hikvision-backup/internal/app/config"
	"github.com/0x131315/hikvision-backup/internal/app/util"
	"io"
	"log"
	"net/http"
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
var logger *log.Logger

func init() {
	logger = util.GetLogger()
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
	logger.Printf("http %s %s %s\n", method, uri,
		strings.Replace(strings.Replace(body, "\n", "", -1), " ", "", -1),
	)
	resp := doRequest(method, uri, body)

	var badCnt int

	for resp.code == http.StatusUnauthorized {
		resp.Stream.Close()
		badCnt++
		if badCnt > retryCnt {
			logger.Println("http retry failed, skipped")
			return nil
		}
		logger.Printf("http error %d, send auth %d/%d\n", resp.code, badCnt, retryCnt)
		updateDigest(resp.digest)
		authHeader = getNextAuthHeader(method, uri)
		resp = doRequest(method, uri, body)
	}

	badCnt = 0
	for resp.code == http.StatusInternalServerError {
		resp.Stream.Close()
		badCnt++
		if badCnt > retryCnt {
			logger.Println("http retry failed, skipped")
			return nil
		}
		logger.Printf("http error %d, resend request %d/%d\n", resp.code, badCnt, retryCnt)
		resp = doRequest(method, uri, body)
	}

	if resp.code != http.StatusOK {
		logger.Printf("http error %d\n", resp.code)
		util.FatalError(fmt.Sprintf("http status code: %d. Reqest: %s %s", resp.code, method, uri))
	}

	return resp
}

func doRequest(method, uri, body string) *Response {
	req, err := http.NewRequest(method, buildUrl(uri), bytes.NewBufferString(body))
	if err != nil {
		util.FatalError("Failed to build request", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		util.FatalError("Failed to send request", err)
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
