package timer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pomodux/pomodux/internal/logger"
)

// SingleTimerCoordinator ensures only one timer runs at a time
type SingleTimerCoordinator struct {
	mu          sync.RWMutex
	activeTimer *EventDrivenTimer
	lockFile    string
	processID   int
	startTime   time.Time
}

// LockInfo represents the information stored in the lock file
type LockInfo struct {
	ProcessID   int       `json:"process_id"`
	StartTime   time.Time `json:"start_time"`
	SessionName string    `json:"session_name"`
	Duration    int       `json:"duration_seconds"`
	Remaining   int       `json:"remaining_seconds"`
}

// NewSingleTimerCoordinator creates a new single timer coordinator
func NewSingleTimerCoordinator() (*SingleTimerCoordinator, error) {
	lockDir, err := getLockDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get lock directory: %w", err)
	}

	// Ensure lock directory exists
	if err := os.MkdirAll(lockDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create lock directory: %w", err)
	}

	lockFile := filepath.Join(lockDir, "pomodux.lock")

	coordinator := &SingleTimerCoordinator{
		lockFile:  lockFile,
		processID: os.Getpid(),
	}

	logger.Debug("Single timer coordinator initialized", map[string]interface{}{
		"lock_file":  lockFile,
		"process_id": coordinator.processID,
	})

	return coordinator, nil
}

// AcquireTimerLock attempts to acquire exclusive timer lock
func (stc *SingleTimerCoordinator) AcquireTimerLock(timer *EventDrivenTimer) error {
	stc.mu.Lock()
	defer stc.mu.Unlock()

	// Check if timer already running in this process
	if stc.activeTimer != nil && stc.activeTimer.GetStatus() == StatusRunning {
		return fmt.Errorf("timer already running (started at %v)",
			stc.startTime.Format("15:04:05"))
	}

	// Check system-wide lock file
	if err := stc.checkSystemLock(); err != nil {
		return err
	}

	// Acquire lock
	if err := stc.createLockFile(timer); err != nil {
		return err
	}

	stc.activeTimer = timer
	stc.startTime = time.Now()

	logger.Info("Timer lock acquired", map[string]interface{}{
		"process_id":   stc.processID,
		"start_time":   stc.startTime,
		"session_name": timer.GetSessionName(),
	})

	return nil
}

// ReleaseLock releases the timer lock
func (stc *SingleTimerCoordinator) ReleaseLock() {
	stc.mu.Lock()
	defer stc.mu.Unlock()

	stc.activeTimer = nil
	stc.removeLockFile()

	logger.Info("Timer lock released", map[string]interface{}{
		"process_id": stc.processID,
	})
}

// IsTimerRunning checks if a timer is currently running
func (stc *SingleTimerCoordinator) IsTimerRunning() bool {
	stc.mu.RLock()
	defer stc.mu.RUnlock()

	return stc.activeTimer != nil &&
		stc.activeTimer.GetStatus() == StatusRunning
}

// GetActiveTimer returns the currently active timer, if any
func (stc *SingleTimerCoordinator) GetActiveTimer() *EventDrivenTimer {
	stc.mu.RLock()
	defer stc.mu.RUnlock()

	return stc.activeTimer
}

// checkSystemLock checks for existing system-wide lock
func (stc *SingleTimerCoordinator) checkSystemLock() error {
	if _, err := os.Stat(stc.lockFile); os.IsNotExist(err) {
		return nil // No lock file exists
	}

	lockInfo, err := stc.readLockFile()
	if err != nil {
		return err
	}
	if lockInfo == nil {
		return nil // Lock file was cleaned up
	}

	return stc.validateExistingLock(*lockInfo)
}

