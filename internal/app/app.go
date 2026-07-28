package app

import (
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	cterm "github.com/charmbracelet/x/term"
	"github.com/necrasov-ilya/trustmefetch/internal/config"
	"github.com/necrasov-ilya/trustmefetch/internal/live"
	"github.com/necrasov-ilya/trustmefetch/internal/render"
	"github.com/necrasov-ilya/trustmefetch/internal/system"
	"github.com/necrasov-ilya/trustmefetch/internal/theme"
	"github.com/necrasov-ilya/trustmefetch/internal/tui"
)

func Run(args []string, version string, stdout, stderr io.Writer) error {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(stderr, "warning: invalid config, using defaults:", err)
		cfg = config.Default()
	}

	if len(args) == 0 {
		return printFetch(cfg, "", stdout)
	}

	switch args[0] {
	case "--question":
		if cfg.Mode == "live" && writerIsTerminal(stdout) {
			return runLive(cfg, "", stdout)
		}
		return printFetch(cfg, "", stdout)
	case "live":
		id := ""
		if len(args) > 1 {
			id = args[1]
		}
		return runLive(cfg, id, stdout)
	case "config", "configure":
		if !writerIsTerminal(stdout) {
			return errors.New("config requires an interactive terminal")
		}
		program := tea.NewProgram(tui.New(system.Collect(), cfg))
		_, err := program.Run()
		return err
	case "themes":
		for _, item := range theme.All() {
			marker := " "
			if item.ID == cfg.Theme {
				marker = "*"
			}
			fmt.Fprintf(stdout, "%s %-20s %s\n", marker, item.ID, item.Name)
		}
		return nil
	case "theme":
		if len(args) != 2 {
			return errors.New("usage: trustmefetch theme <theme-id>")
		}
		if _, ok := theme.ByID(args[1]); !ok {
			return fmt.Errorf("unknown theme %q; run trustmefetch themes", args[1])
		}
		cfg.Theme = args[1]
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Fprintln(stdout, "Theme saved:", args[1])
		return nil
	case "mode":
		if len(args) != 2 || (args[1] != "snapshot" && args[1] != "live") {
			return errors.New("usage: trustmefetch mode <snapshot|live>")
		}
		cfg.Mode = args[1]
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Fprintln(stdout, "Question mode saved:", cfg.Mode)
		return nil
	case "preview":
		id := cfg.Theme
		if len(args) > 1 {
			id = args[1]
		}
		return printFetch(cfg, id, stdout)
	case "random":
		items := theme.All()
		cfg.Theme = items[rand.IntN(len(items))].ID
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Fprintln(stdout, "Reality randomized:", cfg.Theme)
		return printFetch(cfg, "", stdout)
	case "doctor":
		return doctor(stdout)
	case "version", "--version", "-v":
		fmt.Fprintln(stdout, "trustmefetch", version)
		return nil
	case "help", "--help", "-h":
		fmt.Fprint(stdout, usage)
		return nil
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usage)
	}
}

func runLive(cfg config.Config, requested string, stdout io.Writer) error {
	if !writerIsTerminal(stdout) {
		return errors.New("live mode requires an interactive terminal")
	}
	id := cfg.Theme
	if requested != "" {
		id = requested
	}
	selected, ok := theme.ByID(id)
	if !ok {
		return fmt.Errorf("unknown theme %q; run trustmefetch themes", id)
	}
	program := tea.NewProgram(live.New(system.Collect(), selected, cfg.Animation))
	_, err := program.Run()
	return err
}

func printFetch(cfg config.Config, requested string, stdout io.Writer) error {
	id := cfg.Theme
	if requested != "" {
		id = requested
	}
	selected, ok := theme.ByID(id)
	if !ok {
		selected = theme.Must(config.DefaultTheme)
	}
	info := system.Collect()
	color := writerIsTerminal(stdout) && os.Getenv("NO_COLOR") == ""
	width := terminalWidth(stdout)
	options := render.Options{Width: width, Color: color}
	if selected.Rainbow && cfg.Animation && color {
		return animate(stdout, selected, info, options)
	}
	_, err := io.WriteString(stdout, render.Fetch(selected, info, options))
	return err
}

func animate(output io.Writer, selected theme.Theme, info system.Info, options render.Options) error {
	lineCount := 0
	for frame := 0; frame < 14; frame++ {
		options.Frame = frame
		content := render.Fetch(selected, info, options)
		if frame > 0 {
			if _, err := fmt.Fprintf(output, "\x1b[%dA\r\x1b[J", lineCount); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(output, content); err != nil {
			return err
		}
		lineCount = strings.Count(content, "\n")
		time.Sleep(55 * time.Millisecond)
	}
	return nil
}

func doctor(output io.Writer) error {
	path, err := config.Path()
	if err != nil {
		path = "unavailable: " + err.Error()
	}
	checks := [][2]string{
		{"platform", systemPlatform()},
		{"architecture", systemArch()},
		{"config", path},
		{"terminal", os.Getenv("TERM_PROGRAM")},
		{"shell", os.Getenv("SHELL")},
	}
	for _, check := range checks {
		fmt.Fprintf(output, "%-14s %s\n", check[0]+":", fallback(check[1], "unknown"))
	}
	return nil
}

func writerIsTerminal(writer io.Writer) bool {
	file, ok := writer.(*os.File)
	return ok && cterm.IsTerminal(file.Fd())
}

func terminalWidth(writer io.Writer) int {
	if file, ok := writer.(*os.File); ok {
		if width, _, err := cterm.GetSize(file.Fd()); err == nil && width > 0 {
			return width
		}
	}
	if width, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && width > 0 {
		return width
	}
	return 100
}

func fallback(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func systemPlatform() string { return runtime.GOOS }
func systemArch() string     { return runtime.GOARCH }

const usage = `trustmefetch — Trust me, it's Linux.

Usage:
  trustmefetch                    Show the configured fetch
  trustmefetch config             Open the interactive theme picker
  trustmefetch live [theme-id]    Open the live system view
  trustmefetch themes             List all themes
  trustmefetch theme <id>         Save a theme
  trustmefetch mode <mode>        Set question mode: snapshot or live
  trustmefetch preview [id]       Preview a theme
  trustmefetch random             Pick and save a random theme
  trustmefetch doctor             Inspect the installation
  trustmefetch version            Print the version
`
