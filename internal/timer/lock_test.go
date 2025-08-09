package timer

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTimerLockManager_AcquireLock(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()

	// Create lock manager with custom lock file
	lm := &TimerLockManager{
		lockFile: filepath.Join(tmpDir, "test_timer.lock"),
		timerPID: os.Getpid(),
	}

	// Test successful lock acquisition
	err := lm.AcquireLock("test-session", 25*time.Minute)
	if err != nil {
		t.Fatalf("Expected successful lock acquisition, got error: %v", err)
	}

	// Verify lock is marked as acquired
	if !lm.locked {
		t.Error("Expected lock to be marked as acquired")
	}

	// Verify lock file exists
	if _, err := os.Stat(lm.lockFile); os.IsNotExist(err) {
		t.Error("Expected lock file to exist after acquisition")
	}

	// Test duplicate lock acquisition should fail
	lm2 := &TimerLockManager{
		lockFile: filepath.Join(tmpDir, "test_timer.lock"),
		timerPID: os.Getpid() + 1, // Different PID
	}

	err = lm2.AcquireLock("another-session", 30*time.Minute)
	if err == nil {
		t.Error("Expected duplicate lock acquisition to fail")
	}

	// Clean up
	err = lm.ReleaseLock()
	if err != nil {
		t.Errorf("Failed to release lock: %v", err)
	}
}

func TestTimerLockManager_ReleaseLock(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()

	// Create and acquire lock
	lm := &TimerLockManager{
		lockFile: filepath.Join(tmpDir, "test_timer.lock"),
		timerPID: os.Getpid(),
	}

	err := lm.AcquireLock("test-session", 25*time.Minute)
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}

	// Test lock release
	err = lm.ReleaseLock()
	if err != nil {
		t.Fatalf("Expected successful lock release, got error: %v", err)
	}

	// Verify lock is marked as released
	if lm.locked {
		t.Error("Expected lock to be marked as released")
	}

	// Verify lock file is removed
	if _, err := os.Stat(lm.lockFile); !os.IsNotExist(err) {
		t.Error("Expected lock file to be removed after release")
	}

	// Test releasing already released lock should be safe
	err = lm.ReleaseLock()
	if err != nil {
		t.Errorf("Expected releasing already released lock to be safe, got error: %v", err)
	}
}

func TestTimerLockManager_ReadLockState(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()

	// Create and acquire lock
	lm := &TimerLockManager{
		lockFile: filepath.Join(tmpDir, "test_timer.lock"),
		timerPID: os.Getpid(),
	}

	sessionName := "test-session"
	duration := 25 * time.Minute

	err := lm.AcquireLock(sessionName, duration)
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}

	// Read lock state
	state, err := lm.ReadLockState()
	if err != nil {
		t.Fatalf("Failed to read lock state: %v", err)
	}

	// Verify state contents
	if state.PID != os.Getpid() {
		t.Errorf("Expected PID %d, got %d", os.Getpid(), state.PID)
	}

	if state.SessionName != sessionName {
		t.Errorf("Expected session name %s, got %s", sessionName, state.SessionName)
	}

	if state.Duration != int(duration.Seconds()) {
		t.Errorf("Expected duration %d seconds, got %d", int(duration.Seconds()), state.Duration)
	}

	// Clean up
	err = lm.ReleaseLock()
	if err != nil {
		t.Errorf("Failed to release lock: %v", err)
	}
}

func TestTimerLockManager_CorruptedLockFile(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	lockFile := filepath.Join(tmpDir, "test_timer.lock")

	// Create corrupted lock file
	err := os.WriteFile(lockFile, []byte("invalid json"), 0600)
	if err != nil {
		t.Fatalf("Failed to create corrupted lock file: %v", err)
	}

	// Create lock manager
	lm := &TimerLockManager{
		lockFile: lockFile,
		timerPID: os.Getpid(),
	}

	// Reading corrupted lock file should delete it and return error
	_, err = lm.ReadLockState()
	if err == nil {
		t.Error("Expected error when reading corrupted lock file")
	}

	// Verify corrupted lock file was deleted
	if _, err := os.Stat(lockFile); !os.IsNotExist(err) {
		t.Error("Expected corrupted lock file to be deleted")
	}
}
