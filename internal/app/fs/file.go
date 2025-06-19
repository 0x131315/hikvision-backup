package fs

import (
	"github.com/0x131315/hikvision-backup/internal/app/config"
	"log/slog"
	"os"
)

func init() {
	createDir(config.Get().DownloadDir)
}

func RemoveFile(path string) {
	slog.Debug("remove file", "path", path)
	if !IsFileExist(path) {
		slog.Error("file not exist", "path", path)
		os.Exit(1)
	}

	err := os.Remove(path)
	if err != nil {
		slog.Error("Failed to remove file", "path", path, "error", err)
		os.Exit(1)
	}
}

func FileSize(path string) int {
	info, err := getPathInfo(path)
	if err != nil {
		return 0
	}
	return int(info.Size())
}

func IsFileExist(path string) bool {
	info, err := getPathInfo(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		slog.Error("The path is not a file", "path", path, "info", info)
		os.Exit(1)
	}

	return true
}

func isDirExist(path string) bool {
	info, err := getPathInfo(path)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		slog.Error("The path is not a directory", "path", path, "info", info)
		os.Exit(1)
	}

	return true
}

func createDir(path string) {
	if isDirExist(path) {
		return
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		slog.Error("Failed to create a directory", "path", path, "error", err)
		os.Exit(1)
	}
	slog.Debug("Directory created", "path", path)
}

func getPathInfo(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		slog.Error("Failed check path", "path", path, "error", err)
		os.Exit(1)
	}
	return info, err
}
