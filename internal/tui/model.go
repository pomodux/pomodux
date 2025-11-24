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
}

// NewModel creates a new TUI model
func NewModel(t *timer.Timer) Model {
	prog := progress.New(progress.WithDefaultGradient())
	
	return Model{
		timer:    t,
		progress: prog,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
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
			}
			return m, nil
		case "r":
			if m.timer.State() == timer.StatePaused {
				m.timer.Resume()
				return m, tickCmd()
			}
			return m, nil
		case "s", "q":
			m.timer.Stop()
			m.quitting = true
			return m, tea.Quit
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

	case tickMsg:
		if m.timer.State() == timer.StateRunning {
			if m.timer.IsCompleted() {
				// Timer completed
				return m, tea.Quit
			}
			return m, tickCmd()
		}
		return m, nil
	}

	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// This is a placeholder - full rendering will be implemented later
	return "Timer TUI (placeholder)\n"
}

