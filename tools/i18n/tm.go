package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func loadTM(path string) (tmFile, error) {
	var tm tmFile
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return tm, nil
		}
		return tm, err
	}
	if err := json.Unmarshal(data, &tm); err != nil {
		return tm, err
	}
	return tm, nil
}

func writeTM(path string, tm tmFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(tm, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
