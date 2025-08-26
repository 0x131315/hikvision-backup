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

type File struct {
	name string
	path string
	size int
}

type App struct {
	ctx  context.Context
	api  *api.ApiClient
	conf config.Config
}

func NewApp(ctx context.Context, logLvl slog.Level, logHttp bool) *App {
	conf := config.Init(logLvl, logHttp)
	return &App{ctx: ctx, api: api.NewApiClient(ctx, conf), conf: conf}
}

func (app *App) Conf() config.Config {
	return app.conf
}

func (app *App) DownloadVideos() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		videos := app.api.GetVideoList()

		dates := make([]string, 0, len(videos))
		for date := range videos {
			dates = append(dates, date)
		}
		sort.Strings(dates)

		for idx, date := range dates {
			select {
			case <-app.ctx.Done():
				return
			default:
			}

			slog.Info(fmt.Sprintf("Request remote file: %d/%d. Left: %d", idx+1, len(dates), len(dates)-(idx+1)))
			app.saveVideo(videos[date])
		}
	}()

	wg.Wait()
}
func (app *App) saveVideo(video api.Video) {
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

	var videoResp *http.BinaryResponse
	var badCnt, streamSize int
	retryCnt := app.conf.RetryCnt
	for {
		outFile, err := os.Create(file.path)
		if err != nil {
			slog.Error("Failed to create file", "path", file.path, "err", err)
			os.Exit(1)
		}

		slog.Debug("start download", "file", file.name)
		videoResp = app.api.GetVideo(video)
		if videoResp == nil || videoResp.Size() == 0 {
			slog.Debug("download failed", "file", file.name)
			fs.RemoveFile(file.path)
			return
		}
		streamSize = videoResp.Size()

		slog.Debug("writing start", "file", file.name, "expected size", util.FormatFileSize(streamSize))
		stream := io.TeeReader(videoResp.Stream(), util.BuildProgressBar(streamSize, "b"))
		_, err = io.Copy(outFile, stream)
		videoResp.Stream().Close()
		outFile.Close()

		//success branch
		if err == nil {
			break
		}

		//cancel branch
		select {
		case <-app.ctx.Done():
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
	if filesize != streamSize {
		slog.Error("file size mismatch",
			"loaded", fmt.Sprintf("%s/%s", util.FormatFileSize(filesize), util.FormatFileSize(videoResp.Size())),
			"diff", util.FormatFileSize(int(math.Abs(float64(streamSize-filesize)))),
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
