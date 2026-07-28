package system

import (
	"context"
	"encoding/json"
	"fmt"
	"math/bits"
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
	Swap       string
	Disk       string
	Battery    string
	Resolution string
	WM         string
	WMTheme    string
	UITheme    string
	Font       string
	Cursor     string
	LocalIP    string
	Power      string
	Locale     string
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
		Swap:      swap(),
		Disk:      disk(),
		Battery:   battery(),
		WM:        "Quartz Compositor",
		WMTheme:   "Multicolor " + appearance(),
		UITheme:   "Liquid Glass",
		Font:      ".AppleSystemUIFont [System], Helvetica [User]",
		Cursor:    "Fill - Black, Outline - White (32px)",
		LocalIP:   localIP(),
		Power:     powerAdapter(),
		Locale:    localeName(),
	}
	info.Host, info.CPU = hardwareInfo(info.Host, info.CPU)
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
	info.Swap = swap()
	info.Disk = disk()
	info.Battery = battery()
	info.LocalIP = localIP()
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
	speculativePages := uint64(0)
	fileBackedPages := uint64(0)
	for _, line := range strings.Split(command("vm_stat"), "\n") {
		if strings.Contains(line, "page size of") {
			if match := regexp.MustCompile(`page size of ([0-9]+) bytes`).FindStringSubmatch(line); len(match) == 2 {
				pageSize, _ = strconv.ParseUint(match[1], 10, 64)
			}
		}
		fields := strings.Fields(strings.TrimSuffix(line, "."))
		if len(fields) == 0 {
			continue
		}
		value, _ := strconv.ParseUint(fields[len(fields)-1], 10, 64)
		switch {
		case strings.HasPrefix(line, "Pages free:"):
			freePages = value
		case strings.HasPrefix(line, "Pages speculative:"):
			speculativePages = value
		case strings.HasPrefix(line, "File-backed pages:"):
			fileBackedPages = value
		}
	}
	if speculativePages < freePages {
		freePages -= speculativePages
	}
	available := (freePages + fileBackedPages) * pageSize
	if available > total {
		available = 0
	}
	used := total - available
	return fmt.Sprintf("%.1f GiB / %.1f GiB", gib(used), gib(total))
}

func swap() string {
	raw := sysctl("vm.swapusage")
	match := regexp.MustCompile(`total = ([0-9.]+)M\s+used = ([0-9.]+)M`).FindStringSubmatch(raw)
	if len(match) != 3 {
		return "unknown"
	}
	total, _ := strconv.ParseFloat(match[1], 64)
	used, _ := strconv.ParseFloat(match[2], 64)
	return fmt.Sprintf("%.2f GiB / %.2f GiB", used/1024, total/1024)
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
	gpu := firstOr(models, "unknown")
	if cores := firstOr(findStrings(value, "sppci_cores"), ""); cores != "" {
		gpu += " (" + cores + " cores)"
	}
	displays := findMapsWithKey(value, "_spdisplays_pixels")
	if len(displays) == 0 {
		return gpu, "unknown"
	}
	display := displays[0]
	name := stringValue(display, "_name")
	pixels := strings.ReplaceAll(stringValue(display, "_spdisplays_pixels"), " ", "")
	scaled := stringValue(display, "_spdisplays_resolution")
	connection := "External"
	if stringValue(display, "spdisplays_connection_type") == "spdisplays_internal" {
		connection = "Built-in"
	}
	resolution := pixels
	if scaled != "" {
		resolution += " (" + scaled + ")"
	}
	if name != "" {
		resolution = name + ": " + resolution
	}
	return gpu, resolution + " [" + connection + "]"
}

func hardwareInfo(fallbackHost, fallbackCPU string) (string, string) {
	raw := command("system_profiler", "SPHardwareDataType", "-json")
	if raw == "" {
		return fallbackHost, fallbackCPU
	}
	var value any
	if json.Unmarshal([]byte(raw), &value) != nil {
		return fallbackHost, fallbackCPU
	}
	items := findMapsWithKey(value, "machine_model")
	if len(items) == 0 {
		return fallbackHost, fallbackCPU
	}
	item := items[0]
	model := stringValue(item, "machine_model")
	name := stringValue(item, "machine_name")
	host := fallbackHost
	if name != "" && model != "" {
		host = name + " (" + model + ")"
	}
	cpu := fallbackCPU
	if layout := stringValue(item, "number_processors"); layout != "" {
		match := regexp.MustCompile(`proc [0-9]+:([0-9]+):([0-9]+)`).FindStringSubmatch(layout)
		if len(match) == 3 {
			cpu += " (" + match[1] + "+" + match[2] + ")"
		}
	}
	return host, cpu
}

func localIP() string {
	route := command("route", "-n", "get", "default")
	match := regexp.MustCompile(`(?m)^\s*interface:\s+(\S+)`).FindStringSubmatch(route)
	if len(match) != 2 {
		return "unknown"
	}
	name := match[1]
	configuration := command("ifconfig", name)
	inet := regexp.MustCompile(`(?m)^\s*inet\s+(\S+)\s+netmask\s+(0x[0-9a-fA-F]+)`).FindStringSubmatch(configuration)
	if len(inet) != 3 {
		return "unknown"
	}
	mask, _ := strconv.ParseUint(strings.TrimPrefix(inet[2], "0x"), 16, 32)
	prefix := bits.OnesCount32(uint32(mask))
	return fmt.Sprintf("%s/%d (%s)", inet[1], prefix, name)
}

func powerAdapter() string {
	raw := command("system_profiler", "SPPowerDataType", "-json")
	if raw == "" {
		return "not connected"
	}
	var value any
	if json.Unmarshal([]byte(raw), &value) != nil {
		return "unknown"
	}
	items := findMapsWithKey(value, "sppower_ac_charger_name")
	if len(items) == 0 {
		return "not connected"
	}
	return stringValue(items[0], "sppower_ac_charger_name")
}

func appearance() string {
	if command("defaults", "read", "-g", "AppleInterfaceStyle") == "Dark" {
		return "Dark"
	}
	return "Light"
}

func localeName() string {
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return "unknown"
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

func findMapsWithKey(value any, key string) []map[string]any {
	var found []map[string]any
	var walk func(any)
	walk = func(node any) {
		switch node := node.(type) {
		case map[string]any:
			if _, ok := node[key]; ok {
				found = append(found, node)
			}
			for _, child := range node {
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

func stringValue(value map[string]any, key string) string {
	switch item := value[key].(type) {
	case string:
		return item
	case float64:
		return strconv.FormatFloat(item, 'f', -1, 64)
	default:
		return ""
	}
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
