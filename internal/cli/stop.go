package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/pomodux/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the currently running timer",
	Long: `Stop the currently running timer.
	
This command will stop the timer and record the session as interrupted.
The timer must be running (not paused or completed) for this command to work.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
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
					return fmt.Errorf("cannot stop timer running in another process (PID %d): %s session with %s remaining",
						lockState.PID, lockState.SessionName, formatDuration(remaining))
				}
				return fmt.Errorf("timer in another process (PID %d) may have completed. Use 'pomodux status' for details", lockState.PID)
			}
		}

		// Get the global timer
		t := timer.GetGlobalTimer()

		// Check if timer is running
		if t.GetStatus() != timer.StatusRunning {
			return fmt.Errorf("timer is not running (current status: %s)", t.GetStatus())
		}

		// Stop the timer
		if err := t.Stop(); err != nil {
			return fmt.Errorf("failed to stop timer: %v", err)
		}

		fmt.Println("⏹️  Timer stopped.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
