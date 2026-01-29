package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pomodux/pomodux/internal/logger"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Version string            `yaml:"version"`
	Timers  map[string]string `yaml:"timers"`
	Theme   string            `yaml:"theme"`
	Timer   TimerConfig       `yaml:"timer"`
	Logging LoggingConfig     `yaml:"logging"`
	Plugins PluginsConfig     `yaml:"plugins"`
}

// TimerConfig represents timer-specific configuration
type TimerConfig struct {
	BellOnComplete bool `yaml:"bell_on_complete"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

// PluginsConfig represents plugin configuration (Post-MVP)
type PluginsConfig struct {
	Enabled   []string `yaml:"enabled"`
	Directory string   `yaml:"directory"`
}

// Load loads configuration from the XDG-compliant config file location
func Load() (*Config, error) {
	configPath := ConfigPath()
	return LoadFromPath(configPath)
}

// LoadFromPath loads configuration from a specific file path
func LoadFromPath(path string) (*Config, error) {
	logger.WithField("config_path", path).Debug("Loading configuration")

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Info("Config file not found, creating default configuration")
		logger.WithField("config_path", path).Debug("Config will be saved to path")

		defaults := DefaultConfig()

		// Try to save default config to disk
		if err := SaveToPath(defaults, path); err != nil {
			logger.WithError(err).Warn("Could not create config file, using in-memory defaults")
			// Don't fail - return defaults even if we couldn't save
			return defaults, nil
		}

		logger.Info("Default config file created successfully")
		return defaults, nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// Validate and apply defaults
	if err := validateAndApplyDefaults(&config); err != nil {
		logger.WithError(err).Warn("Config validation failed, using defaults for invalid fields")
		applyDefaults(&config)
	}

	logger.Debug("Configuration loaded successfully")
	return &config, nil
}

// Save saves configuration to the XDG-compliant config file location
func Save(config *Config) error {
	configPath := ConfigPath()
	return SaveToPath(config, configPath)
}

// SaveToPath saves configuration to a specific file path
func SaveToPath(config *Config, path string) error {
	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML with proper indentation
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add a helpful comment at the top
	header := "# Pomodux Configuration File\n# Edit this file to customize your timer settings\n\n"
	data = append([]byte(header), data...)

	// Write file with proper permissions
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	logger.WithField("config_path", path).Debug("Configuration saved")
	return nil
}

// ConfigPath returns the XDG-compliant config file path
func ConfigPath() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "pomodux", "config.yaml")
}

// StatePath returns the XDG-compliant state directory path
func StatePath() string {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		home, _ := os.UserHomeDir()
		stateHome = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(stateHome, "pomodux")
}

// HistoryPath returns the path to the history file
func HistoryPath() string {
	return filepath.Join(StatePath(), "history.json")
}

// TimerStatePath returns the path to the timer state file
func TimerStatePath() string {
	return filepath.Join(StatePath(), "timer_state.json")
}

// LogFilePath returns the XDG-compliant default log file path
// Uses $XDG_CACHE_HOME/pomodux/pomodux.log or ~/.cache/pomodux/pomodux.log as fallback
func LogFilePath() string {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			// Fallback if home directory can't be determined
			return filepath.Join(".cache", "pomodux", "pomodux.log")
		}
		cacheHome = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheHome, "pomodux", "pomodux.log")
}

// validateAndApplyDefaults validates the config and applies defaults for missing fields
func validateAndApplyDefaults(config *Config) error {
	// Apply defaults for missing fields
	applyDefaults(config)

	// Validate theme
	if config.Theme == "" {
		config.Theme = "default"
	}

	// Validate logging level
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[config.Logging.Level] {
		logger.Warnf("Invalid log level %q, defaulting to info", config.Logging.Level)
		config.Logging.Level = "info"
	}

	return nil
}

// applyDefaults applies default values to missing config fields
func applyDefaults(config *Config) {
	defaults := DefaultConfig()

	if config.Version == "" {
		config.Version = defaults.Version
	}

	if config.Timers == nil || len(config.Timers) == 0 {
		config.Timers = defaults.Timers
	}

	if config.Theme == "" {
		config.Theme = defaults.Theme
	}

	if config.Logging.Level == "" {
		config.Logging.Level = defaults.Logging.Level
	}

	if config.Logging.File == "" {
		config.Logging.File = defaults.Logging.File
	}

	if config.Plugins.Enabled == nil {
		config.Plugins.Enabled = defaults.Plugins.Enabled
	}

	if config.Plugins.Directory == "" {
		config.Plugins.Directory = defaults.Plugins.Directory
	}
}

