package api

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func TestBuildSearchRequest(t *testing.T) {
	start := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	end := start.Add(2 * time.Hour)

	body, err := buildSearchRequest(10, 50, 3, &start, &end)
	if err != nil {
		t.Fatalf("buildSearchRequest error: %v", err)
	}

	var parsed CMSearchDescription
	if err := xml.Unmarshal([]byte(body), &parsed); err != nil {
		t.Fatalf("unmarshal search xml: %v", err)
	}

	if parsed.MaxResults != 50 {
		t.Fatalf("expected maxResults=50, got %d", parsed.MaxResults)
	}
	if parsed.SearchResultPosition != 10 {
		t.Fatalf("expected searchResultPosition=10, got %d", parsed.SearchResultPosition)
	}
	if parsed.TrackID != TypeVideo {
		t.Fatalf("expected trackID=%d, got %d", TypeVideo, parsed.TrackID)
	}
	if parsed.TimeSpanList.TimeSpan.StartTime != start.Format(timeFormat) {
		t.Fatalf("unexpected startTime: %s", parsed.TimeSpanList.TimeSpan.StartTime)
	}
	if parsed.TimeSpanList.TimeSpan.EndTime != end.Format(timeFormat) {
		t.Fatalf("unexpected endTime: %s", parsed.TimeSpanList.TimeSpan.EndTime)
	}
}

func TestBuildDownloadRequest(t *testing.T) {
	video := Video{Url: "/ISAPI/ContentMgmt/download?size=123"}

	body, err := buildDownloadRequest(video)
	if err != nil {
		t.Fatalf("buildDownloadRequest error: %v", err)
	}
	if !strings.Contains(body, "<playbackURI>"+video.Url+"</playbackURI>") {
		t.Fatalf("expected playbackURI in body, got: %s", body)
	}
}
