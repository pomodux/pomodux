package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pomodux/pomodux/internal/logger"
	"github.com/pomodux/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current timer status",
	Long:  `Show detailed information about the current timer session, including time remaining, session type, and progress.`,
	RunE:  runStatus,
}

var statusJSON bool

func init() {
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "Output status as JSON")
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Check for lock conflicts first
	lockManager, err := timer.NewTimerLockManager()
	if err != nil {
		return fmt.Errorf("failed to initialize lock manager: %v", err)
	}

	// Try to read lock state to provide better status information
	lockState, lockErr := lockManager.ReadLockState()
	if lockErr == nil {
		// There's a lock file, check if it's for another process
		currentPID := os.Getpid()
		if lockState.PID != currentPID {
			// Another process may be running a timer
			remaining := time.Duration(lockState.Duration)*time.Second - time.Since(lockState.StartTime)
			if remaining > 0 {
				return fmt.Errorf("timer running in another process (PID %d): %s session with %s remaining",
					lockState.PID, lockState.SessionName, formatDuration(remaining))
			}
			cmd.PrintErrln("Warning: Timer in another process may have completed, run 'pomodux start' to begin new session")
		}
	}

	t := timer.GetGlobalTimer()

	status := t.GetStatus()
	progress := t.GetProgress()
	sessionName := t.GetSessionName()
	startTime := t.GetStartTime()
	duration := t.GetDuration()
	elapsed := t.GetElapsed()

	remaining := duration - elapsed
	if remaining < 0 {
		remaining = 0
	}

	statusInfo := map[string]interface{}{
		"status":       status,
		"session_name": sessionName,
		"start_time":   startTime.Format(time.RFC3339),
		"duration":     duration.Seconds(),
		"elapsed":      elapsed.Seconds(),
		"remaining":    remaining.Seconds(),
		"progress":     progress,
	}

	if statusJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(statusInfo)
	}

	logger.Debug("Status info", statusInfo)

	fmt.Printf("Status:        %s\n", status)
	fmt.Printf("Session Name:  %s\n", sessionName)
	fmt.Printf("Start Time:    %s\n", startTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Duration:      %s\n", formatDuration(duration))
	fmt.Printf("Elapsed:       %s\n", formatDuration(elapsed))
	fmt.Printf("Remaining:     %s\n", formatDuration(remaining))
	fmt.Printf("Progress:      %3.0f%%\n", progress*100)

	return nil
}
