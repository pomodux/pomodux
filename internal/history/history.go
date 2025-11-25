package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Session represents a completed timer session
type Session struct {
	ID             string    `json:"id"`
	StartedAt      time.Time `json:"started_at"`
	EndedAt        time.Time `json:"ended_at"`
	Duration       string    `json:"duration"` // e.g., "25m"
	Preset         string    `json:"preset,omitempty"`
	Label          string    `json:"label"`
	EndStatus      string    `json:"end_status"` // completed, stopped, cancelled, interrupted
	PausedCount    int       `json:"paused_count"`
	PausedDuration string    `json:"paused_duration"` // e.g., "3m"
}

// History represents the session history
type History struct {
	Version  string    `json:"version"`
	Sessions []Session `json:"sessions"`
}

// Load loads history from the given path
func Load(path string) (*History, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &History{
			Version:  "1.0",
			Sessions: []Session{},
		}, nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	// Parse JSON
	var history History
	if err := json.Unmarshal(data, &history); err != nil {
		// Try to backup corrupted file
		backupPath := path + ".backup"
		os.WriteFile(backupPath, data, 0600)
		return nil, fmt.Errorf("failed to parse history file (backed up to %s): %w", backupPath, err)
	}

	return &history, nil
}

// Save saves history to the given path using atomic write
func Save(history *History, path string) error {
	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create history directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	// Atomic write: write to temp file, then rename
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to save history file: %w", err)
	}

	return nil
}

// AddSession adds a session to the history
func (h *History) AddSession(session Session) {
	h.Sessions = append(h.Sessions, session)
}

// GetRecent returns the most recent N sessions
func (h *History) GetRecent(limit int) []Session {
	if limit <= 0 || limit > len(h.Sessions) {
		limit = len(h.Sessions)
	}

	// Sessions are stored newest first (most recent appended)
	start := len(h.Sessions) - limit
	if start < 0 {
		start = 0
	}

	result := make([]Session, limit)
	copy(result, h.Sessions[start:])
	
	// Reverse to show newest first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}


