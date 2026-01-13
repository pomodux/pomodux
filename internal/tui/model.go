package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pomodux/pomodux/internal/timer"
)

// Model represents the TUI model
type Model struct {
	timer      *timer.Timer
	progress   progress.Model
	width      int
	height     int
	quitting   bool
	sessionID  string
	statePath  string
}

// NewModel creates a new TUI model
func NewModel(t *timer.Timer, sessionID string, statePath string) Model {
	prog := progress.New(progress.WithDefaultGradient())
	
	return Model{
		timer:     t,
		progress:  prog,
		sessionID: sessionID,
		statePath: statePath,
	}
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

	// This is a placeholder - full rendering will be implemented later
	return "Timer TUI (placeholder)\n"
}


