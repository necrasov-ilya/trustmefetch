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

var customLogos = map[string]string{
	"linux": `        ___
    .--'   '--.
   /  100%     \
  |   LINUX!    |
  |  .------.   |
   \ '------'  /
    '--.____.--'`,
	"sudo": `   ┌─────────────┐
   │ $ sudo      │
   │   believe   │
   │             │
   │  [granted]  │
   └─────────────┘`,
	"apple_glitch": "         .:'\n" +
		"      __ :'__\n" +
		"   .'`  `-'  `.\n" +
		"  :             :\n" +
		"  :  NOT macOS  :\n" +
		"   :           :\n" +
		"    `.___ __.'",
	"wayland": `   __      __
   \ \ /\ / /
    \ V  V /
     \ /\ /
      V  V
    WAYLAND
   PRO MAX`,
}

func Logo(id string) []string {
	if logo, ok := customLogos[id]; ok {
		return splitLogo(logo)
	}
	data, err := logoFiles.ReadFile("ascii/" + id + ".txt")
	if err != nil {
		data, _ = logoFiles.ReadFile("ascii/tux.txt")
	}
	return splitLogo(string(data))
}

func splitLogo(value string) []string {
	return strings.Split(strings.Trim(value, "\n"), "\n")
}
