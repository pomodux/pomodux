package main

import (
	"fmt"
	"os"

	"github.com/pomodux/pomodux/internal/config"
	"github.com/pomodux/pomodux/internal/logger"
	"github.com/spf13/cobra"
)

var (
	version   = "0.1.0"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	var limit int
	var today bool
	var all bool

	rootCmd := &cobra.Command{
		Use:     "pomodux-stats",
		Short:   "View timer statistics",
		Long:    "View statistics and history for pomodoro timer sessions",
		Version: fmt.Sprintf("%s (built %s, commit %s)", version, buildTime, gitCommit),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showStats(limit, today, all)
		},
	}

	rootCmd.Flags().IntVarP(&limit, "limit", "l", 20, "Show last N sessions")
	rootCmd.Flags().BoolVarP(&today, "today", "t", false, "Show today's statistics")
	rootCmd.Flags().BoolVar(&all, "all", false, "Show all sessions")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func showStats(limit int, today bool, all bool) error {
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
	fmt.Printf("Options: limit=%d, today=%v, all=%v\n", limit, today, all)

	return nil
}


