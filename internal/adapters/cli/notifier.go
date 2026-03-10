package cli

import (
	"fmt"
	"strings"
)

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorGray    = "\033[90m"
	colorBold    = "\033[1m"
)

// Notifier implements ports.UserNotifier for CLI
type Notifier struct{}

// New creates a new CLI Notifier
func New() *Notifier { return &Notifier{} }

func (n *Notifier) Info(msg string) {
	fmt.Printf("%s%s ℹ%s  %s\n", colorBold, colorBlue, colorReset, msg)
}

func (n *Notifier) Success(msg string) {
	fmt.Printf("%s%s ✓%s  %s\n", colorBold, colorGreen, colorReset, msg)
}

func (n *Notifier) Warning(msg string) {
	fmt.Printf("%s%s ⚠%s  %s\n", colorBold, colorYellow, colorReset, msg)
}

func (n *Notifier) Error(msg string) {
	fmt.Printf("%s%s ✗%s  %s\n", colorBold, colorRed, colorReset, msg)
}

// Progress renders a visual progress bar.
//
//	[██████████░░░░░░░░░░] 50% Scanning packages...
func (n *Notifier) Progress(current, total int, msg string) {
	const barWidth = 20
	pct := 0
	filled := 0
	if total > 0 {
		pct = (current * 100) / total
		filled = (current * barWidth) / total
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	fmt.Printf("\r%s[%s]%s %s%3d%%%s  %s",
		colorCyan, bar, colorReset,
		colorBold, pct, colorReset,
		msg,
	)
	if current >= total {
		fmt.Println()
	}
}

// Confirm asks the user a yes/no question and returns true for y/Y/yes.
func (n *Notifier) Confirm(msg string) bool {
	fmt.Printf("%s%s ?%s  %s %s[y/N]%s: ",
		colorBold, colorYellow, colorReset,
		msg,
		colorGray, colorReset,
	)
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		// stdin closed or unavailable: default to "no"
		return false
	}
	return response == "y" || response == "Y" || response == "yes"
}

// PrintScanSummary renders a box-drawing summary table.
// Call this after a scan completes to display totals.
func (n *Notifier) PrintScanSummary(rows [][2]string) {
	if len(rows) == 0 {
		return
	}

	// Calculate column widths
	maxKey := 0
	maxVal := 0
	for _, row := range rows {
		if len(row[0]) > maxKey {
			maxKey = len(row[0])
		}
		if len(row[1]) > maxVal {
			maxVal = len(row[1])
		}
	}

	// Box drawing chars
	const (
		tl = "╭"
		tr = "╮"
		bl = "╰"
		br = "╯"
		h  = "─"
		v  = "│"
		lj = "├"
		rj = "┤"
	)

	width := maxKey + maxVal + 7 // padding + separator
	topBar := tl + strings.Repeat(h, width) + tr
	botBar := bl + strings.Repeat(h, width) + br

	fmt.Printf("\n%s%s%s\n", colorCyan, topBar, colorReset)
	fmt.Printf("%s%s%s  %s%sScan Results%s  %s%s%s\n",
		colorCyan, v, colorReset,
		colorBold, colorBlue, colorReset,
		colorCyan, v, colorReset,
	)
	fmt.Printf("%s%s%s%s%s\n", colorCyan, lj, strings.Repeat(h, width), rj, colorReset)
	for _, row := range rows {
		keyPad := strings.Repeat(" ", maxKey-len(row[0]))
		valPad := strings.Repeat(" ", maxVal-len(row[1]))
		fmt.Printf("%s%s%s  %s%s%s  %s%s%s  %s%s%s%s\n",
			colorCyan, v, colorReset,
			colorGray, row[0]+keyPad, colorReset,
			colorBold, row[1]+valPad, colorReset,
			colorCyan, v, colorReset, "",
		)
	}
	fmt.Printf("%s%s%s\n\n", colorCyan, botBar, colorReset)
}
