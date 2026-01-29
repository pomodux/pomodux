package timer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// TimerState represents the persisted timer state for crash recovery
type TimerState struct {
	Version        string    `json:"version"`
	SessionID      string    `json:"session_id"`
	PID            int       `json:"pid"`
	StartedAt      time.Time `json:"started_at"`
	Duration       string    `json:"duration"`
	Preset         string    `json:"preset,omitempty"`
	Label          string    `json:"label"`
	Remaining      string    `json:"remaining"`
	IsPaused       bool      `json:"is_paused"`
	PausedCount    int       `json:"paused_count"`
	PausedDuration string    `json:"paused_duration"`
	LastUpdated    time.Time `json:"last_updated"`
}

// SaveState saves the timer state to a JSON file using atomic write
func SaveState(timer *Timer, sessionID string, path string) error {
	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	// Calculate remaining time
	remaining := timer.Remaining()

	// Format durations as strings
	durationStr := FormatDuration(timer.duration)
	remainingStr := FormatDuration(remaining)
	pausedDurationStr := FormatDuration(timer.TotalPausedDuration())

	// Create state struct
	state := TimerState{
		Version:        "1.0",
		SessionID:      sessionID,
		PID:            os.Getpid(),
		StartedAt:      timer.startTime,
		Duration:       durationStr,
		Preset:         timer.preset,
		Label:          timer.label,
		Remaining:      remainingStr,
		IsPaused:       timer.state == StatePaused,
		PausedCount:    timer.pausedCount,
		PausedDuration: pausedDurationStr,
		LastUpdated:    time.Now(),
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Atomic write: write to temp file, then rename
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to save state file: %w", err)
	}

	return nil
}

// LoadState loads timer state from a JSON file
func LoadState(path string) (*TimerState, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("state file does not exist: %s", path)
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	// Parse JSON
	var state TimerState
	if err := json.Unmarshal(data, &state); err != nil {
		// Try to backup corrupted file
		backupPath := path + ".backup"
		os.WriteFile(backupPath, data, 0600)
		return nil, fmt.Errorf("failed to parse state file (backed up to %s): %w", backupPath, err)
	}

	return &state, nil
}

// DeleteState removes the timer state file
func DeleteState(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist, that's fine
		return nil
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete state file: %w", err)
	}

	return nil
}

// IsProcessAlive checks if a process with the given PID is still running
// Uses Signal(0) which doesn't send a signal but checks if process exists
func IsProcessAlive(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Signal 0 doesn't actually send a signal, just checks if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// ResumeFromState reconstructs a Timer from saved state
// Calculates remaining time from wall-clock and handles paused state correctly
func ResumeFromState(state *TimerState) (*Timer, error) {
	// Parse duration
	duration, err := time.ParseDuration(state.Duration)
	if err != nil {
		return nil, fmt.Errorf("failed to parse duration: %w", err)
	}

	// Create timer with saved parameters
	timer, err := NewTimer(duration, state.Label, state.Preset)
	if err != nil {
		return nil, fmt.Errorf("failed to create timer: %w", err)
	}

	// Restore timer state
	timer.startTime = state.StartedAt
	timer.pausedCount = state.PausedCount

	// Parse paused duration
	if state.PausedDuration != "" {
		pausedDuration, err := time.ParseDuration(state.PausedDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to parse paused duration: %w", err)
		}
		timer.totalPaused = pausedDuration
	}

	// Handle paused state
	if state.IsPaused {
		// When paused, we need to set pausedAt such that Remaining() calculates correctly
		// Parse saved remaining time
		savedRemaining, err := time.ParseDuration(state.Remaining)
		if err != nil {
			return nil, fmt.Errorf("failed to parse remaining duration: %w", err)
		}

		// Calculate elapsed time at save: elapsed = duration - remaining
		elapsedAtSave := duration - savedRemaining

		// Calculate what pausedAt should be:
		// At save time: elapsed = (LastUpdated - startTime) - totalPaused - (LastUpdated - pausedAt)
		// elapsed = pausedAt - startTime - totalPaused
		// pausedAt = startTime + totalPaused + elapsed
		timer.pausedAt = state.StartedAt.Add(timer.totalPaused + elapsedAtSave)

		timer.state = StatePaused
	} else {
		// Timer was running, set it to running state
		// The startTime is already set, so Remaining() will calculate correctly
		timer.state = StateRunning
	}

	return timer, nil
}

// FormatDuration formats a time.Duration as a string (e.g., "25m", "1h30m").
// Used when building session history and persisted state.
func FormatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string

	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 && hours == 0 {
		// Only show seconds if less than an hour
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	if len(parts) == 0 {
		return "0s"
	}

	return strings.Join(parts, "")
}
