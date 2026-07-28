package render

import (
	"reflect"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/necrasov-ilya/trustmefetch/internal/system"
	"github.com/necrasov-ilya/trustmefetch/internal/theme"
)

func TestFetchUsesCollectedValues(t *testing.T) {
	info := system.Info{User: "ada", Hostname: "mac", Host: "Mac42,1", Kernel: "25.0", Arch: "arm64", CPU: "Apple M42", Memory: "1 / 2 GiB", Disk: "3 / 4 GiB", OSVersion: "26.0", Build: "A1"}
	output := Fetch(theme.Must("arch"), info, Options{Width: 100, Color: false})
	for _, expected := range []string{"ada@mac", "Mac42,1", "Apple M42", "Arch Linux arm64"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("output does not contain %q:\n%s", expected, output)
		}
	}
}

func TestQuestionModeArguesThatMacIsLinux(t *testing.T) {
	selected := theme.Must("arch-btw")
	output := Fetch(selected, system.Info{}, Options{Width: 120, Color: false, Question: true})
	for _, expected := range []string{
		"YES. 100% LINUX!!!!!!",
		"Evidence: Arch Linux logo detected. KDE Plasma 6 confirmed.",
		"Objection: Darwin kernel dismissed as a harmless typo.",
		"Confidence: 100.00% (independently self-certified)",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("question output does not contain %q:\n%s", expected, output)
		}
	}
}

func TestRegularFetchDoesNotContainQuestionProof(t *testing.T) {
	output := Fetch(theme.Must("arch"), system.Info{}, Options{Width: 120, Color: false})
	if strings.Contains(output, "100% LINUX") || strings.Contains(output, "independently self-certified") {
		t.Fatalf("regular fetch contains question proof:\n%s", output)
	}
}

func TestDataRowsMatchFastfetchMacOSDefaultModules(t *testing.T) {
	rows := dataRows(theme.Must("rgb-linux"), system.Info{})
	labels := make([]string, len(rows))
	for index, row := range rows {
		labels[index] = row[0]
	}
	want := []string{
		"OS", "Host", "Kernel", "Uptime", "Packages", "Shell", "Display",
		"WM", "WM Theme", "Theme", "Font", "Cursor", "Terminal", "CPU", "GPU",
		"Memory", "Swap", "Disk", "Local IP", "Battery", "Power Adapter", "Locale",
	}
	if !reflect.DeepEqual(labels, want) {
		t.Fatalf("module order differs from fastfetch:\ngot  %v\nwant %v", labels, want)
	}
}

func TestColoredFetchHasFastfetchDefaultLineCount(t *testing.T) {
	output := Fetch(theme.Must("arch"), system.Info{}, Options{Width: 200, Color: true, ShowDistroTagline: false})
	if got, want := strings.Count(output, "\n"), 27; got != want {
		t.Fatalf("got %d output lines, want %d", got, want)
	}
}

func TestDistroTaglineCanBeToggled(t *testing.T) {
	selected := theme.Must("ubuntu")
	hidden := Fetch(selected, system.Info{}, Options{Width: 200, Color: false})
	shown := Fetch(selected, system.Info{}, Options{Width: 200, Color: false, ShowDistroTagline: true})
	if strings.Contains(hidden, selected.Tagline) {
		t.Fatal("distro tagline is visible while disabled")
	}
	if !strings.Contains(shown, selected.Tagline) {
		t.Fatal("distro tagline is missing while enabled")
	}
}

func TestJokeThemeAlwaysShowsTagline(t *testing.T) {
	selected := theme.Must("rgb-linux")
	output := Fetch(selected, system.Info{}, Options{Width: 200, Color: false})
	if !strings.Contains(output, selected.Tagline) {
		t.Fatal("joke theme tagline is missing")
	}
}

func TestFetchWrapsInsideRequestedWidth(t *testing.T) {
	info := system.Info{
		User: "ada", Hostname: "mac", Resolution: "Color LCD: 3456x2234 (1728 x 1117 @ 120.00Hz) [Built-in]",
		Font: ".AppleSystemUIFont [System], Helvetica [User]",
	}
	const width = 100
	output := Fetch(theme.Must("rgb-linux"), info, Options{Width: width, Color: false})
	for index, line := range strings.Split(strings.TrimRight(output, "\n"), "\n") {
		if got := lipgloss.Width(line); got > width {
			t.Fatalf("line %d is %d cells wide, want at most %d:\n%s", index+1, got, width, line)
		}
	}
}

func TestEveryThemeRendersWithoutFastfetchMarkers(t *testing.T) {
	info := system.Info{User: "ada", Hostname: "mac", Host: "Mac42,1", Kernel: "25.0", Arch: "arm64", CPU: "Apple M42", CPUUsage: "42%", GPU: "Apple M42", Memory: "1 / 2 GiB", Disk: "3 / 4 GiB", OSVersion: "26.0", Build: "A1"}
	for _, selected := range theme.All() {
		output := Fetch(selected, info, Options{Width: 120, Color: false, ShowDistroTagline: true})
		if strings.Contains(output, "$1") || strings.Contains(output, "$2") || strings.Contains(output, "$3") {
			t.Fatalf("theme %s leaked fastfetch color markers", selected.ID)
		}
		if !strings.Contains(output, selected.Tagline) {
			t.Fatalf("theme %s did not render its tagline", selected.ID)
		}
	}
}
