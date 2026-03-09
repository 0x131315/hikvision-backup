package syncer

import "testing"

func TestIsNumericLikeTrue(t *testing.T) {
	if !isNumericLike("12.34") {
		t.Fatal("expected numeric-like value")
	}
}

func TestParseMarkdownTableInvalid(t *testing.T) {
	if _, _, ok := parseMarkdownTable("just text"); ok {
		t.Fatal("plain text should not parse as table")
	}
}
