package live

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/necrasov-ilya/trustmefetch/internal/render"
	"github.com/necrasov-ilya/trustmefetch/internal/system"
	"github.com/necrasov-ilya/trustmefetch/internal/theme"
)

type refreshDueMsg time.Time
type animationMsg time.Time
type statsMsg system.Info

type Model struct {
	info        system.Info
	themes      []theme.Theme
	selected    int
	animation   bool
	distroJokes bool
	question    bool
	paused      bool
	width       int
	height      int
	offset      int
	frame       int
	updated     time.Time
	quitting    bool
}

func New(info system.Info, selected theme.Theme, animation, distroJokes, question bool) Model {
	items := theme.All()
	index := 0
	for position, item := range items {
		if item.ID == selected.ID {
			index = position
			break
		}
	}
	return Model{info: info, themes: items, selected: index, animation: animation, distroJokes: distroJokes, question: question, width: 100, height: 30, updated: time.Now()}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(scheduleRefresh(), scheduleAnimation())
}

func scheduleRefresh() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return refreshDueMsg(t) })
}

func scheduleAnimation() tea.Cmd {
	return tea.Tick(75*time.Millisecond, func(t time.Time) tea.Msg { return animationMsg(t) })
}

func collect(info system.Info) tea.Cmd {
	return func() tea.Msg { return statsMsg(system.CollectDynamic(info)) }
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.clampOffset()
	case refreshDueMsg:
		if m.paused {
			return m, scheduleRefresh()
		}
		return m, tea.Batch(scheduleRefresh(), collect(m.info))
	case statsMsg:
		m.info = system.Info(msg)
		m.updated = time.Now()
		m.clampOffset()
		return m, nil
	case animationMsg:
		if m.animation && m.current().Rainbow && !m.paused {
			m.frame++
		}
		return m, scheduleAnimation()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "space":
			m.paused = !m.paused
		case "r":
			return m, collect(m.info)
		case "t", "tab":
			m.selected = (m.selected + 1) % len(m.themes)
			m.offset = 0
		case "shift+tab":
			m.selected = (m.selected - 1 + len(m.themes)) % len(m.themes)
			m.offset = 0
		case "up", "k":
			m.offset--
			m.clampOffset()
		case "down", "j":
			m.offset++
			m.clampOffset()
		case "home", "g":
			m.offset = 0
		case "end", "G":
			m.offset = m.maxOffset()
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}
	selected := m.current()
	status := "LIVE"
	if m.paused {
		status = "PAUSED"
	}
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(selected.Primary))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748b"))
	header := headerStyle.Render("trustmefetch "+status) + dim.Render("  "+selected.Name+"  updated "+m.updated.Format("15:04:05"))
	if m.question {
		answer := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(selected.Accent)).Render("YES · 100% LINUX")
		header = answer + "  " + header
	}

	content := render.Fetch(selected, m.info, render.Options{Width: max(40, m.width-4), Color: true, Frame: m.frame, ShowDistroTagline: m.distroJokes, Question: m.question})
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	available := max(3, m.height-4)
	end := min(len(lines), m.offset+available)
	body := strings.Join(lines[m.offset:end], "\n")
	footer := dim.Render(fmt.Sprintf("q quit  space pause  r refresh  t next theme  ↑/↓ scroll   %d/%d", m.offset+1, max(1, len(lines))))

	view := tea.NewView("\n  " + header + "\n\n  " + strings.ReplaceAll(body, "\n", "\n  ") + "\n  " + footer)
	view.AltScreen = true
	view.WindowTitle = "trustmefetch live"
	return view
}

func (m Model) current() theme.Theme {
	return m.themes[m.selected]
}

func (m Model) maxOffset() int {
	selected := m.current()
	content := render.Fetch(selected, m.info, render.Options{Width: max(40, m.width-4), Color: false, Frame: m.frame, ShowDistroTagline: m.distroJokes, Question: m.question})
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	return max(0, len(lines)-max(3, m.height-4))
}

func (m *Model) clampOffset() {
	m.offset = max(0, min(m.offset, m.maxOffset()))
}
