package api

import (
	"testing"
	"time"
)

func TestParseResponseValidXML(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<CMSearchResult xmlns="http://www.hikvision.com/ver20/XMLSchema" version="2.0">
  <searchID>search</searchID>
  <responseStatus>true</responseStatus>
  <responseStatusStrg>OK</responseStatusStrg>
  <totalMatches>1</totalMatches>
  <numOfMatches>1</numOfMatches>
  <matchList>
    <searchMatchItem>
      <sourceID>src</sourceID>
      <trackID>101</trackID>
      <timeSpan>
        <startTime>2025-01-02T03:04:05Z</startTime>
        <endTime>2025-01-02T03:05:05Z</endTime>
      </timeSpan>
      <mediaSegmentDescriptor>
        <contentType>video</contentType>
        <codecType>H.264</codecType>
        <playbackURI>/download?size=100</playbackURI>
      </mediaSegmentDescriptor>
      <metadataMatches>
        <metadataDescriptor>x</metadataDescriptor>
      </metadataMatches>
    </searchMatchItem>
  </matchList>
</CMSearchResult>`

	resp, err := parseResponse(input)
	if err != nil {
		t.Fatalf("parseResponse error: %v", err)
	}
	if resp.TotalMatches != 1 || len(resp.MatchList) != 1 {
		t.Fatalf("unexpected parsed response: %+v", resp)
	}
}

func TestParseResponseInvalidXML(t *testing.T) {
	_, err := parseResponse("<CMSearchResult>")
	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestBuildVideoSuccess(t *testing.T) {
	item := SearchMatchItem{
		TimeSpan: TimeSpan{
			StartTime: "2025-01-02T03:04:05Z",
			EndTime:   "2025-01-02T03:05:05Z",
		},
		MediaSegment: MediaSegmentDescriptor{
			PlaybackURI: "/ISAPI/ContentMgmt/download?size=600",
		},
	}

	video, err := buildVideo(item)
	if err != nil {
		t.Fatalf("buildVideo error: %v", err)
	}
	if video.Size != 600 {
		t.Fatalf("expected size 600, got %d", video.Size)
	}
	if video.Duration != time.Minute {
		t.Fatalf("expected duration 1m, got %s", video.Duration)
	}
}

func TestBuildVideoInvalidSize(t *testing.T) {
	item := SearchMatchItem{
		TimeSpan: TimeSpan{
			StartTime: "2025-01-02T03:04:05Z",
			EndTime:   "2025-01-02T03:05:05Z",
		},
		MediaSegment: MediaSegmentDescriptor{
			PlaybackURI: "/ISAPI/ContentMgmt/download?size=bad",
		},
	}

	_, err := buildVideo(item)
	if err == nil {
		t.Fatalf("expected size parse error")
	}
}
