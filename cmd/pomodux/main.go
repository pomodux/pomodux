package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/pomodux/pomodux/internal/config"
	"github.com/pomodux/pomodux/internal/logger"
	"github.com/pomodux/pomodux/internal/theme"
	"github.com/pomodux/pomodux/internal/timer"
	"github.com/pomodux/pomodux/internal/tui"
	"github.com/spf13/cobra"
)

var (
	version   = "0.1.0"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "pomodux",
		Short:   "Terminal-based Pomodoro timer",
		Long:    "A terminal-based Pomodoro timer application",
		Version: fmt.Sprintf("%s (built %s, commit %s)", version, buildTime, gitCommit),
	}

	startCmd := &cobra.Command{
		Use:   "start <duration|preset> [label]",
		Short: "Start a timer session",
		Long:  "Start a timer session with a duration (e.g., 25m, 1h30m) or preset name, with an optional label",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  startTimer,
	}

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func startTimer(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	if err := logger.Init(logger.Config{
		Level: cfg.Logging.Level,
		File:  cfg.Logging.File,
	}); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Resolve theme from config (fallback to default for unknown names)
	selectedTheme := theme.GetTheme(cfg.Theme)
	if cfg.Theme != "default" && cfg.Theme != "nord" && cfg.Theme != "catppuccin-mocha" {
		logger.WithField("theme", cfg.Theme).Warn("Unknown theme name, using default theme")
	}

	logger.WithFields(map[string]interface{}{
		"component": "pomodux",
		"theme":     selectedTheme.Name,
	}).Infof("Starting pomodux (theme: %s)", selectedTheme.Name)

	// Check for existing timer state (singleton enforcement)
	statePath := config.TimerStatePath()
	var t *timer.Timer
	var sessionID string

	if stateExists(statePath) {
		state, err := timer.LoadState(statePath)
		if err != nil {
			logger.WithError(err).Warn("Failed to load existing timer state, starting new timer")
			// Continue with new timer creation
		} else {
			// Check if process is still alive
			if state.PID > 0 && timer.IsProcessAlive(state.PID) {
				return fmt.Errorf("timer already running in process %d", state.PID)
			}

			// Process is dead, auto-resume
			logger.WithField("session_id", state.SessionID).Info("Resuming interrupted timer")
			t, err = timer.ResumeFromState(state)
			if err != nil {
				return fmt.Errorf("failed to resume timer: %w", err)
			}
			sessionID = state.SessionID
		}
	}

	// If no timer was resumed, create a new one
	if t == nil {
		// Parse duration or preset
		durationOrPreset := args[0]
		var label string
		if len(args) > 1 {
			label = args[1]
		}

		var duration time.Duration
		var preset string

		// Try to parse as duration first
		duration, err = time.ParseDuration(durationOrPreset)
		if err != nil {
			// Not a duration, try as preset
			presetDuration, ok := cfg.Timers[durationOrPreset]
			if !ok {
				return fmt.Errorf("unknown preset %q\nAvailable presets: %v", durationOrPreset, getPresetNames(cfg.Timers))
			}

			duration, err = time.ParseDuration(presetDuration)
			if err != nil {
				return fmt.Errorf("invalid duration in preset %q: %w", durationOrPreset, err)
			}
			preset = durationOrPreset
		}

		// Default label if not provided
		if label == "" {
			if preset != "" {
				label = prettifyPresetName(preset)
			} else {
				label = "Generic timer session"
			}
		}

		// Generate session ID
		sessionID = uuid.New().String()

		// Create and start timer
		t, err = timer.NewTimer(duration, label, preset)
		if err != nil {
			return fmt.Errorf("failed to create timer: %w", err)
		}

		if err := t.Start(); err != nil {
			return fmt.Errorf("failed to start timer: %w", err)
		}

		logger.WithFields(map[string]interface{}{
			"session_id": sessionID,
			"duration":   duration,
			"label":      label,
			"preset":     preset,
		}).Info("Timer started")
	}

	// Save initial state
	if err := timer.SaveState(t, sessionID, statePath); err != nil {
		logger.WithError(err).Warn("Failed to save initial timer state")
	}

	// Redirect logger to file if not configured, to prevent TUI interference
	if cfg.Logging.File == "" {
		logPath := config.LogFilePath()
		if err := logger.RedirectToFile(logPath); err != nil {
			logger.WithError(err).WithField("log_path", logPath).Warn("Failed to redirect logger to file, TUI may be corrupted by log output")
		} else {
			logger.WithField("log_path", logPath).Info("Logger redirected to file for TUI compatibility")
		}
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Start TUI with resolved theme
	model := tui.NewModel(t, sessionID, statePath, selectedTheme)
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Handle signals in a goroutine (this is the exception - signal handling)
	go func() {
		sig := <-sigChan
		logger.WithField("signal", sig.String()).Info("Received interrupt signal, saving state and exiting")

		// Save state before exit
		if err := timer.SaveState(t, sessionID, statePath); err != nil {
			logger.WithError(err).Error("Failed to save state on interrupt")
		}

		program.Quit()
	}()

	// Run TUI
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	// Clean up state file on normal exit
	if err := timer.DeleteState(statePath); err != nil {
		logger.WithError(err).Warn("Failed to delete state file")
	}

	return nil
}

// stateExists checks if the state file exists
func stateExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getPresetNames(timers map[string]string) []string {
	names := make([]string, 0, len(timers))
	for name := range timers {
		names = append(names, name)
	}
	return names
}

func prettifyPresetName(preset string) string {
	// Simple capitalization for now
	if len(preset) == 0 {
		return preset
	}
	return string(preset[0]-32) + preset[1:]
}
