package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pomodux/pomodux/internal/config"
	"github.com/pomodux/pomodux/internal/logger"
	"github.com/pomodux/pomodux/internal/timer"
	"github.com/pomodux/pomodux/internal/tui"
)

// Version information embedded at build time
var (
	Version   = "dev"     // Set via -ldflags during build
	GitCommit = "unknown" // Set via -ldflags during build
	BuildDate = "unknown" // Set via -ldflags during build
)

func main() {
	// Parse command line arguments
	args := os.Args[1:]

	// Handle version command
	if len(args) > 0 && (args[0] == "version" || args[0] == "--version" || args[0] == "-v") {
		fmt.Printf("Pomodux %s (commit: %s, built: %s)\n", Version, GitCommit, BuildDate)
		os.Exit(0)
	}

	// Handle help command
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		printUsage()
		os.Exit(0)
	}

	// Load configuration first to get defaults
	cfg, err := loadConfig(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := initializeLogger(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	// Parse command line arguments for duration and session name
	duration, sessionName, err := parseTimerArgs(args, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		printUsage()
		os.Exit(1)
	}

	// Ensure graceful shutdown
	defer timer.ShutdownGlobalTimer()

	// Launch TUI immediately with timer arguments (TUI instantiates timer singleton)
	if err := tui.RunWithArgs(duration, sessionName); err != nil {
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}
}

// loadConfig loads configuration from file or returns default
func loadConfig(args []string) (*config.Config, error) {
	// Check for --config flag
	configFile := ""
	for i, arg := range args {
		if arg == "--config" && i+1 < len(args) {
			configFile = args[i+1]
			break
		} else if strings.HasPrefix(arg, "--config=") {
			configFile = strings.TrimPrefix(arg, "--config=")
			break
		}
	}

	if configFile != "" {
		return config.LoadFromPath(configFile)
	}
	return config.Load()
}

// initializeLogger sets up logging based on configuration
func initializeLogger(cfg *config.Config) error {
	logConfig := &logger.Config{
		Level:      logger.LogLevel(cfg.Logging.Level),
		Format:     cfg.Logging.Format,
		Output:     cfg.Logging.Output,
		LogFile:    cfg.Logging.LogFile,
		ShowCaller: cfg.Logging.ShowCaller,
	}
	return logger.Init(logConfig)
}

// parseTimerArgs parses command line arguments for timer duration and session name
func parseTimerArgs(args []string, cfg *config.Config) (time.Duration, string, error) {
	// Filter out config-related args and legacy "start" command
	timerArgs := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--config" {
			i++ // skip the next argument (config file path)
			continue
		}
		if strings.HasPrefix(arg, "--config=") {
			continue
		}
		// Skip legacy "start" command for backwards compatibility
		if arg == "start" {
			continue
		}
		timerArgs = append(timerArgs, arg)
	}

	// Default values from config
	duration := cfg.Timer.DefaultWorkDuration
	sessionName := cfg.Timer.DefaultSessionName
	if sessionName == "" {
		sessionName = "work" // fallback
	}

	// Parse duration argument
	if len(timerArgs) > 0 {
		var err error
		duration, err = parseDuration(timerArgs[0])
		if err != nil {
			return 0, "", fmt.Errorf("invalid duration '%s': %v", timerArgs[0], err)
		}
	}

	// Parse session name argument
	if len(timerArgs) > 1 {
		sessionName = timerArgs[1]
	}

	if duration <= 0 {
		return 0, "", fmt.Errorf("duration must be positive")
	}

	return duration, sessionName, nil
}

// parseDuration parses duration strings like "25m", "1h30m", "90m", etc.
func parseDuration(durationStr string) (time.Duration, error) {
	// Try parsing as Go duration first (e.g., "25m", "1h30m")
	if duration, err := time.ParseDuration(durationStr); err == nil {
		return duration, nil
	}

	// Try parsing as plain number (assume minutes)
	if minutes, err := strconv.Atoi(durationStr); err == nil {
		return time.Duration(minutes) * time.Minute, nil
	}

	return 0, fmt.Errorf("unable to parse duration")
}

// printUsage prints usage information
func printUsage() {
	fmt.Printf(`Pomodux %s - Terminal-based timer application

USAGE:
    pomodux [duration] [session_name] [flags]

ARGUMENTS:
    duration      Timer duration (e.g., 25m, 1h30m, 90) [default: from config]
    session_name  Name for this session [default: from config]

FLAGS:
    --config FILE    Configuration file path
    --version, -v    Show version information
    --help, -h       Show this help message

EXAMPLES:
    pomodux                           # Start with default duration and session name
    pomodux 25m                       # 25-minute session with default name  
    pomodux 25m "work"                # 25-minute "work" session
    pomodux 5m "break"                # 5-minute "break" session
    pomodux 1h30m "deep focus"        # 1.5-hour custom session
    pomodux 90 "coding"               # 90-minute session (number = minutes)

DESCRIPTION:
    Pomodux launches an interactive TUI timer that supports:
    • Real-time progress display with smooth animations
    • Keyboard controls: [P]ause, [R]esume, [S]top, [Q]uit
    • Automatic session recording and history tracking
    • Plugin system for notifications and extensions
    • Single timer enforcement (only one timer at a time)

    The TUI provides immediate visual feedback and all timer controls
    are available within the interface. No separate CLI commands needed.

`, Version)
}
