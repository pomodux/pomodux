package timer

import "time"

// TimerEngine defines the core timer interface
// for Pomodux timer operations.
type TimerEngine interface {
	Start(duration time.Duration) error
	StartWithSessionName(duration time.Duration, sessionName string) error
	Stop() error
	Pause() error
	Resume() error
	GetStatus() TimerStatus
	GetProgress() float64
	GetSessionName() string
}

// TimerStatus represents the state of the timer.
type TimerStatus string

const (
	StatusIdle      TimerStatus = "idle"
	StatusRunning   TimerStatus = "running"
	StatusPaused    TimerStatus = "paused"
	StatusCompleted TimerStatus = "completed"
)

// Session type constants removed - sessions now use generic string names
