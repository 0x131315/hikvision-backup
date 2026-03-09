package api

import (
	"context"
	"fmt"
	"log/slog"
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
	if api.conf.ScanLastDays == 0 {
		return api.getVideoListUnbounded()
	}

	var listVideos = make(VideoList)
	var timebreak = time.Now().UTC()
	var timestart = time.Now().UTC().AddDate(0, 0, -1*api.conf.ScanLastDays)
	var timeend = timestart.AddDate(0, 1, 0)
	var cnt int
	var resp *CMSearchResult
	var offset int

	for {
		for {
			select {
			case <-api.ctx.Done():
				return listVideos
			default:
			}

			requestStr, err := buildSearchRequest(offset, limit, api.conf.ScanLastDays, &timestart, &timeend)
			if err != nil {
				slog.Error("failed to build search request", "error", err)
				return listVideos
			}
			response := api.httpClient.Send("POST", searchPath, requestStr)
			if response == nil {
				slog.Error("failed to request file list")
				return listVideos
			}
			resp, err = parseResponse(response.Value())
			if err != nil {
				slog.Error("failed to parse search response", "error", err)
				return listVideos
			}
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
				video, err := buildVideo(item)
				if err != nil {
					slog.Warn("skip invalid video metadata", "error", err)
					continue
				}
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

func (api *ApiClient) getVideoListUnbounded() VideoList {
	var listVideos = make(VideoList)
	timestart := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	timeend := time.Now().UTC()
	var cnt int
	var offset int

	for {
		select {
		case <-api.ctx.Done():
			return listVideos
		default:
		}

		requestStr, err := buildSearchRequest(offset, limit, api.conf.ScanLastDays, &timestart, &timeend)
		if err != nil {
			slog.Error("failed to build search request", "error", err)
			return listVideos
		}
		response := api.httpClient.Send("POST", searchPath, requestStr)
		if response == nil {
			slog.Error("failed to request file list")
			return listVideos
		}
		resp, err := parseResponse(response.Value())
		if err != nil {
			slog.Error("failed to parse search response", "error", err)
			return listVideos
		}
		if offset == 0 {
			slog.Info("found files", "count", resp.TotalMatches)
		}

		if len(resp.MatchList) == 0 {
			break
		}

		cnt += len(resp.MatchList)
		slog.Debug("requested file info",
			"count", fmt.Sprintf("%d/%d", cnt, resp.TotalMatches),
			"total", len(listVideos)+len(resp.MatchList),
		)

		for _, item := range resp.MatchList {
			video, err := buildVideo(item)
			if err != nil {
				slog.Warn("skip invalid video metadata", "error", err)
				continue
			}
			listVideos[video.Time.Format(timeFormat)] = video
		}

		if cnt >= resp.TotalMatches {
			break
		}
		offset += limit
	}

	slog.Info("Request remote file list complete", "count", len(listVideos))

	return listVideos
}

func (api *ApiClient) GetVideo(video Video) *http.BinaryResponse {
	request, err := buildDownloadRequest(video)
	if err != nil {
		slog.Error("failed to build download request", "error", err)
		return nil
	}

	return api.httpClient.GetStream("GET", downloadPath, request)
}
