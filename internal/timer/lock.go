package timer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/pomodux/pomodux/internal/logger"
)

// TimerLockManager manages file-based locking to ensure single timer instance
type TimerLockManager struct {
	lockFile     string
	lockFd       *os.File
	locked       bool
	timerPID     int
	sessionName  string
	mu           sync.Mutex
}

// LockFileState represents the state stored in the lock file
type LockFileState struct {
	PID         int       `json:"pid"`
	SessionName string    `json:"session_name"`
	StartTime   time.Time `json:"start_time"`
	Duration    int       `json:"duration_seconds"`
	Locked      time.Time `json:"locked_at"`
}

// TimerConflictError represents an error when timer is already running
type TimerConflictError struct {
	PID         int
	SessionName string
	StartTime   time.Time
	Duration    time.Duration
}

func (e *TimerConflictError) Error() string {
	elapsed := time.Since(e.StartTime)
	remaining := e.Duration - elapsed
	
	if remaining <= 0 {
		return fmt.Sprintf("Timer process %d may have completed. Try: pomodux status", e.PID)
	}
	
	return fmt.Sprintf(
		"Timer already running in process %d (%s remaining in '%s' session)\n" +
		"Try: pomodux status | pomodux stop | pomodux pause",
		e.PID,
		formatDuration(remaining),
		e.SessionName,
	)
}

var ErrTimerAlreadyRunning = fmt.Errorf("timer already running")

// NewTimerLockManager creates a new timer lock manager
func NewTimerLockManager() (*TimerLockManager, error) {
	lockDir, err := getLockDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get lock directory: %w", err)
	}

	// Ensure lock directory exists
	if err := os.MkdirAll(lockDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create lock directory: %w", err)
	}

	lockFile := filepath.Join(lockDir, "timer.lock")
	
	logger.Debug("Lock manager initialized", map[string]interface{}{
		"lock_file": lockFile,
		"process_id": os.Getpid(),
	})

	return &TimerLockManager{
		lockFile: lockFile,
		timerPID: os.Getpid(),
	}, nil
}

// AcquireLock attempts to acquire exclusive timer lock
//nolint:funlen // Complex lock acquisition logic is better kept together
func (lm *TimerLockManager) AcquireLock(sessionName string, duration time.Duration) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	logger.Debug("Attempting to acquire timer lock", map[string]interface{}{
		"session_name": sessionName,
		"duration": duration.String(),
		"lock_file": lm.lockFile,
	})

	// Check if lock file exists and is valid
	if _, err := os.Stat(lm.lockFile); err == nil {
		// Lock file exists, check if process is still running
		if err := lm.recoverOrphanedLock(); err != nil {
			if err == ErrTimerAlreadyRunning {
				// Read current state and return conflict error
				state, readErr := lm.ReadLockState()
				if readErr != nil {
					logger.Error("Failed to read lock state during conflict", readErr, map[string]interface{}{
						"lock_file": lm.lockFile,
					})
					return fmt.Errorf("timer already running (unable to read details)")
				}
				
				logger.Warn("Timer already running in another process", map[string]interface{}{
					"running_pid": state.PID,
					"running_session": state.SessionName,
					"running_start_time": state.StartTime,
					"attempted_session": sessionName,
					"attempted_duration": duration.String(),
				})
				
				return &TimerConflictError{
					PID:         state.PID,
					SessionName: state.SessionName,
					StartTime:   state.StartTime,
					Duration:    time.Duration(state.Duration) * time.Second,
				}
			}
			return fmt.Errorf("failed to recover orphaned lock: %w", err)
		}
	}

	// Create lock file with timer state
	state := LockFileState{
		PID:         lm.timerPID,
		SessionName: sessionName,
		StartTime:   time.Now(),
		Duration:    int(duration.Seconds()),
		Locked:      time.Now(),
	}

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal lock state: %w", err)
	}

	// Create lock file atomically
	lockFd, err := os.OpenFile(lm.lockFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		// If file exists, another process beat us to it
		if os.IsExist(err) {
			return ErrTimerAlreadyRunning
		}
		return fmt.Errorf("failed to create lock file: %w", err)
	}

	if _, err := lockFd.Write(data); err != nil {
		lockFd.Close()
		os.Remove(lm.lockFile)
		return fmt.Errorf("failed to write lock state: %w", err)
	}

	lm.lockFd = lockFd
	lm.locked = true
	lm.sessionName = sessionName

	logger.Info("Timer lock acquired successfully", map[string]interface{}{
		"session_name": sessionName,
		"duration": duration.String(),
		"lock_file": lm.lockFile,
		"process_id": lm.timerPID,
	})

	return nil
}

