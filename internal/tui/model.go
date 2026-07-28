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
	saved    config.Config
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
	return Model{items: items, info: info, cfg: cfg, saved: cfg, cursor: cursor, width: 110, height: 32}
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
			m.refreshDirty()
		case "m":
			if m.cfg.Mode == "live" {
				m.cfg.Mode = "snapshot"
			} else {
				m.cfg.Mode = "live"
			}
			m.status = "Question mode: " + m.cfg.Mode
			m.refreshDirty()
		case "d":
			m.cfg.DistroJokes = !m.cfg.DistroJokes
			m.status = fmt.Sprintf("Distro jokes: %s", onOff(m.cfg.DistroJokes))
			m.refreshDirty()
		case "enter", "s":
			m.cfg.Theme = m.current().ID
			if err := config.Save(m.cfg); err != nil {
				m.status = "Save failed: " + err.Error()
			} else {
				m.status = "Saved " + m.current().ID
				m.saved = m.cfg
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
	mode := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#22d3ee")).Render("question: " + m.cfg.Mode)
	jokes := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f59e0b")).Render("distro jokes: " + onOff(m.cfg.DistroJokes))
	header := title + "  " + subtitle + "  " + mode + "  " + jokes
	if m.width < 110 {
		header = title + "\n" + subtitle + "\n" + mode + "  " + jokes
	}

	if m.width < 100 {
		panelWidth := max(30, m.width-6)
		panelHeight := max(7, (m.height-8)/2)
		body := m.renderList(panelWidth, panelHeight) + "\n" + m.renderPreview(panelWidth, panelHeight)
		return background.Render(header + "\n\n" + body + "\n" + m.help())
	}
	layoutWidth := min(144, m.width-4)
	listWidth := 34
	if layoutWidth < 100 {
		listWidth = 25
	}
	previewWidth := max(48, layoutWidth-listWidth-3)
	panelHeight := m.panelHeight()
	left := m.renderList(listWidth, panelHeight)
	right := m.renderPreview(previewWidth, panelHeight)
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, "   ", right)
	return background.Render(header + "\n\n" + body + "\n" + m.help())
}

func (m Model) renderList(width, height int) string {
	border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#475569")).Width(width).Height(height)
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

func (m Model) renderPreview(width, height int) string {
	item := m.current()
	showTagline := item.Joke || m.cfg.DistroJokes
	content := render.Fetch(item, m.info, render.Options{Width: width - 4, Color: true, Frame: m.frame, HideTagline: true})
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	limit := max(3, height-4)
	if showTagline {
		limit = max(3, limit-2)
	}
	if len(lines) > limit {
		lines = lines[:limit]
	}
	label := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(item.Primary)).Render(item.Name)
	meta := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748b")).Render("  " + item.ID)
	preview := label + meta + "\n\n" + strings.Join(lines, "\n")
	if showTagline {
		tagline := lipgloss.NewStyle().Foreground(lipgloss.Color(item.Accent)).Render(item.Tagline)
		preview += "\n\n" + tagline
	}
	border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(item.Primary)).Width(width).Height(height)
	return border.Padding(0, 1).Render(preview)
}

func (m Model) help() string {
	help := "↑/↓ navigate  enter/s save  r random  a animation  m question mode  d distro jokes  q quit"
	if m.width < 100 {
		help = "↑/↓ navigate  enter/s save  r random  q quit\na animation  m question mode  d distro jokes"
	}
	if m.status != "" {
		help += "   • " + m.status
	} else if m.dirty {
		help += "   • unsaved"
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#94a3b8")).Render(help)
}

func (m *Model) move(delta int) {
	m.cursor = max(0, min(len(m.items)-1, m.cursor+delta))
	m.refreshDirty()
	m.status = ""
	m.clampOffset()
}

func (m *Model) refreshDirty() {
	m.dirty = m.saved.Theme != m.current().ID || m.saved.Animation != m.cfg.Animation || m.saved.Mode != m.cfg.Mode || m.saved.DistroJokes != m.cfg.DistroJokes
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
	if m.width < 100 {
		return max(3, max(7, (m.height-8)/2)-4)
	}
	return max(8, m.panelHeight()-4)
}

func (m Model) panelHeight() int {
	return max(14, min(28, m.height-7))
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
