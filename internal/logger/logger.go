package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// Logger is the global logger instance
	Logger *logrus.Logger
)

// Config represents logger configuration
type Config struct {
	Level string // debug, info, warn, error
	File  string // Empty = stderr only, or path to log file
}

// Init initializes the global logger with the given configuration
func Init(config Config) error {
	Logger = logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		// Default to info if invalid level
		level = logrus.InfoLevel
		Logger.WithError(err).Warn("Invalid log level, defaulting to info")
	}
	Logger.SetLevel(level)

	// Set output
	var output io.Writer = os.Stderr
	if config.File != "" {
		file, err := os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			Logger.WithError(err).Warn("Failed to open log file, using stderr")
		} else {
			output = file
		}
	}
	Logger.SetOutput(output)

	// Set formatter
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return nil
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	if Logger == nil {
		Init(Config{Level: "info", File: ""})
	}
	Logger.Debug(args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	if Logger == nil {
		Init(Config{Level: "info", File: ""})
	}
	Logger.Debugf(format, args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	if Logger == nil {
		Init(Config{Level: "info", File: ""})
	}
	Logger.Info(args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	if Logger == nil {
		Init(Config{Level: "info", File: ""})
	}
	Logger.Infof(format, args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	if Logger == nil {
		Init(Config{Level: "info", File: ""})
	}
	Logger.Warn(args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	if Logger == nil {
		Init(Config{Level: "info", File: ""})
	}
	Logger.Warnf(format, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	if Logger == nil {
		Init(Config{Level: "info", File: ""})
	}
	Logger.Error(args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	if Logger == nil {
		Init(Config{Level: "info", File: ""})
	}
	Logger.Errorf(format, args...)
}

// WithField creates a logger entry with a field
func WithField(key string, value interface{}) *logrus.Entry {
	if Logger == nil {
		// Initialize with defaults if not already initialized
		Init(Config{Level: "info", File: ""})
	}
	return Logger.WithField(key, value)
}

// WithFields creates a logger entry with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	if Logger == nil {
		// Initialize with defaults if not already initialized
		Init(Config{Level: "info", File: ""})
	}
	return Logger.WithFields(fields)
}

// WithError creates a logger entry with an error
func WithError(err error) *logrus.Entry {
	if Logger == nil {
		// Initialize with defaults if not already initialized
		Init(Config{Level: "info", File: ""})
	}
	return Logger.WithError(err)
}

