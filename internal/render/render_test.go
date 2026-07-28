package render

import (
	"strings"
	"testing"

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
