package render

import (
	"fmt"
	"math"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/necrasov-ilya/trustmefetch/internal/system"
	"github.com/necrasov-ilya/trustmefetch/internal/theme"
)

type Options struct {
	Width             int
	Color             bool
	Frame             int
	ShowDistroTagline bool
}

func Fetch(selected theme.Theme, info system.Info, options Options) string {
	if options.Width <= 0 {
		options.Width = 100
	}
	logo := theme.Logo(selected.Logo)
	rows := dataRows(selected, info)
	primary := style(selected.Primary, true, options.Color)
	secondary := style(selected.Secondary, false, options.Color)
	accent := style(selected.Accent, false, options.Color)

	header := primary.Render(info.Identity())
	separator := primary.Render(strings.Repeat("-", max(20, lipgloss.Width(info.Identity()))))
	data := []string{header, separator}
	for _, row := range rows {
		label := primary.Render(row[0])
		data = append(data, label+accent.Render(": ")+secondary.Render(row[1]))
	}
	if selected.Joke || options.ShowDistroTagline {
		data = append(data, "", accent.Render(selected.Tagline))
	}
	if options.Color {
		data = append(data, "", colorBar(false), colorBar(true))
	}

	if options.Width < 76 {
		coloredLogo := colorLogo(logo, selected, options)
		return strings.Join(append(append(coloredLogo, ""), data...), "\n") + "\n"
	}

	leftWidth := 0
	for _, line := range logo {
		leftWidth = max(leftWidth, lipgloss.Width(stripLogoCodes(line)))
	}
	leftWidth += 5
	height := max(len(logo), len(data))
	lines := make([]string, 0, height)
	for index := 0; index < height; index++ {
		left := ""
		if index < len(logo) {
			left = paintLogoLine(logo[index], selected, options, index)
		}
		left = padVisible(left, leftWidth)
		right := ""
		if index < len(data) {
			right = data[index]
		}
		lines = append(lines, left+right)
	}
	return strings.Join(lines, "\n") + "\n"
}

func dataRows(selected theme.Theme, info system.Info) [][2]string {
	rows := [][2]string{
		{"OS", selected.Distro + " " + info.Arch},
		{"Host", info.Host},
		{"Kernel", "Darwin " + info.Kernel},
		{"Uptime", info.Uptime},
		{"Packages", info.Packages},
		{"Shell", info.Shell},
		{"Display", info.Resolution},
		{"WM", info.WM},
		{"WM Theme", info.WMTheme},
		{"Theme", selected.Desktop + " / " + info.UITheme},
		{"Font", info.Font},
		{"Cursor", info.Cursor},
		{"Terminal", info.Terminal},
		{"CPU", info.CPU},
		{"GPU", info.GPU},
		{"Memory", info.Memory},
		{"Swap", info.Swap},
		{"Disk", info.Disk},
		{"Local IP", info.LocalIP},
		{"Battery", info.Battery},
		{"Power Adapter", info.Power},
		{"Locale", info.Locale},
	}
	return rows
}

func colorLogo(lines []string, selected theme.Theme, options Options) []string {
	result := make([]string, len(lines))
	for index, line := range lines {
		result[index] = paintLogoLine(line, selected, options, index)
	}
	return result
}

func paintLogoLine(line string, selected theme.Theme, options Options, index int) string {
	if selected.Rainbow {
		color := rainbowHex((options.Frame*19 + index*31) % 360)
		return style(color, true, options.Color).Render(stripLogoCodes(line))
	}
	if !strings.Contains(line, "$") {
		return style(selected.Primary, true, options.Color).Render(line)
	}

	colors := []string{selected.Primary, selected.Accent, selected.Secondary}
	current := colors[0]
	var result strings.Builder
	start := 0
	for position := 0; position+1 < len(line); position++ {
		if line[position] != '$' || line[position+1] < '1' || line[position+1] > '9' {
			continue
		}
		if position > start {
			result.WriteString(style(current, true, options.Color).Render(line[start:position]))
		}
		colorIndex := int(line[position+1]-'1') % len(colors)
		current = colors[colorIndex]
		position++
		start = position + 1
	}
	if start < len(line) {
		result.WriteString(style(current, true, options.Color).Render(line[start:]))
	}
	return result.String()
}

func stripLogoCodes(line string) string {
	var result strings.Builder
	for position := 0; position < len(line); position++ {
		if line[position] == '$' && position+1 < len(line) && line[position+1] >= '1' && line[position+1] <= '9' {
			position++
			continue
		}
		result.WriteByte(line[position])
	}
	return result.String()
}

func style(color string, bold bool, enabled bool) lipgloss.Style {
	value := lipgloss.NewStyle()
	if !enabled {
		return value
	}
	value = value.Foreground(lipgloss.Color(color))
	if bold {
		value = value.Bold(true)
	}
	return value
}

func rainbowHex(hue int) string {
	h := float64((hue%360+360)%360) / 60
	x := 1 - math.Abs(math.Mod(h, 2)-1)
	var r, g, b float64
	switch int(h) {
	case 0:
		r, g = 1, x
	case 1:
		r, g = x, 1
	case 2:
		g, b = 1, x
	case 3:
		g, b = x, 1
	case 4:
		r, b = x, 1
	default:
		r, b = 1, x
	}
	return fmt.Sprintf("#%02x%02x%02x", int(r*255), int(g*255), int(b*255))
}

func padVisible(value string, width int) string {
	visible := lipgloss.Width(value)
	if visible >= width {
		return value
	}
	return value + strings.Repeat(" ", width-visible)
}

func colorBar(bright bool) string {
	var result strings.Builder
	start := 0
	if bright {
		start = 8
	}
	for index := 0; index < 8; index++ {
		result.WriteString(lipgloss.NewStyle().Background(lipgloss.Color(fmt.Sprintf("%d", start+index))).Render("   "))
	}
	return result.String()
}
