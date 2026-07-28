package theme

import "testing"

func TestThemeCatalog(t *testing.T) {
	items := All()
	if len(items) < 32 {
		t.Fatalf("expected at least 32 themes, got %d", len(items))
	}
	seen := make(map[string]bool, len(items))
	for _, item := range items {
		if item.ID == "" || item.Name == "" || item.Distro == "" {
			t.Fatalf("incomplete theme: %+v", item)
		}
		if seen[item.ID] {
			t.Fatalf("duplicate theme id: %s", item.ID)
		}
		seen[item.ID] = true
		if len(Logo(item.Logo)) == 0 {
			t.Fatalf("theme %s has no logo", item.ID)
		}
		if !item.Joke && len(Logo(item.Logo)) < 10 {
			t.Fatalf("distro theme %s uses a reduced logo", item.ID)
		}
	}
}
