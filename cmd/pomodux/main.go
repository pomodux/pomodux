package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pomodux/pomodux/internal/config"
	"github.com/pomodux/pomodux/internal/logger"
	"github.com/pomodux/pomodux/internal/timer"
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

	logger.WithField("component", "pomodux").Info("Starting pomodux")

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

	// Create and start timer
	t, err := timer.NewTimer(duration, label, preset)
	if err != nil {
		return fmt.Errorf("failed to create timer: %w", err)
	}

	if err := t.Start(); err != nil {
		return fmt.Errorf("failed to start timer: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"duration": duration,
		"label":    label,
		"preset":   preset,
	}).Info("Timer started")

	// TODO: Start TUI
	fmt.Printf("Timer started for %v with label: %s\n", duration, label)
	fmt.Println("TUI will be implemented in the next phase")

	return nil
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
