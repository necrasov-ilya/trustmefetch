package system

import "testing"

func TestFindStrings(t *testing.T) {
	value := map[string]any{
		"display": []any{
			map[string]any{"sppci_model": "Apple M42"},
			map[string]any{"nested": map[string]any{"sppci_model": "External GPU"}},
		},
	}
	got := findStrings(value, "sppci_model")
	if len(got) != 2 || got[0] != "Apple M42" || got[1] != "External GPU" {
		t.Fatalf("unexpected values: %#v", got)
	}
}

func TestLines(t *testing.T) {
	if got := lines("one\ntwo\n"); got != 2 {
		t.Fatalf("got %d, want 2", got)
	}
	if got := lines("\n"); got != 0 {
		t.Fatalf("got %d, want 0", got)
	}
}
