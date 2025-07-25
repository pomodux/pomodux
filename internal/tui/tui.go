package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pomodux/pomodux/internal/timer"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Align(lipgloss.Center)
	boxStyle      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 4).Align(lipgloss.Center)
	barStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	controlsStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	statusStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81")).Align(lipgloss.Center)
)

// Model represents the state of the TUI application.
type Model struct {
	Timer        *timer.Timer
	Paused       bool
	Quitting     bool
	ForceStopped bool
	Width        int
	Height       int
}

func NewModel(t *timer.Timer) Model {
	return Model{
		Timer:        t,
		Paused:       t.GetStatus() == timer.StatusPaused,
		Quitting:     false,
		ForceStopped: false,
	}
}

type tickMsg struct{}

func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case tickMsg:
		status := m.Timer.GetStatus()
		if status == timer.StatusCompleted || status == timer.StatusIdle {
			m.Quitting = true
			return m, tea.Quit
		}
		if status == timer.StatusPaused {
			m.Paused = true
		} else {
			m.Paused = false
		}
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg{} })
	case tea.KeyMsg:
		key := strings.ToLower(msg.String())
		switch key {
		case "p":
			if m.Timer.GetStatus() == timer.StatusRunning {
				_ = m.Timer.Pause()
			}
		case "r":
			if m.Timer.GetStatus() == timer.StatusPaused {
				_ = m.Timer.Resume()
			}
		case "q", "s":
			m.ForceStopped = true
			_ = m.Timer.Stop()
			m.Quitting = true
			return m, tea.Quit
		case "ctrl+c":
			m.ForceStopped = true
			_ = m.Timer.Stop()
			m.Quitting = true
			return m, tea.Quit
		}
	}
	if m.Quitting {
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) renderProgressBar(progress float64, barWidth int) string {
	filled := int(progress * float64(barWidth))
	return barStyle.Render(strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled))
}

func (m Model) renderControls(barWidth int) string {
	controlsFull := "[P]ause  [R]esume  [S]top  [Q]uit"
	controlsShort := "[P]ause [R]esume [S]top [Q]uit"
	controls := controlsFull
	if barWidth+5 < len(controlsFull) {
		controls = controlsShort
	}
	return controlsStyle.Render(controls)
}

func (m Model) View() string {
	if m.Quitting {
		return "Exiting Pomodux...\n"
	}
	sessionType := m.Timer.GetSessionType()
	title := titleStyle.Render(strings.ToUpper(string(sessionType)) + " SESSION")
	status := "[RUNNING]"
	if m.Paused {
		status = "[PAUSED]"
	}
	statusLine := statusStyle.Render(status)
	dur := m.Timer.GetDuration()
	rem := dur - m.Timer.GetElapsed()
	timeStr := fmt.Sprintf("%02dm %02ds remaining", int(rem.Minutes()), int(rem.Seconds())%60)
	progress := float64(dur-rem) / float64(dur)
	barWidth := 40
	percent := int(progress * 100)
	percentStr := fmt.Sprintf("%3d%%", percent)
	if m.Width > 0 && m.Width-16 < barWidth+5 {
		barWidth = m.Width - 16 - 5
		if barWidth < 10 {
			barWidth = 10
		}
	}
	bar := m.renderProgressBar(progress, barWidth)
	controls := m.renderControls(barWidth)

	box := lipgloss.JoinVertical(lipgloss.Left,
		"",
		title,
		statusLine,
		"",
		timeStr,
		"",
		bar+" "+percentStr,
		"",
		controls,
		"",
	)
	box = boxStyle.Width(barWidth + 8 + 5).Render(box)

	padTop := 0
	padLeft := 0
	if m.Height > 0 {
		padTop = (m.Height - lipgloss.Height(box)) / 2
		if padTop < 0 {
			padTop = 0
		}
	}
	if m.Width > 0 {
		padLeft = (m.Width - lipgloss.Width(box)) / 2
		if padLeft < 0 {
			padLeft = 0
		}
	}
	return strings.Repeat("\n", padTop) + lipgloss.NewStyle().MarginLeft(padLeft).Render(box)
}

// RunTUI launches the Bubbletea TUI app with the given timer instance.
func RunTUI(t *timer.Timer) error {
	model := NewModel(t)
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	if m, ok := finalModel.(Model); ok {
		if !m.ForceStopped && m.Timer.GetStatus() == timer.StatusCompleted {
			fmt.Println("Pomodux session complete!")
		}
	}
	return nil
}
