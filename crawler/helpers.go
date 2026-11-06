package crawler

import (
	"math/rand"
	"strings"
	"unicode/utf8"
)

func defaultContextLimit(opts CrawlOptions) int {
	if opts.MaxContext > 0 {
		return opts.MaxContext
	}
	return 4000
}

func buildContextFromPages(pageContexts []string, limit int) string {
	if len(pageContexts) == 0 {
		return ""
	}

	shuffled := make([]string, len(pageContexts))
	copy(shuffled, pageContexts)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	perPage := limit / len(shuffled)
	// With more pages than the character budget, integer division gives 0,
	// which would truncate every page to nothing and yield an all-spaces
	// context. Fall back to taking at least a little from each page and let
	// the overall limit below do the trimming.
	if perPage < 1 {
		perPage = 1
	}

	var b strings.Builder
	for _, page := range shuffled {
		if b.Len() >= limit {
			break
		}
		b.WriteString(truncateRunes(page, perPage))
		b.WriteString(" ")
	}

	return truncateRunes(b.String(), limit)
}

// truncateRunes returns at most n runes of s, never splitting a multibyte
// character. Slicing by bytes risks emitting invalid UTF-8 into the context
// that is handed to the LLM.
func truncateRunes(s string, n int) string {
	if n <= 0 || s == "" {
		return ""
	}
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	count := 0
	for i := range s {
		if count == n {
			return s[:i]
		}
		count++
	}
	return s
}

func mapKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
