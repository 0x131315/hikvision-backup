package api

import (
	"encoding/xml"
	"github.com/0x131315/hikvision-backup/internal/app/util"
	"net/url"
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
		util.FatalError("Failed to parse XML data", err)
	}

	return &result
}

func buildVideo(item SearchMatchItem) Video {
	u, err := url.Parse(item.MediaSegment.PlaybackURI)
	if err != nil {
		util.FatalError("Parse url error", err)
	}
	sizeStr := u.Query().Get("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		util.FatalError("invalid size: ", err)
	}

	startTime, err := time.Parse(timeFormat, item.TimeSpan.StartTime)
	if err != nil {
		util.FatalError("invalid starttime: ", err)
	}
	endTime, err := time.Parse(timeFormat, item.TimeSpan.EndTime)
	if err != nil {
		util.FatalError("invalid endtime: ", err)
	}
	duration := endTime.Sub(startTime)

	return Video{
		Url:      item.MediaSegment.PlaybackURI,
		Time:     startTime.Local(),
		Duration: duration,
		Size:     size,
	}
}
