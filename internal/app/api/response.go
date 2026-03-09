package api

import (
	"encoding/xml"
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

func parseResponse(data string) (*CMSearchResult, error) {
	var result CMSearchResult
	err := xml.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func buildVideo(item SearchMatchItem) (Video, error) {
	u, err := url.Parse(item.MediaSegment.PlaybackURI)
	if err != nil {
		return Video{}, err
	}
	sizeStr := u.Query().Get("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return Video{}, err
	}

	startTime, err := time.Parse(timeFormat, item.TimeSpan.StartTime)
	if err != nil {
		return Video{}, err
	}
	endTime, err := time.Parse(timeFormat, item.TimeSpan.EndTime)
	if err != nil {
		return Video{}, err
	}
	duration := endTime.Sub(startTime)

	return Video{
		Url:      item.MediaSegment.PlaybackURI,
		Time:     startTime.Local(),
		Duration: duration,
		Size:     size,
	}, nil
}
