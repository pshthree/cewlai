package crawler

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const cacheDir = ".cewlai-cache"

type cachedResult struct {
	Result    CrawlResult `json:"result"`
	Timestamp time.Time   `json:"timestamp"`
}

// cacheKey derives the cache filename from the URL, depth, and every option
// that changes what the crawl actually produces. Flags like --secrets,
// --email, --offsite or --ai-context alter the stored CrawlResult, so leaving
// them out of the key would serve a result that does not match the request
// (for example an empty Secrets list after a prior run without --secrets).
func cacheKey(url string, depth int, opts CrawlOptions) string {
	seed := fmt.Sprintf(
		"%s:%d:em=%t:mt=%t:sc=%t:off=%t:cp=%t:csub=%t:cdom=%t:wn=%t:ctx=%d:mp=%d:mf=%d",
		url, depth,
		opts.ExtractEmails, opts.ExtractMeta, opts.ExtractSecrets, opts.Offsite,
		opts.CapturePaths, opts.CaptureSubdomains, opts.CaptureDomain,
		opts.WithNumbers, opts.MaxContext, opts.MaxPages, opts.MaxFiles,
	)
	h := sha256.Sum256([]byte(seed))
	return fmt.Sprintf("%x", h[:16])
}

func loadCache(url string, depth int, ttl time.Duration, opts CrawlOptions) (*CrawlResult, bool) {
	path := filepath.Join(cacheDir, cacheKey(url, depth, opts)+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var cached cachedResult
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, false
	}

	if time.Since(cached.Timestamp) > ttl {
		_ = os.Remove(path)
		return nil, false
	}

	return &cached.Result, true
}

func saveCache(url string, depth int, result *CrawlResult, opts CrawlOptions) {
	_ = os.MkdirAll(cacheDir, 0755)
	path := filepath.Join(cacheDir, cacheKey(url, depth, opts)+".json")

	cached := cachedResult{
		Result:    *result,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0644)
}