// ReleaseLock releases the timer lock
func (lm *TimerLockManager) ReleaseLock() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if !lm.locked {
		return nil // Not locked, nothing to do
	}

	logger.Debug("Releasing timer lock", map[string]interface{}{
		"session_name": lm.sessionName,
		"lock_file": lm.lockFile,
		"process_id": lm.timerPID,
	})

	// Close file descriptor
	if lm.lockFd != nil {
		if err := lm.lockFd.Close(); err != nil {
			logger.Warn("Failed to close lock file descriptor", map[string]interface{}{
				"error": err.Error(),
				"lock_file": lm.lockFile,
			})
		}
		lm.lockFd = nil
	}

	// Remove lock file
	if err := os.Remove(lm.lockFile); err != nil && !os.IsNotExist(err) {
		logger.Error("Failed to remove lock file", err, map[string]interface{}{
			"lock_file": lm.lockFile,
		})
		return fmt.Errorf("failed to remove lock file: %w", err)
	}

	lm.locked = false

	logger.Info("Timer lock released successfully", map[string]interface{}{
		"session_name": lm.sessionName,
		"lock_file": lm.lockFile,
		"process_id": lm.timerPID,
	})

	return nil
}

// ReadLockState reads the current lock file state
func (lm *TimerLockManager) ReadLockState() (*LockFileState, error) {
	data, err := os.ReadFile(lm.lockFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	var state LockFileState
	if err := json.Unmarshal(data, &state); err != nil {
		logger.Warn("Lock file corrupted, will delete", map[string]interface{}{
			"lock_file": lm.lockFile,
			"error": err.Error(),
		})
		// Delete corrupted lock file
		if removeErr := os.Remove(lm.lockFile); removeErr != nil {
			logger.Error("Failed to delete corrupted lock file", removeErr, map[string]interface{}{
				"lock_file": lm.lockFile,
			})
		}
		return nil, fmt.Errorf("lock file corrupted and deleted: %w", err)
	}

	return &state, nil
}

// validateProcess checks if a process is still running and is a valid pomodux process
func (lm *TimerLockManager) validateProcess(pid int) bool {
	logger.Debug("Validating process for lock ownership", map[string]interface{}{
		"pid": pid,
		"validator_pid": os.Getpid(),
	})
	
	// Check if process exists
	process, err := os.FindProcess(pid)
	if err != nil {
		logger.Debug("Process not found during validation", map[string]interface{}{
			"pid": pid,
			"error": err.Error(),
		})
		return false
	}
	
	// On Unix systems, check if process is still alive
	if err := process.Signal(syscall.Signal(0)); err != nil {
		logger.Debug("Process not responding to signal", map[string]interface{}{
			"pid": pid,
			"error": err.Error(),
		})
		return false
	}
	
	logger.Debug("Process validation complete", map[string]interface{}{
		"pid": pid,
		"valid": true,
	})
	return true
}

// recoverOrphanedLock attempts to recover from orphaned lock
func (lm *TimerLockManager) recoverOrphanedLock() error {
	logger.Info("Attempting to recover orphaned lock", map[string]interface{}{
		"lock_file": lm.lockFile,
	})
	
	state, err := lm.ReadLockState()
	if err != nil {
		logger.Error("Failed to read lock state during recovery", err, map[string]interface{}{
			"lock_file": lm.lockFile,
		})
		return err
	}
	
	if !lm.validateProcess(state.PID) {
		logger.Warn("Recovering orphaned timer lock", map[string]interface{}{
			"orphaned_pid": state.PID,
			"session_name": state.SessionName,
			"start_time": state.StartTime,
			"lock_age": time.Since(state.Locked),
		})
		
		if err := lm.forceReleaseLock(); err != nil {
			logger.Error("Failed to force release orphaned lock", err, map[string]interface{}{
				"orphaned_pid": state.PID,
			})
			return err
		}
		
		logger.Info("Orphaned lock recovered successfully", map[string]interface{}{
			"orphaned_pid": state.PID,
			"session_name": state.SessionName,
		})
		return nil
	}
	
	logger.Debug("Lock owner process is still valid", map[string]interface{}{
		"owner_pid": state.PID,
		"session_name": state.SessionName,
	})
	return ErrTimerAlreadyRunning
}

// forceReleaseLock forcefully removes the lock file
func (lm *TimerLockManager) forceReleaseLock() error {
	return os.Remove(lm.lockFile)
}

// getLockDir returns the XDG-compliant lock directory
func getLockDir() (string, error) {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		// Fallback to state directory
		stateDir, err := getStateDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(stateDir, "runtime"), nil
	}
	return filepath.Join(runtimeDir, "pomodux"), nil
}