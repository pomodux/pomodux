package timer

import (
	"fmt"
	"time"
)

// State represents the timer state
type State string

const (
	StateIdle      State = "idle"
	StateRunning   State = "running"
	StatePaused    State = "paused"
	StateCompleted State = "completed"
	StateStopped   State = "stopped"
)

// Timer represents a timer instance
type Timer struct {
	duration    time.Duration
	label       string
	preset      string
	startTime   time.Time
	pausedAt    time.Time
	totalPaused time.Duration
	pausedCount int
	state       State
}

// NewTimer creates a new timer with the given duration and label
func NewTimer(duration time.Duration, label string, preset string) (*Timer, error) {
	if duration <= 0 {
		return nil, fmt.Errorf("duration must be positive, got %v", duration)
	}

	if duration > 24*time.Hour {
		return nil, fmt.Errorf("duration exceeds maximum (24h), got %v", duration)
	}

	if len(label) > 200 {
		return nil, fmt.Errorf("label too long (max 200 chars), got %d", len(label))
	}

	return &Timer{
		duration: duration,
		label:    label,
		preset:   preset,
		state:    StateIdle,
	}, nil
}

// Start starts the timer
func (t *Timer) Start() error {
	if t.state != StateIdle {
		return fmt.Errorf("timer cannot be started from state %s", t.state)
	}

	t.startTime = time.Now()
	t.state = StateRunning
	return nil
}

// Pause pauses the timer
func (t *Timer) Pause() error {
	if t.state != StateRunning {
		return fmt.Errorf("timer cannot be paused from state %s", t.state)
	}

	t.pausedAt = time.Now()
	t.pausedCount++
	t.state = StatePaused
	return nil
}

// Resume resumes the timer
func (t *Timer) Resume() error {
	if t.state != StatePaused {
		return fmt.Errorf("timer cannot be resumed from state %s", t.state)
	}

	// Add the paused duration to total
	t.totalPaused += time.Since(t.pausedAt)
	t.pausedAt = time.Time{}
	t.state = StateRunning
	return nil
}

// Stop stops the timer early
func (t *Timer) Stop() error {
	if t.state != StateRunning && t.state != StatePaused {
		return fmt.Errorf("timer cannot be stopped from state %s", t.state)
	}

	t.state = StateStopped
	return nil
}

// Remaining returns the remaining time
func (t *Timer) Remaining() time.Duration {
	if t.state == StateIdle || t.state == StateCompleted || t.state == StateStopped {
		return 0
	}

	elapsed := time.Since(t.startTime) - t.totalPaused
	if t.state == StatePaused {
		elapsed -= time.Since(t.pausedAt)
	}

	remaining := t.duration - elapsed
	if remaining < 0 {
		return 0
	}

	return remaining
}

// IsCompleted checks if the timer has completed
func (t *Timer) IsCompleted() bool {
	return t.Remaining() == 0 && t.state == StateRunning
}

// Duration returns the configured duration
func (t *Timer) Duration() time.Duration {
	return t.duration
}

// Label returns the timer label
func (t *Timer) Label() string {
	return t.label
}

// Preset returns the preset name (empty if custom duration)
func (t *Timer) Preset() string {
	return t.preset
}

// State returns the current timer state
func (t *Timer) State() State {
	return t.state
}

// PausedCount returns the number of times the timer has been paused
func (t *Timer) PausedCount() int {
	return t.pausedCount
}

// TotalPausedDuration returns the total time spent paused
func (t *Timer) TotalPausedDuration() time.Duration {
	if t.state == StatePaused {
		return t.totalPaused + time.Since(t.pausedAt)
	}
	return t.totalPaused
}

// StartTime returns when the timer started
func (t *Timer) StartTime() time.Time {
	return t.startTime
}

// FromState creates a Timer instance from a TimerState
// Used for resuming a timer from saved state
func FromState(state *TimerState) (*Timer, error) {
	return ResumeFromState(state)
}

// ToState converts a Timer to TimerState for persistence
func (t *Timer) ToState(sessionID string) *TimerState {
	remaining := t.Remaining()
	durationStr := formatDuration(t.duration)
	remainingStr := formatDuration(remaining)
	pausedDurationStr := formatDuration(t.TotalPausedDuration())

	return &TimerState{
		Version:        "1.0",
		SessionID:      sessionID,
		PID:            0, // Will be set by SaveState
		StartedAt:      t.startTime,
		Duration:       durationStr,
		Preset:         t.preset,
		Label:          t.label,
		Remaining:      remainingStr,
		IsPaused:       t.state == StatePaused,
		PausedCount:    t.pausedCount,
		PausedDuration: pausedDurationStr,
		LastUpdated:    time.Now(),
	}
}
