package output

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	CyanStyle    lipgloss.Style
	YellowStyle  lipgloss.Style
	GreenStyle   lipgloss.Style
	MagentaStyle lipgloss.Style
	DimStyle     lipgloss.Style
	ErrorStyle   lipgloss.Style
	WarnStyle    lipgloss.Style
	BoldStyle    lipgloss.Style

	// NoColor indicates if color output is disabled
	NoColor bool
)

func init() {
	NoColor = isNoColor()
	initStyles()
}

func isNoColor() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return true
	}
	return false
}

// SetNoColor forces color on or off.
func SetNoColor(v bool) {
	NoColor = v
	initStyles()
}

func initStyles() {
	if NoColor {
		CyanStyle = lipgloss.NewStyle()
		YellowStyle = lipgloss.NewStyle()
		GreenStyle = lipgloss.NewStyle()
		MagentaStyle = lipgloss.NewStyle()
		DimStyle = lipgloss.NewStyle()
		ErrorStyle = lipgloss.NewStyle()
		WarnStyle = lipgloss.NewStyle()
		BoldStyle = lipgloss.NewStyle()
		return
	}

	CyanStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	YellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	GreenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	MagentaStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	DimStyle = lipgloss.NewStyle().Faint(true)
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	WarnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	BoldStyle = lipgloss.NewStyle().Bold(true)
}

// Cyan renders text in cyan.
func Cyan(s string) string { return CyanStyle.Render(s) }

// Yellow renders text in yellow.
func Yellow(s string) string { return YellowStyle.Render(s) }

// Green renders text in green.
func Green(s string) string { return GreenStyle.Render(s) }

// Magenta renders text in magenta.
func Magenta(s string) string { return MagentaStyle.Render(s) }

// Dim renders text dimly.
func Dim(s string) string { return DimStyle.Render(s) }

// Bold renders text bold.
func Bold(s string) string { return BoldStyle.Render(s) }

// ProgressBar renders an ASCII progress bar.
func ProgressBar(percent, width int) string {
	filled := width * percent / 100
	if filled > width {
		filled = width
	}
	empty := width - filled
	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}
