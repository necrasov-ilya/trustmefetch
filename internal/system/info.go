package system

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Info struct {
	User       string
	Hostname   string
	OSVersion  string
	Build      string
	Host       string
	Kernel     string
	Arch       string
	Uptime     string
	Packages   string
	Shell      string
	Terminal   string
	CPU        string
	CPUUsage   string
	GPU        string
	Memory     string
	Disk       string
	Battery    string
	Resolution string
}

func Collect() Info {
	info := Info{
		User:      currentUser(),
		Hostname:  hostname(),
		OSVersion: command("sw_vers", "-productVersion"),
		Build:     command("sw_vers", "-buildVersion"),
		Host:      sysctl("hw.model"),
		Kernel:    command("uname", "-r"),
		Arch:      command("uname", "-m"),
		Uptime:    uptime(),
		Packages:  packages(),
		Shell:     shell(),
		Terminal:  terminal(),
		CPU:       sysctl("machdep.cpu.brand_string"),
		CPUUsage:  cpuUsage(),
		Memory:    memory(),
		Disk:      disk(),
		Battery:   battery(),
	}
	info.GPU, info.Resolution = displayInfo()
	if info.CPU == "" {
		info.CPU = runtime.GOARCH
	}
	if info.Host == "" {
		info.Host = "Apple computer"
	}
	return info
}

func CollectDynamic(info Info) Info {
	info.Uptime = uptime()
	info.CPUUsage = cpuUsage()
	info.Memory = memory()
	info.Disk = disk()
	info.Battery = battery()
	return info
}

func (i Info) Identity() string {
	if i.User == "" {
		return i.Hostname
	}
	return i.User + "@" + i.Hostname
}

func command(name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func sysctl(key string) string {
	return command("sysctl", "-n", key)
}

func currentUser() string {
	if value := os.Getenv("USER"); value != "" {
		return value
	}
	u, err := user.Current()
	if err == nil {
		return u.Username
	}
	return "user"
}

func hostname() string {
	value, err := os.Hostname()
	if err != nil {
		return "mac"
	}
	return strings.TrimSuffix(value, ".local")
}

func shell() string {
	value := os.Getenv("SHELL")
	if value == "" {
		return "unknown"
	}
	return filepath.Base(value)
}

func terminal() string {
	for _, key := range []string{"TERM_PROGRAM", "LC_TERMINAL", "TERM"} {
		if value := os.Getenv(key); value != "" {
			if version := os.Getenv("TERM_PROGRAM_VERSION"); version != "" && key == "TERM_PROGRAM" {
				return value + " " + version
			}
			return value
		}
	}
	return "unknown"
}

func uptime() string {
	raw := sysctl("kern.boottime")
	match := regexp.MustCompile(`sec = ([0-9]+)`).FindStringSubmatch(raw)
	if len(match) != 2 {
		return "unknown"
	}
	seconds, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return "unknown"
	}
	d := time.Since(time.Unix(seconds, 0))
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	parts := make([]string, 0, 3)
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	parts = append(parts, fmt.Sprintf("%dm", minutes))
	return strings.Join(parts, " ")
}

func packages() string {
	if _, err := exec.LookPath("brew"); err != nil {
		return "0 (brew not found)"
	}
	formulae := lines(command("brew", "list", "--formula"))
	casks := lines(command("brew", "list", "--cask"))
	if casks == 0 {
		return fmt.Sprintf("%d (brew)", formulae)
	}
	return fmt.Sprintf("%d (brew), %d (brew-cask)", formulae, casks)
}

func cpuUsage() string {
	raw := command("ps", "-A", "-o", "%cpu=")
	if raw == "" {
		return "unknown"
	}
	total := 0.0
	for _, field := range strings.Fields(raw) {
		value, err := strconv.ParseFloat(field, 64)
		if err == nil {
			total += value
		}
	}
	if cores := runtime.NumCPU(); cores > 0 {
		total /= float64(cores)
	}
	if total > 100 {
		total = 100
	}
	return fmt.Sprintf("%.1f%%", total)
}

func lines(value string) int {
	if strings.TrimSpace(value) == "" {
		return 0
	}
	return len(strings.Split(strings.TrimSpace(value), "\n"))
}

func memory() string {
	total, _ := strconv.ParseUint(sysctl("hw.memsize"), 10, 64)
	if total == 0 {
		return "unknown"
	}
	pageSize := uint64(4096)
	freePages := uint64(0)
	for _, line := range strings.Split(command("vm_stat"), "\n") {
		if strings.Contains(line, "page size of") {
			if match := regexp.MustCompile(`page size of ([0-9]+) bytes`).FindStringSubmatch(line); len(match) == 2 {
				pageSize, _ = strconv.ParseUint(match[1], 10, 64)
			}
		}
		if strings.HasPrefix(line, "Pages free:") || strings.HasPrefix(line, "Pages inactive:") || strings.HasPrefix(line, "Pages speculative:") {
			fields := strings.Fields(strings.TrimSuffix(line, "."))
			if len(fields) > 0 {
				value, _ := strconv.ParseUint(fields[len(fields)-1], 10, 64)
				freePages += value
			}
		}
	}
	available := freePages * pageSize
	if available > total {
		available = 0
	}
	used := total - available
	return fmt.Sprintf("%.1f GiB / %.1f GiB", gib(used), gib(total))
}

func disk() string {
	path := "/"
	if _, err := os.Stat("/System/Volumes/Data"); err == nil {
		path = "/System/Volumes/Data"
	}
	out := command("df", "-k", path)
	rows := strings.Split(out, "\n")
	if len(rows) < 2 {
		return "unknown"
	}
	fields := strings.Fields(rows[len(rows)-1])
	if len(fields) < 5 {
		return "unknown"
	}
	total, _ := strconv.ParseUint(fields[1], 10, 64)
	used, _ := strconv.ParseUint(fields[2], 10, 64)
	return fmt.Sprintf("%.1f GiB / %.1f GiB (%s)", float64(used)/1024/1024, float64(total)/1024/1024, fields[4])
}

func battery() string {
	raw := command("pmset", "-g", "batt")
	match := regexp.MustCompile(`([0-9]+)%`).FindStringSubmatch(raw)
	if len(match) != 2 {
		return "not present"
	}
	state := "battery"
	if strings.Contains(raw, "AC Power") || strings.Contains(raw, "charging") || strings.Contains(raw, "charged") {
		state = "AC power"
	}
	return match[1] + "% (" + state + ")"
}

func displayInfo() (string, string) {
	raw := command("system_profiler", "SPDisplaysDataType", "-json")
	if raw == "" {
		return "unknown", "unknown"
	}
	var value any
	if json.Unmarshal([]byte(raw), &value) != nil {
		return "unknown", "unknown"
	}
	models := findStrings(value, "sppci_model")
	resolutions := findStrings(value, "_spdisplays_resolution")
	return firstOr(models, "unknown"), firstOr(resolutions, "unknown")
}

func findStrings(value any, key string) []string {
	var found []string
	var walk func(any)
	walk = func(node any) {
		switch node := node.(type) {
		case map[string]any:
			for k, child := range node {
				if k == key {
					if text, ok := child.(string); ok && text != "" {
						found = append(found, text)
					}
				}
				walk(child)
			}
		case []any:
			for _, child := range node {
				walk(child)
			}
		}
	}
	walk(value)
	return found
}

func firstOr(values []string, fallback string) string {
	if len(values) == 0 {
		return fallback
	}
	return values[0]
}

func gib(value uint64) float64 {
	return float64(value) / 1024 / 1024 / 1024
}