// readLockFile reads and validates the lock file, cleaning up corrupted files
func (stc *SingleTimerCoordinator) readLockFile() (*LockInfo, error) {
	data, err := os.ReadFile(stc.lockFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	var lockInfo LockInfo
	if err := json.Unmarshal(data, &lockInfo); err != nil {
		stc.handleCorruptedLockFile(err)
		return nil, nil
	}

	return &lockInfo, nil
}

// handleCorruptedLockFile removes corrupted lock files
func (stc *SingleTimerCoordinator) handleCorruptedLockFile(err error) {
	logger.Warn("Corrupted lock file detected, removing", map[string]interface{}{
		"lock_file": stc.lockFile,
		"error":     err.Error(),
	})
	if removeErr := os.Remove(stc.lockFile); removeErr != nil {
		logger.Warn("Failed to remove corrupted lock file", map[string]interface{}{
			"lock_file": stc.lockFile,
			"error":     removeErr.Error(),
		})
	}
}

// validateExistingLock checks if an existing lock is still valid
func (stc *SingleTimerCoordinator) validateExistingLock(lockInfo LockInfo) error {
	if !isProcessRunning(lockInfo.ProcessID) {
		stc.handleStaleLockFile(lockInfo.ProcessID)
		return nil
	}

	elapsed := time.Since(lockInfo.StartTime)
	remaining := time.Duration(lockInfo.Duration)*time.Second - elapsed

	if remaining <= 0 {
		stc.handleExpiredLockFile(elapsed, time.Duration(lockInfo.Duration)*time.Second)
		return nil
	}

	return stc.createActiveTimerError(lockInfo, remaining)
}

// handleStaleLockFile removes stale lock files
func (stc *SingleTimerCoordinator) handleStaleLockFile(stalePID int) {
	logger.Warn("Stale lock file detected, removing", map[string]interface{}{
		"lock_file":   stc.lockFile,
		"stale_pid":   stalePID,
		"current_pid": stc.processID,
	})
	if removeErr := os.Remove(stc.lockFile); removeErr != nil {
		logger.Warn("Failed to remove stale lock file", map[string]interface{}{
			"lock_file": stc.lockFile,
			"error":     removeErr.Error(),
		})
	}
}

// handleExpiredLockFile removes expired lock files
func (stc *SingleTimerCoordinator) handleExpiredLockFile(elapsed, duration time.Duration) {
	logger.Warn("Expired lock file detected, removing", map[string]interface{}{
		"lock_file": stc.lockFile,
		"elapsed":   elapsed,
		"duration":  duration,
	})
	if removeErr := os.Remove(stc.lockFile); removeErr != nil {
		logger.Warn("Failed to remove expired lock file", map[string]interface{}{
			"lock_file": stc.lockFile,
			"error":     removeErr.Error(),
		})
	}
}

// createActiveTimerError creates an error for active timer conflicts
func (stc *SingleTimerCoordinator) createActiveTimerError(lockInfo LockInfo, remaining time.Duration) error {
	return fmt.Errorf("timer already running in process %d\n"+
		"Session: %s\n"+
		"Started: %v\n"+
		"Remaining: %v\n"+
		"Use 'pomodux status' to check progress or 'pomodux stop' to stop",
		lockInfo.ProcessID,
		lockInfo.SessionName,
		lockInfo.StartTime.Format("15:04:05"),
		remaining.Round(time.Second))
}

// createLockFile creates the lock file with timer information
func (stc *SingleTimerCoordinator) createLockFile(timer *EventDrivenTimer) error {
	lockInfo := LockInfo{
		ProcessID:   stc.processID,
		StartTime:   time.Now(),
		SessionName: timer.GetSessionName(),
		Duration:    int(timer.GetDuration().Seconds()),
		Remaining:   int((timer.GetDuration() - timer.GetElapsed()).Seconds()),
	}

	data, err := json.Marshal(lockInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal lock info: %w", err)
	}

	// Create lock file atomically
	if err := os.WriteFile(stc.lockFile, data, 0600); err != nil {
		return fmt.Errorf("failed to create lock file: %w", err)
	}

	logger.Debug("Lock file created", map[string]interface{}{
		"lock_file":    stc.lockFile,
		"process_id":   lockInfo.ProcessID,
		"session_name": lockInfo.SessionName,
		"duration":     lockInfo.Duration,
	})

	return nil
}

// removeLockFile removes the lock file
func (stc *SingleTimerCoordinator) removeLockFile() {
	if err := os.Remove(stc.lockFile); err != nil && !os.IsNotExist(err) {
		logger.Warn("Failed to remove lock file", map[string]interface{}{
			"lock_file": stc.lockFile,
			"error":     err.Error(),
		})
	} else {
		logger.Debug("Lock file removed", map[string]interface{}{
			"lock_file": stc.lockFile,
		})
	}
}

// isProcessRunning checks if a process with the given PID is running
func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Unix systems, sending signal 0 checks if process exists
	err = process.Signal(os.Signal(nil))
	return err == nil
}
