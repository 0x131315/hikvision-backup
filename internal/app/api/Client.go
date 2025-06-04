package api

import (
	"github.com/0x131315/hikvision-backup/internal/app/config"
	"github.com/0x131315/hikvision-backup/internal/app/http"
	"github.com/0x131315/hikvision-backup/internal/app/util"
	"io"
	"log"
	"time"
)

const (
	TypeVideo = 101
	TypePhoto = 103
)
const timeFormat = "2006-01-02T15:04:05Z"
const limit = 200
const searchPath = "/ISAPI/ContentMgmt/search"
const downloadPath = "/ISAPI/ContentMgmt/download"

var logger *log.Logger

type Video struct {
	Url      string
	Time     time.Time
	Duration time.Duration
	Size     int
}
type VideoList map[string]Video

func init() {
	logger = util.GetLogger()
}

func GetVideoList() VideoList {
	var result = make(VideoList)
	var timebreak = time.Now().UTC()
	var timestart = time.Now().UTC().AddDate(0, 0, -1*config.Get().ScanLastDays)
	var timeend = timestart.AddDate(0, 1, 0)
	var cnt int
	var resp *CMSearchResult
	var video Video
	var offset int

	for {
		for {
			resp = parseResponse(
				respToStr(
					http.Send("POST", searchPath, buildSearchRequest(offset, limit, &timestart, &timeend)).Stream,
				),
			)
			if offset == 0 {
				logger.Printf("found files: %d\n", resp.TotalMatches)
			}
			offset += limit

			if len(resp.MatchList) == 0 {
				break
			}

			cnt += len(resp.MatchList)
			logger.Printf("requested file info: %d/%d, total: %d\n", cnt, resp.TotalMatches, len(result)+len(resp.MatchList))

			for _, item := range resp.MatchList {
				video = buildVideo(item)
				result[video.Time.Format(timeFormat)] = video
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

	return result
}

func GetVideo(video Video) *http.Response {
	return http.Send("GET", downloadPath, buildDownloadRequest(video))
}

func respToStr(resp io.ReadCloser) string {
	defer resp.Close()
	buff, err := io.ReadAll(resp)
	if err != nil {
		util.FatalError("Failed to read response body", err)
	}
	return string(buff)
}
