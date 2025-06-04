package fs

import (
	"fmt"
	"github.com/0x131315/hikvision-backup/internal/app/config"
	"github.com/0x131315/hikvision-backup/internal/app/util"
	"os"
)

func init() {
	createDir(config.Get().DownloadDir)
}

func RemoveFile(path string) {
	if !IsFileExist(path) {
		util.FatalError("Not found file for remove: " + path)
	}

	err := os.Remove(path)
	if err != nil {
		util.FatalError("Failed to remove file: "+path, err)
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
		util.FatalError("The path is not a file: " + path)
	}

	return true
}

func isDirExist(path string) bool {
	info, err := getPathInfo(path)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		util.FatalError("The path is not a directory: " + path)
	}

	return true
}

func createDir(path string) {
	if isDirExist(path) {
		return
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		util.FatalError("Failed to create a directory: "+path, err)
	}
	fmt.Println("Directory created: ", path)
}

func getPathInfo(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		util.FatalError("Failed check path: "+path, err)
	}
	return info, err
}
