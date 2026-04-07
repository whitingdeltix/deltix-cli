package output

import (
	"fmt"
	"strings"
)

// Colors (ANSI)
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

func ScoreColor(score int) string {
	if score >= 80 {
		return Green
	} else if score >= 60 {
		return Yellow
	}
	return Red
}

func PassFail(passed bool) string {
	if passed {
		return Green + "PASS" + Reset
	}
	return Red + "FAIL" + Reset
}

func Status(status string) string {
	switch status {
	case "completed":
		return Green + status + Reset
	case "running", "queued":
		return Blue + status + Reset
	case "failed":
		return Red + status + Reset
	case "cancelled":
		return Dim + status + Reset
	default:
		return status
	}
}

func Header(text string) {
	fmt.Println()
	fmt.Println(Bold + text + Reset)
	fmt.Println(Dim + strings.Repeat("─", 50) + Reset)
}

func Section(label string) {
	fmt.Println()
	fmt.Println(Dim + strings.ToUpper(label) + Reset)
}

func Check(text string) {
	fmt.Println("  " + Green + "✓" + Reset + " " + text)
}

func Cross(text string) {
	fmt.Println("  " + Red + "✗" + Reset + " " + text)
}

func Warn(text string) {
	fmt.Println("  " + Yellow + "⚠" + Reset + " " + text)
}

func Info(text string) {
	fmt.Println("  " + Dim + text + Reset)
}

func Divider() {
	fmt.Println(Dim + strings.Repeat("─", 50) + Reset)
}
