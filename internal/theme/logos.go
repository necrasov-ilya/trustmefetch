package theme

import (
	"embed"
	"strings"
)

// The distro artwork in ascii/ comes from fastfetch and is distributed under
// its MIT license. See THIRD_PARTY_NOTICES.md at the repository root.
//
//go:embed ascii/*.txt
var logoFiles embed.FS

func Logo(id string) []string {
	data, err := logoFiles.ReadFile("ascii/" + id + ".txt")
	if err != nil {
		data, _ = logoFiles.ReadFile("ascii/tux.txt")
	}
	return splitLogo(string(data))
}

func splitLogo(value string) []string {
	return strings.Split(strings.Trim(value, "\n"), "\n")
}
