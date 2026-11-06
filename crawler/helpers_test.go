package crawler

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// When there are more pages than the character budget, perPage is 0 under naive
// integer division and the context collapses to only spaces. The guard must
// keep real content in the output.
func TestBuildContextFromPages_MorePagesThanLimit(t *testing.T) {
	pages := make([]string, 50)
	for i := range pages {
		pages[i] = "alpha"
	}
	limit := 10

	got := buildContextFromPages(pages, limit)
	if strings.TrimSpace(got) == "" {
		t.Fatalf("context collapsed to whitespace with %d pages and limit %d", len(pages), limit)
	}
	if len(got) > limit {
		t.Errorf("context length %d exceeds limit %d", len(got), limit)
	}
}

func TestBuildContextFromPages_Empty(t *testing.T) {
	if got := buildContextFromPages(nil, 100); got != "" {
		t.Errorf("expected empty context for no pages, got %q", got)
	}
}

// truncateRunes must never split a multibyte rune, so the result is always
// valid UTF-8 even when the cut lands inside a multibyte sequence.
func TestTruncateRunes_ValidUTF8(t *testing.T) {
	s := "héllo wörld" // mixed multibyte
	for n := 0; n <= utf8.RuneCountInString(s)+2; n++ {
		out := truncateRunes(s, n)
		if !utf8.ValidString(out) {
			t.Errorf("truncateRunes(%q, %d) produced invalid UTF-8: %q", s, n, out)
		}
		if rc := utf8.RuneCountInString(out); rc > n && n > 0 {
			t.Errorf("truncateRunes(%q, %d) returned %d runes", s, n, rc)
		}
	}
}
