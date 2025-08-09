package cli

import (
	"fmt"
	"time"

	"github.com/pomodux/pomodux/internal/config"
	"github.com/pomodux/pomodux/internal/timer"
	"github.com/pomodux/pomodux/internal/tui"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [duration] [session_name]",
	Short: "Start a timer with optional session name",
	Long: `Start a timer for the specified duration and session name.
	
Examples:
  pomodux start 25m              # Start a 25-minute session with default name
  pomodux start 25m "work"       # Start a 25-minute "work" session
  pomodux start 5m "break"       # Start a 5-minute "break" session
  pomodux start 1h30m "deep work" # Start a 1.5-hour custom session
  
If no duration is specified, uses the default work duration from config.
If no session name is specified, uses the default session name from config.`,
	Args: cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration with hot-reload
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %v", err)
		}

		// Parse duration argument
		var duration time.Duration
		if len(args) > 0 {
			duration, err = time.ParseDuration(args[0])
			if err != nil {
				return fmt.Errorf("invalid duration: %v", err)
			}
		} else {
			// Use default work duration from config
			duration = cfg.Timer.DefaultWorkDuration
		}

		if duration <= 0 {
			return fmt.Errorf("duration must be positive")
		}

		// Parse session name argument
		var sessionName string
		if len(args) > 1 {
			sessionName = args[1]
		} else {
			// Use default session name from config
			sessionName = cfg.Timer.DefaultSessionName
		}

		if sessionName == "" {
			sessionName = "work" // fallback if config is empty
		}

		// Initialize lock manager
		lockManager, err := timer.NewTimerLockManager()
		if err != nil {
			return fmt.Errorf("failed to initialize lock manager: %v", err)
		}

		// Attempt to acquire lock
		err = lockManager.AcquireLock(sessionName, duration)
		if err != nil {
			// Handle timer conflict error specifically
			if conflictErr, ok := err.(*timer.TimerConflictError); ok {
				return fmt.Errorf("timer already running: %v", conflictErr.Error())
			}
			return fmt.Errorf("failed to acquire timer lock: %v", err)
		}

		// Ensure lock is released when command exits
		defer func() {
			if releaseErr := lockManager.ReleaseLock(); releaseErr != nil {
				cmd.PrintErrln("Warning: Failed to release timer lock:", releaseErr)
			}
		}()

		// Get the global timer
		t := timer.GetGlobalTimer()

		// Start the timer using the new session name API
		err = t.StartWithSessionName(duration, sessionName)
		if err != nil {
			return fmt.Errorf("failed to start timer: %v", err)
		}

		// Launch TUI immediately (blocking)
		err = tui.RunTUI(t)
		if err != nil {
			return fmt.Errorf("TUI error: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
