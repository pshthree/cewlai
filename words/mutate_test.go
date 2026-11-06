package words

import "testing"

func TestMutateWords(t *testing.T) {
	cfg := DefaultMutateConfig()
	result := MutateWords([]string{"admin"}, cfg)

	expected := []string{
		"admin", "ADMIN", "Admin",
		"nimda",
		"4dm1n",
		"admin1", "admin123", "admin!", "admin2026",
		"Admin!", "Admin123", "Admin2026",
	}
	resultSet := make(map[string]struct{})
	for _, w := range result {
		resultSet[w] = struct{}{}
	}

	for _, e := range expected {
		if _, ok := resultSet[e]; !ok {
			t.Errorf("MutateWords missing expected word: %q", e)
		}
	}
}

func TestMutateWords_NoDuplicates(t *testing.T) {
	cfg := DefaultMutateConfig()
	result := MutateWords([]string{"test", "Test", "TEST"}, cfg)
	seen := make(map[string]struct{})
	for _, w := range result {
		if _, exists := seen[w]; exists {
			t.Errorf("MutateWords produced duplicate: %q", w)
		}
		seen[w] = struct{}{}
	}
}

func TestMutateWords_LengthFilter(t *testing.T) {
	cfg := DefaultMutateConfig()
	cfg.MinLength = 5
	cfg.MaxLength = 8
	result := MutateWords([]string{"admin"}, cfg)
	for _, w := range result {
		if len(w) < 5 || len(w) > 8 {
			t.Errorf("MutateWords produced word outside length bounds: %q (len %d)", w, len(w))
		}
	}
}

func TestMutateWords_LengthFilterRunes(t *testing.T) {
	// "café" is 4 runes / 5 bytes. With MinLength=4 the base word must be
	// kept; a byte-based check would drop it only above 5, and MaxLength=4
	// must still admit it. This guards the rune-count fix.
	cfg := DefaultMutateConfig()
	cfg.Capitalize = false
	cfg.Uppercase = false
	cfg.Lowercase = false
	cfg.Reverse = false
	cfg.LeetEnabled = false
	cfg.Suffixes = nil
	cfg.Prefixes = nil
	cfg.MinLength = 4
	cfg.MaxLength = 4

	result := MutateWords([]string{"café"}, cfg)
	found := false
	for _, w := range result {
		if w == "café" {
			found = true
		}
	}
	if !found {
		t.Errorf("MutateWords dropped 4-rune word %q under rune bounds 4..4: got %v", "café", result)
	}
}

func TestMutateWords_DisableFeatures(t *testing.T) {
	cfg := DefaultMutateConfig()
	cfg.Reverse = false
	cfg.LeetEnabled = false
	cfg.Uppercase = false
	cfg.Suffixes = nil
	cfg.Prefixes = nil
	result := MutateWords([]string{"admin"}, cfg)

	resultSet := make(map[string]struct{})
	for _, w := range result {
		resultSet[w] = struct{}{}
	}

	if _, ok := resultSet["nimda"]; ok {
		t.Error("reverse should be disabled")
	}
	if _, ok := resultSet["4dm1n"]; ok {
		t.Error("leet should be disabled")
	}
	if _, ok := resultSet["ADMIN"]; ok {
		t.Error("uppercase should be disabled")
	}
}

func TestMutateWords_CustomConfig(t *testing.T) {
	cfg := MutateConfig{
		Leet:        map[string]string{"a": "@"},
		Suffixes:    []string{"2026"},
		Prefixes:    []string{"pre_"},
		Capitalize:  true,
		LeetEnabled: true,
	}
	result := MutateWords([]string{"admin"}, cfg)

	resultSet := make(map[string]struct{})
	for _, w := range result {
		resultSet[w] = struct{}{}
	}

	if _, ok := resultSet["@dmin"]; !ok {
		t.Error("custom leet a->@ missing")
	}
	if _, ok := resultSet["admin2026"]; !ok {
		t.Error("custom suffix missing")
	}
	if _, ok := resultSet["pre_admin"]; !ok {
		t.Error("custom prefix missing")
	}
}

func BenchmarkMutateWords(b *testing.B) {
	cfg := DefaultMutateConfig()
	input := []string{"people", "who", "go", "all", "in", "are", "rare", "keep", "them"}
	for b.Loop() {
		_ = MutateWords(input, cfg)
	}
}
