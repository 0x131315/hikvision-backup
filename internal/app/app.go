package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/0x131315/hikvision-backup/internal/app/api"
	"github.com/0x131315/hikvision-backup/internal/app/config"
	"github.com/0x131315/hikvision-backup/internal/app/fs"
	"github.com/0x131315/hikvision-backup/internal/app/http"
	"github.com/0x131315/hikvision-backup/internal/app/util"
)

const retryCnt = 3

type File struct {
	name string
	path string
	size int
}

func DownloadVideos(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	videos := api.GetVideoList(ctx)

	dates := make([]string, 0, len(videos))
	for date := range videos {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for idx, date := range dates {
		select {
		case <-ctx.Done():
			return
		default:
		}

		slog.Info(fmt.Sprintf("Request remote file: %d/%d. Left: %d", idx+1, len(dates), len(dates)-(idx+1)))
		saveVideo(ctx, videos[date])
	}

}
func saveVideo(ctx context.Context, video api.Video) {
	file := buildFile(video)
	slog.Info("processing file", "file", file.name, "size", util.FormatFileSize(file.size))

	if fs.IsFileExist(file.path) {
		slog.Debug("file exist", "path", file.path)
		filesize := fs.FileSize(file.path)
		if filesize >= file.size {
			slog.Debug("file valid, skip downloading")
			return
		}
		slog.Debug("file invalid",
			"loaded", fmt.Sprintf("%s/%s", util.FormatFileSize(filesize), util.FormatFileSize(file.size)),
			"diff", util.FormatFileSize(int(math.Abs(float64(file.size-filesize)))),
		)
		fs.RemoveFile(file.path)
	}

	var videoResp *http.Response
	var badCnt int
	for {
		outFile, err := os.Create(file.path)
		if err != nil {
			slog.Error("Failed to create file", "path", file.path, "err", err)
			os.Exit(1)
		}

		slog.Debug("start download", "file", file.name)
		videoResp = api.GetVideo(ctx, video)
		if videoResp == nil {
			slog.Debug("download failed", "file", file.name)
			fs.RemoveFile(file.path)
			return
		}

		slog.Debug("writing start", "file", file.name, "expected size", util.FormatFileSize(videoResp.Size))
		writer := io.MultiWriter(outFile, util.BuildProgressBar(videoResp.Size, "b"))
		_, err = io.Copy(writer, videoResp.Stream)
		videoResp.Stream.Close()
		outFile.Close()

		//success branch
		if err == nil {
			break
		}

		//cancel branch
		select {
		case <-ctx.Done():
			fs.RemoveFile(file.path)
			return
		default:
		}

		//error branch
		slog.Error("Failed to write file", "file", file.name, "err", err)
		badCnt++
		if badCnt > retryCnt {
			slog.Debug("bad file skipped", "file", file.name)
			break
		}
		slog.Debug(fmt.Sprintf("retry: %d/%d", badCnt, retryCnt))
		fs.RemoveFile(file.path)
	}
	filesize := fs.FileSize(file.path)
	slog.Debug("writing complete", "file", file.name, "size", util.FormatFileSize(filesize))
	if badCnt > retryCnt {
		return
	}

	slog.Debug("validate")
	if filesize != videoResp.Size {
		slog.Error("file size mismatch",
			"loaded", fmt.Sprintf("%s/%s", util.FormatFileSize(filesize), util.FormatFileSize(videoResp.Size)),
			"diff", util.FormatFileSize(int(math.Abs(float64(videoResp.Size-filesize)))),
		)
		fs.RemoveFile(file.path)
		return
	}

	slog.Debug("complete")
}

func buildFile(video api.Video) File {
	filename := video.Time.Format("2006-01-02T15-04-05") + ".mp4"
	return File{
		name: filename,
		path: filepath.Join(config.Get().DownloadDir, filename),
		size: video.Size,
	}
}
