package app

import (
	"github.com/0x131315/hikvision-backup/internal/app/api"
	"github.com/0x131315/hikvision-backup/internal/app/config"
	"github.com/0x131315/hikvision-backup/internal/app/fs"
	"github.com/0x131315/hikvision-backup/internal/app/http"
	"github.com/0x131315/hikvision-backup/internal/app/util"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
)

const retryCnt = 3

var logger *log.Logger

type File struct {
	name string
	path string
	size int
}

func init() {
	logger = util.GetLogger()
}

func DownloadVideos() {
	logger.Println("Request remote file list...")

	videos := api.GetVideoList()
	logger.Printf("Remote files: %d", len(videos))

	dates := make([]string, 0, len(videos))
	for date := range videos {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for idx, date := range dates {
		logger.Printf("Request remote file: %d/%d. Left: %d", idx+1, len(dates), len(dates)-(idx+1))
		saveVideo(videos[date])
		logger.Println("")
	}

}
func saveVideo(video api.Video) {
	file := buildFile(video)
	logger.Printf("processing %s", file.name)

	if fs.IsFileExist(file.path) {
		filesize := fs.FileSize(file.path)
		if filesize >= file.size {
			logger.Println("exist and valid")
			return
		}
		logger.Printf("exist and invalid: loaded %s/%s, diff %s. Removed",
			util.FormatFileSize(filesize),
			util.FormatFileSize(file.size),
			util.FormatFileSize(int(math.Abs(float64(file.size-filesize)))),
		)
		fs.RemoveFile(file.path)
	}

	var videoResp *http.Response
	var badCnt int
	for {
		outFile, err := os.Create(file.path)
		if err != nil {
			util.FatalError("Failed to create file: "+file.path, err)
		}

		logger.Println("start download")
		videoResp = api.GetVideo(video)
		if videoResp == nil {
			logger.Println("removed")
			fs.RemoveFile(file.path)
			return
		}

		logger.Printf("writing start, %s", util.FormatFileSize(videoResp.Size))
		writer := io.MultiWriter(outFile, util.BuildProgressBar(videoResp.Size, "b"))
		_, err = io.Copy(writer, videoResp.Stream)
		videoResp.Stream.Close()
		outFile.Close()
		if err == nil {
			break
		}
		logger.Printf("Failed to write file %s: %v", file.name, err)
		badCnt++
		if badCnt > retryCnt {
			logger.Println("bad file skipped")
			break
		}
		logger.Printf("retry: %d/%d\n", badCnt, retryCnt)
		logger.Println("removed")
		fs.RemoveFile(file.path)
	}
	filesize := fs.FileSize(file.path)
	logger.Printf("writing complete, %s", util.FormatFileSize(filesize))
	if badCnt > retryCnt {
		return
	}

	logger.Println("validate")
	if filesize != videoResp.Size {
		logger.Printf("validate failed: write %s/%s, diff %s. File removed",
			util.FormatFileSize(filesize),
			util.FormatFileSize(videoResp.Size),
			util.FormatFileSize(int(math.Abs(float64(videoResp.Size-filesize)))),
		)
		fs.RemoveFile(file.path)
		return
	}

	logger.Println("complete")
}

func buildFile(video api.Video) File {
	filename := video.Time.Format("2006-01-02T15-04-05") + ".mp4"
	return File{
		name: filename,
		path: filepath.Join(config.Get().DownloadDir, filename),
		size: video.Size,
	}
}
