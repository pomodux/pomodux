package config

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Version: "1.0",
		Timers: map[string]string{
			"work":      "25m",
			"break":     "5m",
			"longbreak": "15m",
		},
		Theme: "default",
		Timer: TimerConfig{
			BellOnComplete: false,
		},
		Logging: LoggingConfig{
			Level: "info",
			File:  "",
		},
		Plugins: PluginsConfig{
			Enabled:  []string{},
			Directory: "",
		},
	}
}

