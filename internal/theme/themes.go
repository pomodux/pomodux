package theme

import "github.com/charmbracelet/lipgloss"

// DefaultTheme returns the default theme
func DefaultTheme() *Theme {
	return &Theme{
		Name: "default",
		Colors: Colors{
			Background:      lipgloss.Color("0"),
			Foreground:      lipgloss.Color("15"),
			Primary:         lipgloss.Color("6"),  // Cyan
			Secondary:       lipgloss.Color("5"),  // Magenta
			Success:         lipgloss.Color("2"),  // Green
			Warning:         lipgloss.Color("3"),  // Yellow
			Error:           lipgloss.Color("1"),  // Red
			Border:          lipgloss.Color("8"),  // Dark gray
			ProgressFilled:  lipgloss.Color("6"),  // Cyan
			ProgressEmpty:   lipgloss.Color("8"),  // Dark gray
			TextMuted:       lipgloss.Color("8"),  // Dark gray
		},
		Progress: ProgressStyle{
			FilledChar:     "█",
			EmptyChar:      "░",
			ShowPercentage: true,
		},
		Border: BorderStyle{
			Style: "rounded",
		},
	}
}

// NordTheme returns the Nord color scheme theme
func NordTheme() *Theme {
	return &Theme{
		Name: "nord",
		Colors: Colors{
			Background:      lipgloss.Color("#2e3440"),
			Foreground:      lipgloss.Color("#d8dee9"),
			Primary:         lipgloss.Color("#88c0d0"), // Blue
			Secondary:       lipgloss.Color("#81a1c1"), // Light blue
			Success:         lipgloss.Color("#a3be8c"), // Green
			Warning:         lipgloss.Color("#ebcb8b"), // Yellow
			Error:           lipgloss.Color("#bf616a"), // Red
			Border:          lipgloss.Color("#4c566a"), // Gray
			ProgressFilled:  lipgloss.Color("#88c0d0"), // Blue
			ProgressEmpty:   lipgloss.Color("#3b4252"), // Dark blue
			TextMuted:       lipgloss.Color("#6c7086"), // Muted gray
		},
		Progress: ProgressStyle{
			FilledChar:     "█",
			EmptyChar:      "░",
			ShowPercentage: true,
		},
		Border: BorderStyle{
			Style: "rounded",
		},
	}
}

// CatppuccinMochaTheme returns the Catppuccin Mocha color scheme theme
func CatppuccinMochaTheme() *Theme {
	return &Theme{
		Name: "catppuccin-mocha",
		Colors: Colors{
			Background:      lipgloss.Color("#1e1e2e"),
			Foreground:      lipgloss.Color("#cdd6f4"),
			Primary:         lipgloss.Color("#89b4fa"), // Blue
			Secondary:       lipgloss.Color("#cba6f7"), // Mauve
			Success:         lipgloss.Color("#a6e3a1"), // Green
			Warning:         lipgloss.Color("#f9e2af"), // Yellow
			Error:           lipgloss.Color("#f38ba8"), // Red
			Border:          lipgloss.Color("#45475a"), // Gray
			ProgressFilled:  lipgloss.Color("#89b4fa"), // Blue
			ProgressEmpty:   lipgloss.Color("#313244"), // Surface
			TextMuted:       lipgloss.Color("#6c7086"), // Muted
		},
		Progress: ProgressStyle{
			FilledChar:     "█",
			EmptyChar:      "░",
			ShowPercentage: true,
		},
		Border: BorderStyle{
			Style: "rounded",
		},
	}
}


