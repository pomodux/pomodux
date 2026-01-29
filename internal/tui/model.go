package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pomodux/pomodux/internal/logger"
	"github.com/pomodux/pomodux/internal/theme"
	"github.com/pomodux/pomodux/internal/timer"
)

const (
	minWidth  = 80
	minHeight = 24
)

// Model represents the TUI model
type Model struct {
	timer                        *timer.Timer
	progress                     progress.Model
	theme                        *theme.Theme
	width                        int
	height                       int
	quitting                     bool
	sessionID                    string
	statePath                    string
	showConfirmation             bool
	wasRunningBeforeConfirmation bool
	showCompletion               bool
	completionCountdown          int
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
	logger.WithFields(map[string]interface{}{
		"component": "tui",
		"event":     "tui_initialized",
		"session_id": m.sessionID,
		"timer_state": string(m.timer.State()),
	}).Info("TUI initialized")
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

// completionTickMsg is a message sent every second during completion countdown
type completionTickMsg struct{}

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
		logger.WithFields(map[string]interface{}{
			"component": "tui",
			"event":     "window_resize",
			"width":     msg.Width,
			"height":    msg.Height,
			"progress_width": m.progress.Width,
		}).Debug("Window resized")
		return m, nil

	case tea.KeyMsg:
		key := msg.String()
		logger.WithFields(map[string]interface{}{
			"component":        "tui",
			"event":            "keypress",
			"key":              key,
			"timer_state":      string(m.timer.State()),
			"show_confirmation": m.showConfirmation,
			"show_completion":  m.showCompletion,
		}).Debug("Key pressed")

		// Handle confirmation state first
		if m.showConfirmation {
			switch key {
			case "y", "Y":
				// Confirm stop
				logger.WithFields(map[string]interface{}{
					"component": "tui",
					"event":     "confirmation_confirmed",
					"session_id": m.sessionID,
				}).Info("Stop confirmed by user")
				m.timer.Stop()
				m.saveState()
				m.quitting = true
				return m, tea.Quit
			case "n", "N", "esc":
				// Cancel confirmation
				logger.WithFields(map[string]interface{}{
					"component":        "tui",
					"event":            "confirmation_cancelled",
					"was_running":      m.wasRunningBeforeConfirmation,
				}).Info("Stop cancelled by user")
				m.showConfirmation = false
				if m.wasRunningBeforeConfirmation {
					m.timer.Resume()
					m.saveState()
					return m, tea.Batch(
						tickCmd(),
						saveStateCmd(),
					)
				}
				return m, nil
			default:
				// Other keys ignored in confirmation state
				logger.WithFields(map[string]interface{}{
					"component": "tui",
					"event":     "key_ignored",
					"key":       key,
					"reason":    "in_confirmation_state",
				}).Debug("Key ignored in confirmation state")
				return m, nil
			}
		}

		// Normal key handling (not in confirmation state)
		switch key {
		case "p":
			if m.timer.State() == timer.StateRunning {
				logger.WithFields(map[string]interface{}{
					"component": "tui",
					"event":     "pause",
					"session_id": m.sessionID,
				}).Info("Timer paused")
				m.timer.Pause()
				// Save state on pause
				m.saveState()
				// Stop periodic saves while paused
				return m, nil
			}
			logger.WithFields(map[string]interface{}{
				"component": "tui",
				"event":     "key_ignored",
				"key":       key,
				"reason":    "timer_not_running",
				"timer_state": string(m.timer.State()),
			}).Debug("Pause key ignored - timer not running")
			return m, nil
		case "r":
			if m.timer.State() == timer.StatePaused {
				logger.WithFields(map[string]interface{}{
					"component": "tui",
					"event":     "resume",
					"session_id": m.sessionID,
				}).Info("Timer resumed")
				m.timer.Resume()
				// Save state on resume
				m.saveState()
				// Resume periodic saves
				return m, tea.Batch(
					tickCmd(),
					saveStateCmd(),
				)
			}
			logger.WithFields(map[string]interface{}{
				"component": "tui",
				"event":     "key_ignored",
				"key":       key,
				"reason":    "timer_not_paused",
				"timer_state": string(m.timer.State()),
			}).Debug("Resume key ignored - timer not paused")
			return m, nil
		case "s", "q":
			// Enter confirmation state
			logger.WithFields(map[string]interface{}{
				"component":        "tui",
				"event":            "stop_requested",
				"session_id":       m.sessionID,
				"timer_state":      string(m.timer.State()),
			}).Info("Stop requested - entering confirmation")
			m.showConfirmation = true
			m.wasRunningBeforeConfirmation = (m.timer.State() == timer.StateRunning)
			if m.timer.State() == timer.StateRunning {
				m.timer.Pause()
				m.saveState()
			}
			return m, nil
		case "ctrl+c":
			// Emergency exit - bypass confirmation
			logger.WithFields(map[string]interface{}{
				"component": "tui",
				"event":     "emergency_exit",
				"session_id": m.sessionID,
			}).Warn("Emergency exit (Ctrl+C) - bypassing confirmation")
			m.timer.Stop()
			// Save state on interrupt
			m.saveState()
			m.quitting = true
			return m, tea.Quit
		default:
			logger.WithFields(map[string]interface{}{
				"component": "tui",
				"event":     "key_ignored",
				"key":       key,
				"reason":    "unknown_key",
			}).Debug("Unknown key ignored")
			return m, nil
		}

	case tickMsg:
		if m.timer.State() == timer.StateRunning {
			if m.timer.IsCompleted() {
				// Enter completion countdown state
				if !m.showCompletion {
					logger.WithFields(map[string]interface{}{
						"component": "tui",
						"event":     "timer_completed",
						"session_id": m.sessionID,
					}).Info("Timer completed - starting countdown")
					m.showCompletion = true
					m.completionCountdown = 3
					m.saveState()
					return m, tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
						return completionTickMsg{}
					})
				}
			}
			return m, tickCmd()
		}
		return m, nil

	case completionTickMsg:
		if m.showCompletion {
			m.completionCountdown--
			logger.WithFields(map[string]interface{}{
				"component":         "tui",
				"event":             "completion_countdown",
				"session_id":        m.sessionID,
				"countdown":         m.completionCountdown,
			}).Debug("Completion countdown tick")
			if m.completionCountdown <= 0 {
				logger.WithFields(map[string]interface{}{
					"component": "tui",
					"event":     "completion_exit",
					"session_id": m.sessionID,
				}).Info("Completion countdown finished - exiting")
				return m, tea.Quit
			}
			return m, tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
				return completionTickMsg{}
			})
		}
		return m, nil

	case saveStateMsg:
		// Periodic state save (every 5 seconds while running)
		if m.timer.State() == timer.StateRunning {
			logger.WithFields(map[string]interface{}{
				"component": "tui",
				"event":     "periodic_state_save",
				"session_id": m.sessionID,
			}).Debug("Periodic state save")
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
		logger.WithError(err).WithFields(map[string]interface{}{
			"component": "tui",
			"event":     "state_save_error",
			"session_id": m.sessionID,
			"state_path": m.statePath,
		}).Error("Failed to save timer state")
	} else {
		logger.WithFields(map[string]interface{}{
			"component": "tui",
			"event":     "state_saved",
			"session_id": m.sessionID,
			"timer_state": string(m.timer.State()),
		}).Debug("Timer state saved")
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
	// Control legend, confirmation prompt, or completion message
	mutedStyle := lipgloss.NewStyle().Foreground(th.Colors.TextMuted)
	successStyle := lipgloss.NewStyle().Foreground(th.Colors.Success)
	warningStyle := lipgloss.NewStyle().Foreground(th.Colors.Warning)

	state := m.timer.State()
	var bottomLine string

	// Priority: completion > confirmation > normal state
	if m.showCompletion {
		bottomLine = successStyle.Render(fmt.Sprintf("Session saved! Closing in %d.", m.completionCountdown))
	} else if m.showConfirmation {
		bottomLine = warningStyle.Render("Stop timer and exit? [y]es / [n]o")
	} else if state == timer.StateCompleted {
		bottomLine = successStyle.Render("Session saved!")
	} else {
		if state == timer.StatePaused {
			bottomLine = mutedStyle.Render("[r]esume  [s]top")
		} else {
			bottomLine = mutedStyle.Render("[p]ause  [s]top")
		}
	}

	inner := lipgloss.JoinVertical(lipgloss.Left,
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


