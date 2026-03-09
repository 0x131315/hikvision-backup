package fs

import (
	"fmt"
	"log/slog"
	"os"
)

func Init(downloadDir string) error {
	return createDir(downloadDir)
}

func RemoveFile(path string) {
	slog.Debug("remove file", "path", path)
	if !IsFileExist(path) {
		slog.Debug("file not found, skip remove", "path", path)
		return
	}

	err := os.Remove(path)
	if err != nil {
		slog.Error("Failed to remove file", "path", path, "error", err)
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
		return false
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
		return false
	}

	return true
}

func createDir(path string) error {
	if isDirExist(path) {
		return nil
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		slog.Error("Failed to create a directory", "path", path, "error", err)
		return err
	}
	slog.Debug("Directory created", "path", path)
	return nil
}

func getPathInfo(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		slog.Error("Failed check path", "path", path, "error", err)
		return nil, fmt.Errorf("failed to stat path %s: %w", path, err)
	}
	return info, err
}
