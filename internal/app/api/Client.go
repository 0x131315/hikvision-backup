package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/0x131315/hikvision-backup/internal/app/config"
	"github.com/0x131315/hikvision-backup/internal/app/http"
)

const (
	TypeVideo = 101
	TypePhoto = 103
)
const timeFormat = "2006-01-02T15:04:05Z"
const limit = 200
const searchPath = "/ISAPI/ContentMgmt/search"
const downloadPath = "/ISAPI/ContentMgmt/download"

type Video struct {
	Url      string
	Time     time.Time
	Duration time.Duration
	Size     int
}
type VideoList map[string]Video

type ApiClient struct {
	ctx        context.Context
	httpClient *http.Client
	conf       config.Config
}

func NewApiClient(ctx context.Context, conf config.Config) *ApiClient {
	return &ApiClient{ctx: ctx, httpClient: http.NewHttpClient(ctx, conf), conf: conf}
}

func (api *ApiClient) GetVideoList() VideoList {
	slog.Info("Request remote file list...")
	var listVideos = make(VideoList)
	var timebreak = time.Now().UTC()
	var timestart = time.Now().UTC().AddDate(0, 0, -1*api.conf.ScanLastDays)
	var timeend = timestart.AddDate(0, 1, 0)
	var cnt int
	var resp *CMSearchResult
	var video Video
	var offset int

	for {
		for {
			select {
			case <-api.ctx.Done():
				return listVideos
			default:
			}

			resp = parseResponse(
				respToStr(
					api.httpClient.Send("POST", searchPath, buildSearchRequest(offset, limit, api.conf.ScanLastDays, &timestart, &timeend)).Stream,
				),
			)
			if offset == 0 {
				slog.Info("found files", "count", resp.TotalMatches)
			}
			offset += limit

			if len(resp.MatchList) == 0 {
				break
			}

			cnt += len(resp.MatchList)
			slog.Debug("requested file info",
				"count", fmt.Sprintf("%d/%d", cnt, resp.TotalMatches),
				"total", len(listVideos)+len(resp.MatchList),
			)

			for _, item := range resp.MatchList {
				video = buildVideo(item)
				listVideos[video.Time.Format(timeFormat)] = video
			}

			if cnt >= resp.TotalMatches {
				break
			}
		}

		cnt = 0
		offset = 0
		timeend = timeend.AddDate(0, 1, 0)
		timestart = timestart.AddDate(0, 1, 0)

		if timestart.After(timebreak) {
			break
		}
	}

	slog.Info("Request remote file list complete", "count", len(listVideos))

	return listVideos
}

func (api *ApiClient) GetVideo(video Video) *http.Response {
	resp := api.httpClient.Send("GET", downloadPath, buildDownloadRequest(video))

	if resp != nil {
		resp.Stream = &ctxReadCloser{
			ctx:    api.ctx,
			reader: resp.Stream,
		}
	}

	return resp
}

func respToStr(resp io.ReadCloser) string {
	defer resp.Close()
	buff, err := io.ReadAll(resp)
	if err != nil {
		slog.Error("Failed to read response body", "error", err)
		os.Exit(1)
	}
	return string(buff)
}
