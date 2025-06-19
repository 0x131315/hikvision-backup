package api

import (
	"encoding/xml"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"time"
)

type CMSearchResult struct {
	XMLName            xml.Name          `xml:"http://www.hikvision.com/ver20/XMLSchema CMSearchResult"`
	Version            string            `xml:"version,attr"`
	SearchID           string            `xml:"searchID"`
	ResponseStatus     bool              `xml:"responseStatus"`
	ResponseStatusStrg string            `xml:"responseStatusStrg"`
	TotalMatches       int               `xml:"totalMatches"`
	NumOfMatches       int               `xml:"numOfMatches"`
	MatchList          []SearchMatchItem `xml:"matchList>searchMatchItem"`
}

type SearchMatchItem struct {
	SourceID     string                 `xml:"sourceID"`
	TrackID      int                    `xml:"trackID"`
	TimeSpan     TimeSpan               `xml:"timeSpan"`
	MediaSegment MediaSegmentDescriptor `xml:"mediaSegmentDescriptor"`
	Metadata     MetadataMatches        `xml:"metadataMatches"`
}

type TimeSpan struct {
	StartTime string `xml:"startTime"`
	EndTime   string `xml:"endTime"`
}

type MediaSegmentDescriptor struct {
	ContentType string `xml:"contentType"`
	CodecType   string `xml:"codecType"`
	PlaybackURI string `xml:"playbackURI"`
}

type MetadataMatches struct {
	MetadataDescriptor string `xml:"metadataDescriptor"`
}

func parseResponse(data string) *CMSearchResult {
	var result CMSearchResult
	err := xml.Unmarshal([]byte(data), &result)
	if err != nil {
		slog.Error("Failed to parse XML data", "error", err)
		os.Exit(1)
	}

	return &result
}

func buildVideo(item SearchMatchItem) Video {
	u, err := url.Parse(item.MediaSegment.PlaybackURI)
	if err != nil {
		slog.Error("Failed to parse url", "url", item.MediaSegment.PlaybackURI, "error", err)
		os.Exit(1)
	}
	sizeStr := u.Query().Get("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		slog.Error("Failed to parse size", "size", sizeStr, "error", err)
		os.Exit(1)
	}

	startTime, err := time.Parse(timeFormat, item.TimeSpan.StartTime)
	if err != nil {
		slog.Error("Failed to parse start time", "time", item.TimeSpan.StartTime, "error", err)
		os.Exit(1)
	}
	endTime, err := time.Parse(timeFormat, item.TimeSpan.EndTime)
	if err != nil {
		slog.Error("Failed to parse end time", "time", item.TimeSpan.EndTime, "error", err)
		os.Exit(1)
	}
	duration := endTime.Sub(startTime)

	return Video{
		Url:      item.MediaSegment.PlaybackURI,
		Time:     startTime.Local(),
		Duration: duration,
		Size:     size,
	}
}
