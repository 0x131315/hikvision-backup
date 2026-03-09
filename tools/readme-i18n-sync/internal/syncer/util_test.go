package syncer

import "testing"

func TestIsQuotaExceeded(t *testing.T) {
	if !isQuotaExceeded(assertErr("Quota exceeded")) {
		t.Fatal("expected true")
	}
	if isQuotaExceeded(assertErr("other error")) {
		t.Fatal("expected false")
	}
}

func TestHTMLUnescape(t *testing.T) {
	in := "a &amp; b &lt;tag&gt; &quot;x&quot; &#39;y&#39;"
	out := htmlUnescape(in)
	want := "a & b <tag> \"x\" 'y'"
	if out != want {
		t.Fatalf("unexpected htmlUnescape: %q", out)
	}
}

type simpleErr string

func (e simpleErr) Error() string { return string(e) }

func assertErr(s string) error { return simpleErr(s) }
