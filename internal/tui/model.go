package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pomodux/pomodux/internal/theme"
	"github.com/pomodux/pomodux/internal/timer"
)

const (
	minWidth  = 80
	minHeight = 24
)

// Model represents the TUI model
type Model struct {
	timer      *timer.Timer
	progress   progress.Model
	theme      *theme.Theme
	width      int
	height     int
	quitting   bool
	sessionID  string
	statePath  string
}

// NewModel creates a new TUI model. If theme is nil, the default theme is used.
func NewModel(t *timer.Timer, sessionID string, statePath string, th *theme.Theme) Model {
	if th == nil {
		th = theme.GetTheme("default")
	}

	fullRune := firstRune(th.Progress.FilledChar, '█')
	emptyRune := firstRune(th.Progress.EmptyChar, '░')

	opts := []progress.Option{
		progress.WithFillCharacters(fullRune, emptyRune),
		progress.WithSolidFill(string(th.Colors.ProgressFilled)),
		progress.WithWidth(76), // Updated on WindowSizeMsg
	}
	if !th.Progress.ShowPercentage {
		opts = append(opts, progress.WithoutPercentage())
	}

	prog := progress.New(opts...)
	prog.EmptyColor = string(th.Colors.ProgressEmpty)

	return Model{
		timer:     t,
		progress:  prog,
		theme:     th,
		sessionID: sessionID,
		statePath: statePath,
	}
}

// firstRune returns the first rune of s, or fallback if s is empty.
func firstRune(s string, fallback rune) rune {
	if s == "" {
		return fallback
	}
	for _, r := range s {
		return r
	}
	return fallback
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Save initial state and start periodic saves
	return tea.Batch(
		tickCmd(),
		saveStateCmd(),
	)
}

// tickMsg is a message sent periodically to update the timer display
type tickMsg struct {
	Time time.Time
}

// tickCmd returns a command that sends a tick message after 250ms
func tickCmd() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{Time: t}
	})
}

// saveStateMsg is a message to trigger state persistence
type saveStateMsg struct{}

// saveStateCmd returns a command that sends a save state message after 5 seconds
func saveStateCmd() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return saveStateMsg{}
	})
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 4
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			if m.timer.State() == timer.StateRunning {
				m.timer.Pause()
				// Save state on pause
				m.saveState()
				// Stop periodic saves while paused
				return m, nil
			}
			return m, nil
		case "r":
			if m.timer.State() == timer.StatePaused {
				m.timer.Resume()
				// Save state on resume
				m.saveState()
				// Resume periodic saves
				return m, tea.Batch(
					tickCmd(),
					saveStateCmd(),
				)
			}
			return m, nil
		case "s", "q":
			m.timer.Stop()
			// Save final state before exit
			m.saveState()
			m.quitting = true
			return m, tea.Quit
		case "ctrl+c":
			// Save state on interrupt
			m.saveState()
			m.quitting = true
			return m, tea.Quit
		}

	case tickMsg:
		if m.timer.State() == timer.StateRunning {
			if m.timer.IsCompleted() {
				// Timer completed - save final state
				m.saveState()
				return m, tea.Quit
			}
			return m, tickCmd()
		}
		return m, nil

	case saveStateMsg:
		// Periodic state save (every 5 seconds while running)
		if m.timer.State() == timer.StateRunning {
			m.saveState()
			// Continue periodic saves
			return m, saveStateCmd()
		}
		// Timer is paused or stopped, don't continue periodic saves
		return m, nil
	}

	return m, nil
}

// saveState saves the current timer state to disk
func (m Model) saveState() {
	if err := timer.SaveState(m.timer, m.sessionID, m.statePath); err != nil {
		// Log error but don't fail - state persistence is best-effort
		// In a real implementation, we might want to show this in the UI
	}
}

// View renders the TUI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	th := m.theme
	if th == nil {
		th = theme.GetTheme("default")
	}

	// Terminal too small: show warning
	if m.width < minWidth || m.height < minHeight {
		warningStyle := lipgloss.NewStyle().Foreground(th.Colors.Warning)
		msg := fmt.Sprintf("Terminal too small!\nMinimum: %dx%d\nCurrent: %dx%d", minWidth, minHeight, m.width, m.height)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, warningStyle.Render(msg))
	}

	// Session header: "{Preset} Session: {Label}" or "Session: {Label}"
	sessionHeader := m.sessionHeader(th)
	// Progress bar
	duration := m.timer.Duration()
	remaining := m.timer.Remaining()
	elapsed := duration - remaining
	var percent float64
	if duration > 0 {
		percent = float64(elapsed) / float64(duration)
		if percent > 1 {
			percent = 1
		}
	}
	progressBar := m.progress.ViewAs(percent)
	// Time remaining MM:SS
	timeDisplay := m.formatRemaining(remaining)
	primaryStyle := lipgloss.NewStyle().Foreground(th.Colors.Primary)
	timeLine := primaryStyle.Render(timeDisplay)
	// Status
	statusLine := m.statusLine(th)
	// Control legend or completion message
	mutedStyle := lipgloss.NewStyle().Foreground(th.Colors.TextMuted)
	successStyle := lipgloss.NewStyle().Foreground(th.Colors.Success)

	state := m.timer.State()
	var bottomLine string
	if state == timer.StateCompleted {
		bottomLine = successStyle.Render("Session saved!")
	} else {
		if state == timer.StatePaused {
			bottomLine = mutedStyle.Render("[r] resume  [s] stop  [q] quit  [Ctrl+C] emergency exit")
		} else {
			bottomLine = mutedStyle.Render("[p] pause  [s] stop  [q] quit  [Ctrl+C] emergency exit")
		}
	}

	titleLine := th.TitleStyle().Render("Pomodoro Timer")
	inner := lipgloss.JoinVertical(lipgloss.Left,
		titleLine,
		"",
		sessionHeader,
		"",
		progressBar,
		"",
		timeLine,
		"",
		statusLine,
		"",
		bottomLine,
	)

	windowStyle := th.BorderStyle().Padding(1, 2)
	content := windowStyle.Render(inner)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) sessionHeader(th *theme.Theme) string {
	label := m.timer.Label()
	preset := m.timer.Preset()
	var text string
	if preset != "" {
		text = fmt.Sprintf("%s Session: %s", prettifyPreset(preset), label)
	} else {
		text = fmt.Sprintf("Session: %s", label)
	}
	return th.TitleStyle().Render(text)
}

func prettifyPreset(preset string) string {
	if preset == "" {
		return preset
	}
	// Simple capitalization: first letter upper, rest unchanged
	runes := []rune(preset)
	if len(runes) == 0 {
		return preset
	}
	if runes[0] >= 'a' && runes[0] <= 'z' {
		runes[0] -= 'a' - 'A'
	}
	return string(runes)
}

func (m Model) formatRemaining(d time.Duration) string {
	total := int(d.Round(time.Second).Seconds())
	if total < 0 {
		total = 0
	}
	minutes := total / 60
	seconds := total % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (m Model) statusLine(th *theme.Theme) string {
	state := m.timer.State()
	var statusText string
	var statusKey string
	switch state {
	case timer.StateRunning:
		statusText = "RUNNING"
		statusKey = "running"
	case timer.StatePaused:
		statusText = "⏸ PAUSED"
		statusKey = "paused"
	case timer.StateCompleted:
		statusText = "✓ COMPLETED"
		statusKey = "completed"
	case timer.StateStopped:
		statusText = "STOPPED"
		statusKey = "stopped"
	default:
		statusText = "—"
		statusKey = ""
	}
	return th.StatusStyle(statusKey).Render("Status: " + statusText)
}


