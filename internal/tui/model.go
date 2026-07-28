package tui

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/necrasov-ilya/trustmefetch/internal/config"
	"github.com/necrasov-ilya/trustmefetch/internal/render"
	"github.com/necrasov-ilya/trustmefetch/internal/system"
	"github.com/necrasov-ilya/trustmefetch/internal/theme"
)

type tickMsg time.Time

type Model struct {
	items    []theme.Theme
	info     system.Info
	cfg      config.Config
	cursor   int
	offset   int
	width    int
	height   int
	frame    int
	status   string
	dirty    bool
	quitting bool
}

func New(info system.Info, cfg config.Config) Model {
	items := theme.All()
	cursor := 0
	for index, item := range items {
		if item.ID == cfg.Theme {
			cursor = index
			break
		}
	}
	return Model{items: items, info: info, cfg: cfg, cursor: cursor, width: 110, height: 32}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(70*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.clampOffset()
	case tickMsg:
		if m.current().Rainbow && m.cfg.Animation {
			m.frame++
		}
		return m, tick()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			m.move(-1)
		case "down", "j":
			m.move(1)
		case "pgup":
			m.move(-m.visibleRows())
		case "pgdown":
			m.move(m.visibleRows())
		case "home", "g":
			m.cursor = 0
			m.clampOffset()
		case "end", "G":
			m.cursor = len(m.items) - 1
			m.clampOffset()
		case "r":
			m.cursor = rand.IntN(len(m.items))
			m.status = "Random reality selected"
			m.dirty = true
			m.clampOffset()
		case "a":
			m.cfg.Animation = !m.cfg.Animation
			m.status = fmt.Sprintf("Animation: %s", onOff(m.cfg.Animation))
			m.dirty = true
		case "enter", "s":
			m.cfg.Theme = m.current().ID
			if err := config.Save(m.cfg); err != nil {
				m.status = "Save failed: " + err.Error()
			} else {
				m.status = "Saved " + m.current().ID
				m.dirty = false
			}
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}
	canvas := m.renderView()
	view := tea.NewView(canvas)
	view.AltScreen = true
	view.WindowTitle = "trustmefetch config"
	return view
}

func (m Model) renderView() string {
	background := lipgloss.NewStyle().Padding(1, 2)
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#a78bfa")).Render("trustmefetch config")
	subtitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#94a3b8")).Render("Choose the Linux your Mac believes it is today")
	header := title + "  " + subtitle

	if m.width < 70 {
		body := m.renderList(max(30, m.width-6)) + "\n\n" + m.renderPreview(max(42, m.width-6))
		return background.Render(header + "\n\n" + body + "\n" + m.help())
	}
	listWidth := 34
	if m.width < 100 {
		listWidth = 25
	}
	previewWidth := max(48, m.width-listWidth-10)
	left := m.renderList(listWidth)
	right := m.renderPreview(previewWidth)
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, "   ", right)
	return background.Render(header + "\n\n" + body + "\n" + m.help())
}

func (m Model) renderList(width int) string {
	border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#475569")).Width(width)
	selected := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0f172a")).Background(lipgloss.Color("#a78bfa"))
	normal := lipgloss.NewStyle().Foreground(lipgloss.Color("#cbd5e1"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748b"))
	rows := []string{lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("Themes  %d/%d", m.cursor+1, len(m.items))), ""}
	end := min(len(m.items), m.offset+m.visibleRows())
	for index := m.offset; index < end; index++ {
		item := m.items[index]
		kind := "DISTRO"
		if item.Joke {
			kind = "JOKE"
		}
		nameWidth := max(8, width-12)
		line := fmt.Sprintf(" %-*s %-6s", nameWidth, truncate(item.Name, nameWidth), kind)
		if index == m.cursor {
			rows = append(rows, selected.Width(width-2).Render(line))
		} else {
			rows = append(rows, normal.Render(fmt.Sprintf(" %-*s ", nameWidth, truncate(item.Name, nameWidth)))+dim.Render(fmt.Sprintf("%-6s", kind)))
		}
	}
	return border.Padding(0, 1).Render(strings.Join(rows, "\n"))
}

func (m Model) renderPreview(width int) string {
	item := m.current()
	content := render.Fetch(item, m.info, render.Options{Width: width - 4, Color: true, Frame: m.frame})
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	limit := max(10, m.height-10)
	if len(lines) > limit {
		lines = lines[:limit]
	}
	label := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(item.Primary)).Render(item.Name)
	meta := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748b")).Render("  " + item.ID)
	border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(item.Primary)).Width(width)
	return border.Padding(0, 1).Render(label + meta + "\n\n" + strings.Join(lines, "\n"))
}

func (m Model) help() string {
	help := "↑/↓ navigate  enter/s save  r random  a animation  q quit"
	if m.status != "" {
		help += "   • " + m.status
	} else if m.dirty {
		help += "   • unsaved"
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#94a3b8")).Render(help)
}

func (m *Model) move(delta int) {
	m.cursor = max(0, min(len(m.items)-1, m.cursor+delta))
	m.dirty = m.cfg.Theme != m.current().ID
	m.status = ""
	m.clampOffset()
}

func (m *Model) clampOffset() {
	visible := m.visibleRows()
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}
	m.offset = max(0, min(m.offset, max(0, len(m.items)-visible)))
}

func (m Model) visibleRows() int {
	if m.width < 70 {
		return max(5, min(8, (m.height-12)/2))
	}
	return max(8, min(24, m.height-11))
}

func (m Model) current() theme.Theme {
	return m.items[m.cursor]
}

func onOff(value bool) string {
	if value {
		return "on"
	}
	return "off"
}

func truncate(value string, width int) string {
	if len([]rune(value)) <= width {
		return value
	}
	runes := []rune(value)
	return string(runes[:width-1]) + "…"
}
