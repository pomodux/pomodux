package main

import (
	"fmt"
	"os"

	"github.com/pomodux/pomodux/internal/config"
	"github.com/pomodux/pomodux/internal/logger"
	"github.com/urfave/cli/v2"
)

var (
	version   = "0.1.0"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	app := &cli.App{
		Name:    "pomodux-stats",
		Usage:   "View timer statistics",
		Version: fmt.Sprintf("%s (built %s, commit %s)", version, buildTime, gitCommit),
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Show last N sessions",
				Value:   20,
			},
			&cli.BoolFlag{
				Name:    "today",
				Aliases: []string{"t"},
				Usage:   "Show today's statistics",
			},
			&cli.BoolFlag{
				Name:  "all",
				Usage: "Show all sessions",
			},
		},
		Action: showStats,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func showStats(c *cli.Context) error {
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

	logger.WithField("component", "pomodux-stats").Info("Starting pomodux-stats")

	// TODO: Load and display history
	fmt.Println("Statistics view will be implemented in the next phase")
	fmt.Printf("Options: limit=%d, today=%v, all=%v\n",
		c.Int("limit"), c.Bool("today"), c.Bool("all"))

	return nil
}

