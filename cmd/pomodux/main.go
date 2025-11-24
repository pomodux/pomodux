package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pomodux/pomodux/internal/config"
	"github.com/pomodux/pomodux/internal/logger"
	"github.com/pomodux/pomodux/internal/timer"
	"github.com/urfave/cli/v2"
)

var (
	version   = "0.1.0"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	app := &cli.App{
		Name:    "pomodux",
		Usage:   "Terminal-based Pomodoro timer",
		Version: fmt.Sprintf("%s (built %s, commit %s)", version, buildTime, gitCommit),
		Commands: []*cli.Command{
			{
				Name:      "start",
				Usage:     "Start a timer session",
				UsageText: "pomodux start <duration|preset> [label]",
				Action:    startTimer,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func startTimer(c *cli.Context) error {
	args := c.Args()
	if args.Len() == 0 {
		return cli.Exit("Error: duration or preset required\n"+
			"Usage: pomodux start <duration|preset> [label]", 2)
	}

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
	durationOrPreset := args.Get(0)
	label := args.Get(1)

	var duration time.Duration
	var preset string

	// Try to parse as duration first
	duration, err = time.ParseDuration(durationOrPreset)
	if err != nil {
		// Not a duration, try as preset
		presetDuration, ok := cfg.Timers[durationOrPreset]
		if !ok {
			return cli.Exit(fmt.Sprintf("Error: unknown preset %q\n"+
				"Available presets: %v", durationOrPreset, getPresetNames(cfg.Timers)), 2)
		}

		duration, err = time.ParseDuration(presetDuration)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Error: invalid duration in preset %q: %v", durationOrPreset, err), 2)
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
		return cli.Exit(fmt.Sprintf("Error: %v", err), 1)
	}

	if err := t.Start(); err != nil {
		return cli.Exit(fmt.Sprintf("Error: failed to start timer: %v", err), 1)
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

