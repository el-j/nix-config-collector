package cli

import "fmt"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// Notifier implements ports.UserNotifier for CLI
type Notifier struct{}

// New creates a new CLI Notifier
func New() *Notifier { return &Notifier{} }

func (n *Notifier) Info(msg string)    { fmt.Printf("%s[INFO]%s %s\n", colorBlue, colorReset, msg) }
func (n *Notifier) Success(msg string) { fmt.Printf("%s[OK]%s %s\n", colorGreen, colorReset, msg) }
func (n *Notifier) Warning(msg string) { fmt.Printf("%s[WARN]%s %s\n", colorYellow, colorReset, msg) }
func (n *Notifier) Error(msg string)   { fmt.Printf("%s[ERR]%s %s\n", colorRed, colorReset, msg) }

func (n *Notifier) Progress(current, total int, msg string) {
	pct := 0
	if total > 0 {
		pct = (current * 100) / total
	}
	fmt.Printf("%s[%d%%]%s %s\n", colorCyan, pct, colorReset, msg)
}

func (n *Notifier) Confirm(msg string) bool {
	fmt.Printf("%s[?]%s %s [y/N]: ", colorYellow, colorReset, msg)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y" || response == "yes"
}
