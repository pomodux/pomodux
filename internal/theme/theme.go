package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme represents a visual theme
type Theme struct {
	Name      string
	Colors    Colors
	Progress  ProgressStyle
	Border    BorderStyle
}

// Colors defines the color palette
type Colors struct {
	Background  lipgloss.Color
	Foreground  lipgloss.Color
	Primary     lipgloss.Color
	Secondary   lipgloss.Color
	Success     lipgloss.Color
	Warning     lipgloss.Color
	Error       lipgloss.Color
	Border      lipgloss.Color
	ProgressFilled lipgloss.Color
	ProgressEmpty  lipgloss.Color
	TextMuted   lipgloss.Color
}

// ProgressStyle defines progress bar styling
type ProgressStyle struct {
	FilledChar    string
	EmptyChar     string
	ShowPercentage bool
}

// BorderStyle defines border styling
type BorderStyle struct {
	Style string // rounded, square, double, none
}

// GetTheme returns a theme by name
func GetTheme(name string) *Theme {
	switch name {
	case "nord":
		return NordTheme()
	case "catppuccin-mocha":
		return CatppuccinMochaTheme()
	default:
		return DefaultTheme()
	}
}

// TitleStyle returns a style for titles
func (t *Theme) TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Colors.Primary).
		Bold(true)
}

// StatusStyle returns a style for status text
func (t *Theme) StatusStyle(status string) lipgloss.Style {
	var color lipgloss.Color
	switch status {
	case "running":
		color = t.Colors.Success
	case "paused":
		color = t.Colors.Warning
	case "completed":
		color = t.Colors.Success
	case "stopped":
		color = t.Colors.Error
	default:
		color = t.Colors.Foreground
	}

	return lipgloss.NewStyle().Foreground(color)
}

// ProgressFilledStyle returns a style for filled progress bar
func (t *Theme) ProgressFilledStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Colors.ProgressFilled)
}

// ProgressEmptyStyle returns a style for empty progress bar
func (t *Theme) ProgressEmptyStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Colors.ProgressEmpty)
}

// BorderStyle returns a border style
func (t *Theme) BorderStyle() lipgloss.Style {
	style := lipgloss.NewStyle().
		BorderForeground(t.Colors.Border)

	switch t.Border.Style {
	case "rounded":
		style = style.BorderStyle(lipgloss.RoundedBorder())
	case "square":
		style = style.BorderStyle(lipgloss.NormalBorder())
	case "double":
		style = style.BorderStyle(lipgloss.DoubleBorder())
	case "none":
		style = style.BorderStyle(lipgloss.HiddenBorder())
	default:
		style = style.BorderStyle(lipgloss.RoundedBorder())
	}

	return style
}

