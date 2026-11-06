package crawler

import "testing"

// cacheKey must change when any crawl option that alters the stored result
// changes, otherwise a run with new extraction flags would be served a stale
// cached result that does not match the request.
func TestCacheKeyVariesByOptions(t *testing.T) {
	base := CrawlOptions{URL: "https://example.com", Depth: 2}
	baseKey := cacheKey(base.URL, base.Depth, base)

	variants := map[string]CrawlOptions{
		"secrets":    withOpt(base, func(o *CrawlOptions) { o.ExtractSecrets = true }),
		"emails":     withOpt(base, func(o *CrawlOptions) { o.ExtractEmails = true }),
		"meta":       withOpt(base, func(o *CrawlOptions) { o.ExtractMeta = true }),
		"offsite":    withOpt(base, func(o *CrawlOptions) { o.Offsite = true }),
		"paths":      withOpt(base, func(o *CrawlOptions) { o.CapturePaths = true }),
		"subdomains": withOpt(base, func(o *CrawlOptions) { o.CaptureSubdomains = true }),
		"domain":     withOpt(base, func(o *CrawlOptions) { o.CaptureDomain = true }),
		"context":    withOpt(base, func(o *CrawlOptions) { o.MaxContext = 8000 }),
		"maxpages":   withOpt(base, func(o *CrawlOptions) { o.MaxPages = 10 }),
	}

	for name, v := range variants {
		if got := cacheKey(v.URL, v.Depth, v); got == baseKey {
			t.Errorf("cacheKey did not change for option %q; stale cache would be served", name)
		}
	}
}

func TestCacheKeyStableForSameOptions(t *testing.T) {
	o := CrawlOptions{URL: "https://example.com", Depth: 3, ExtractSecrets: true, MaxContext: 4000}
	if cacheKey(o.URL, o.Depth, o) != cacheKey(o.URL, o.Depth, o) {
		t.Error("cacheKey is not deterministic for identical options")
	}
}

func withOpt(base CrawlOptions, mut func(*CrawlOptions)) CrawlOptions {
	cp := base
	mut(&cp)
	return cp
}
