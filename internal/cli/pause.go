package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/pomodux/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause the currently running timer",
	Long: `Pause the currently running timer. The timer will stop counting down
but retain its current progress. Use 'pomodux resume' to continue the timer.`,
	RunE: runPause,
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}

func runPause(cmd *cobra.Command, args []string) error {
	// Check for lock conflicts first
	lockManager, err := timer.NewTimerLockManager()
	if err != nil {
		return fmt.Errorf("failed to initialize lock manager: %v", err)
	}

	// Try to read lock state to provide better error messages
	lockState, lockErr := lockManager.ReadLockState()
	if lockErr == nil {
		currentPID := os.Getpid()
		if lockState.PID != currentPID {
			// Another process is running the timer
			remaining := time.Duration(lockState.Duration)*time.Second - time.Since(lockState.StartTime)
			if remaining > 0 {
				return fmt.Errorf("cannot pause timer running in another process (PID %d): %s session with %s remaining",
					lockState.PID, lockState.SessionName, formatDuration(remaining))
			}
			return fmt.Errorf("timer in another process (PID %d) may have completed. Use 'pomodux status' for details", lockState.PID)
		}
	}

	t := timer.GetGlobalTimer()

	status := t.GetStatus()
	if status != timer.StatusRunning {
		cmd.PrintErrln("Cannot pause timer: timer is not running (current status:", status, ")")
		return fmt.Errorf("cannot pause timer: timer is not running (current status: %v)", status)
	}

	if err := t.Pause(); err != nil {
		cmd.PrintErrln("Failed to pause timer:", err)
		return fmt.Errorf("failed to pause timer: %w", err)
	}
	fmt.Println("Timer paused")
	return nil
}
