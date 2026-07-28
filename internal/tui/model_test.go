package tui

import (
	"strings"
	"testing"

	"github.com/necrasov-ilya/trustmefetch/internal/config"
	"github.com/necrasov-ilya/trustmefetch/internal/system"
	"github.com/necrasov-ilya/trustmefetch/internal/theme"
)

func TestDistroTaglineIsPinnedInPreview(t *testing.T) {
	selected := theme.Must("nixos")
	cfg := config.Default()
	cfg.Theme = selected.ID
	cfg.DistroJokes = true
	model := New(system.Info{}, cfg)
	preview := model.renderPreview(100, 28)
	if !strings.Contains(preview, selected.Tagline) {
		t.Fatal("enabled distro tagline is missing from the preview")
	}
}

func TestDisabledDistroTaglineIsHiddenFromPreview(t *testing.T) {
	selected := theme.Must("nixos")
	cfg := config.Default()
	cfg.Theme = selected.ID
	cfg.DistroJokes = false
	model := New(system.Info{}, cfg)
	preview := model.renderPreview(100, 28)
	if strings.Contains(preview, selected.Tagline) {
		t.Fatal("disabled distro tagline is visible in the preview")
	}
}
