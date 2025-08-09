package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// expandPath expands environment variables and tilde in a path
func expandPath(path string) string {
	// Expand environment variables first
	expanded := os.ExpandEnv(path)

	// Expand tilde to home directory
	if strings.HasPrefix(expanded, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// If we can't get home directory, return as-is
			return expanded
		}
		expanded = filepath.Join(homeDir, expanded[2:])
	}

	return expanded
}

// Config represents the Pomodux configuration structure
type Config struct {
	Timer struct {
		DefaultWorkDuration      time.Duration `yaml:"default_work_duration"`
		DefaultBreakDuration     time.Duration `yaml:"default_break_duration"`
		DefaultLongBreakDuration time.Duration `yaml:"default_long_break_duration"`
		DefaultSessionName       string        `yaml:"default_session_name"`
		AutoStartBreaks          bool          `yaml:"auto_start_breaks"`
	} `yaml:"timer"`

	TUI struct {
		Theme       string            `yaml:"theme"`
		KeyBindings map[string]string `yaml:"key_bindings"`
	} `yaml:"tui"`

	Notifications struct {
		Enabled bool   `yaml:"enabled"`
		Sound   bool   `yaml:"sound"`
		Message string `yaml:"message"`
	} `yaml:"notifications"`

	Plugins struct {
		Directory string          `yaml:"directory"`
		Enabled   map[string]bool `yaml:"enabled"`
	} `yaml:"plugins"`

	// Raw plugins section for plugin-specific config
	PluginsRaw map[string]interface{} // no yaml tag

	Logging struct {
		Level      string `yaml:"level"`
		Format     string `yaml:"format"`
		Output     string `yaml:"output"`
		LogFile    string `yaml:"log_file"`
		ShowCaller bool   `yaml:"show_caller"`
	} `yaml:"logging"`
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	config := &Config{}

	// Timer defaults
	config.Timer.DefaultWorkDuration = 25 * time.Minute
	config.Timer.DefaultBreakDuration = 5 * time.Minute
	config.Timer.DefaultLongBreakDuration = 15 * time.Minute
	config.Timer.DefaultSessionName = "work"
	config.Timer.AutoStartBreaks = false

	// TUI defaults
	config.TUI.Theme = "default"
	config.TUI.KeyBindings = map[string]string{
		"start":  "s",
		"stop":   "q",
		"pause":  "p",
		"resume": "r",
	}

	// Notification defaults
	config.Notifications.Enabled = true
	config.Notifications.Sound = false
	config.Notifications.Message = "Timer completed!"

	// Plugins directory default
	config.Plugins.Directory = defaultPluginsDir()
	config.Plugins.Enabled = make(map[string]bool) // Initialize Enabled map

	// Logging defaults
	config.Logging.Level = "info"
	config.Logging.Format = "text"
	config.Logging.Output = "file" // Changed from "console" to "file"
	config.Logging.LogFile = ""
	config.Logging.ShowCaller = false

	return config
}

// defaultPluginsDir returns the default plugins directory (XDG_CONFIG_HOME/pomodux/plugins)
func defaultPluginsDir() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "./plugins" // fallback
		}
		configHome = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configHome, "pomodux", "plugins")
}

// Load loads configuration from the default XDG location
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// If config file doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := Save(config); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}

	// Load existing config
	data, err := os.ReadFile(configPath) // #nosec G304 -- configPath is from trusted XDG_CONFIG_HOME or user home
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Unmarshal the plugins section into PluginsRaw
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err == nil {
		if pluginsSection, ok := raw["plugins"]; ok {
			if pluginsMap, ok := pluginsSection.(map[string]interface{}); ok {
				config.PluginsRaw = pluginsMap
			}
		}
	}

	// Ensure plugins directory is set and expand paths
	if config.Plugins.Directory == "" {
		config.Plugins.Directory = defaultPluginsDir()
	} else {
		config.Plugins.Directory = expandPath(config.Plugins.Directory)
	}

	// Validate configuration
	if err := Validate(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Save saves configuration to the default XDG location
func Save(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// SaveToPath saves configuration to a specific path
func SaveToPath(config *Config, path string) error {
	// Ensure directory exists
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadFromPath loads configuration from a specific path
func LoadFromPath(path string) (*Config, error) {
	// Validate path is not empty and doesn't contain path traversal
	if path == "" || strings.Contains(path, "..") {
		return nil, fmt.Errorf("invalid config path: %s", path)
	}
	data, err := os.ReadFile(path) // #nosec G304 -- path is validated above
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Ensure plugins directory is set and expand paths
	if config.Plugins.Directory == "" {
		config.Plugins.Directory = defaultPluginsDir()
	} else {
		config.Plugins.Directory = expandPath(config.Plugins.Directory)
	}

	// Validate configuration
	if err := Validate(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// getConfigPath returns the path to the configuration file
func getConfigPath() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configHome = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configHome, "pomodux", "config.yaml"), nil
}

// Validate validates the configuration
func Validate(config *Config) error {
	// Validate timer configuration
	if config.Timer.DefaultWorkDuration <= 0 {
		return fmt.Errorf("default work duration must be positive")
	}
	if config.Timer.DefaultBreakDuration <= 0 {
		return fmt.Errorf("default break duration must be positive")
	}
	if config.Timer.DefaultLongBreakDuration <= 0 {
		return fmt.Errorf("default long break duration must be positive")
	}
	if config.Timer.DefaultSessionName == "" {
		return fmt.Errorf("default session name cannot be empty")
	}

	// Validate logging configuration
	if config.Logging.Level != "" {
		validLevels := map[string]bool{
			"debug": true,
			"info":  true,
			"warn":  true,
			"error": true,
		}
		if !validLevels[config.Logging.Level] {
			return fmt.Errorf("invalid log level: %s", config.Logging.Level)
		}
	}

	if config.Logging.Format != "" {
		validFormats := map[string]bool{
			"text": true,
			"json": true,
		}
		if !validFormats[config.Logging.Format] {
			return fmt.Errorf("invalid log format: %s", config.Logging.Format)
		}
	}

	if config.Logging.Output != "" {
		validOutputs := map[string]bool{
			"console": true,
			"file":    true,
			"both":    true,
		}
		if !validOutputs[config.Logging.Output] {
			return fmt.Errorf("invalid log output: %s", config.Logging.Output)
		}
	}

	return nil
}
