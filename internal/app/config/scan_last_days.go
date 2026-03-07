package config

import (
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const localFileTimeLayout = "2006-01-02T15-04-05"

func adjustScanLastDaysFromLocal(downloadDir string, scanLastDays int, backfillDays int) int {
	latestTime, ok := findLatestLocalFileTime(downloadDir)
	if !ok {
		slog.Debug("no local video files found for scan adjustment", "download_dir", downloadDir)
		return scanLastDays
	}

	target := latestTime.AddDate(0, 0, -backfillDays)
	days := daysSince(target)
	if days < 1 {
		days = 1
	}

	if scanLastDays == 0 {
		slog.Info("adjusted scan window from local files", "scan_last_days", days)
		return days
	}

	if days < scanLastDays {
		slog.Info("adjusted scan window from local files", "scan_last_days", days, "scan_last_days_env", scanLastDays)
		return days
	}

	return scanLastDays
}

func findLatestLocalFileTime(downloadDir string) (time.Time, bool) {
	entries, err := os.ReadDir(downloadDir)
	if err != nil {
		slog.Debug("failed to read download dir for scan adjustment", "download_dir", downloadDir, "error", err)
		return time.Time{}, false
	}

	var latest time.Time
	found := false
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.ToLower(filepath.Ext(name)) != ".mp4" {
			continue
		}

		base := strings.TrimSuffix(name, filepath.Ext(name))
		t, err := time.ParseInLocation(localFileTimeLayout, base, time.Local)
		if err != nil {
			continue
		}

		if !found || t.After(latest) {
			latest = t
			found = true
		}
	}

	return latest, found
}

func daysSince(t time.Time) int {
	diff := time.Since(t)
	if diff < 0 {
		return 0
	}

	return int(math.Ceil(diff.Hours() / 24))
}
