package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/pomodux/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a paused timer",
	Long:  `Resume a paused timer. The timer will continue counting down from where it was paused. Use 'pomodux pause' to pause again if needed.`,
	RunE:  runResume,
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}

func runResume(cmd *cobra.Command, args []string) error {
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
				return fmt.Errorf("cannot resume timer running in another process (PID %d): %s session with %s remaining",
					lockState.PID, lockState.SessionName, formatDuration(remaining))
			}
			return fmt.Errorf("timer in another process (PID %d) may have completed. Use 'pomodux status' for details", lockState.PID)
		}
	}

	t := timer.GetGlobalTimer()

	status := t.GetStatus()
	if status != timer.StatusPaused {
		cmd.PrintErrln("Cannot resume timer: timer is not paused (current status:", status, ")")
		return fmt.Errorf("cannot resume timer: timer is not paused (current status: %v)", status)
	}

	if err := t.Resume(); err != nil {
		cmd.PrintErrln("Failed to resume timer:", err)
		return fmt.Errorf("failed to resume timer: %w", err)
	}

	fmt.Println("Timer resumed.")
	return nil
}
