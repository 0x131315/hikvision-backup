package api

import (
	"encoding/xml"
	"log/slog"
	"os"
	"time"
)

type TimeSpanList struct {
	TimeSpan TimeSpan `xml:"timeSpan"`
}

type CMSearchDescription struct {
	XMLName              xml.Name     `xml:"CMSearchDescription"`
	MaxResults           int          `xml:"maxResults"`
	SearchResultPosition int          `xml:"searchResultPosition"`
	TimeSpanList         TimeSpanList `xml:"timeSpanList"`
	TrackID              int          `xml:"trackID"`
	SearchID             string       `xml:"searchID"`
}

type DownloadRequest struct {
	XMLName     xml.Name `xml:"downloadRequest"`
	PlaybackURI string   `xml:"playbackURI"`
}

func buildSearchRequest(offset, limit, lastDays int, timestart, timeend *time.Time) string {
	if timestart == nil {
		old := time.Now().AddDate(0, 0, -1*lastDays)
		timestart = &old
	}
	if timeend == nil {
		now := time.Now().UTC()
		timeend = &now
	}
	req := &CMSearchDescription{
		MaxResults:           limit,
		SearchResultPosition: offset,
		TrackID:              TypeVideo,
		SearchID:             "search",
		TimeSpanList: TimeSpanList{
			TimeSpan: TimeSpan{
				StartTime: timestart.Format(timeFormat),
				EndTime:   timeend.Format(timeFormat),
			},
		},
	}

	return buildXml(req)
}

func buildDownloadRequest(file Video) string {
	req := &DownloadRequest{
		PlaybackURI: file.Url,
	}

	return buildXml(req)
}

func buildXml[T DownloadRequest | CMSearchDescription](req *T) string {
	str, err := xml.MarshalIndent(req, "", "  ")
	if err != nil {
		slog.Error("xml.MarshalIndent() failed", "error", err)
		os.Exit(1)
	}
	body := string(str)

	return body
}
