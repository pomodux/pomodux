package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "1.0", config.Version)
	assert.Equal(t, "25m", config.Timers["work"])
	assert.Equal(t, "5m", config.Timers["break"])
	assert.Equal(t, "15m", config.Timers["longbreak"])
	assert.Equal(t, "default", config.Theme)
	assert.False(t, config.Timer.BellOnComplete)
	assert.Equal(t, "info", config.Logging.Level)
	assert.Empty(t, config.Logging.File)
}

func TestLoadFromPath_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent.yaml")

	config, err := LoadFromPath(configPath)
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "1.0", config.Version)
}

func TestLoadFromPath_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `
version: "1.0"
timers:
  work: 30m
  break: 10m
theme: "nord"
timer:
  bell_on_complete: true
logging:
  level: "debug"
  file: ""
`

	err := os.WriteFile(configPath, []byte(yamlContent), 0600)
	require.NoError(t, err)

	config, err := LoadFromPath(configPath)
	require.NoError(t, err)
	assert.Equal(t, "30m", config.Timers["work"])
	assert.Equal(t, "10m", config.Timers["break"])
	assert.Equal(t, "nord", config.Theme)
	assert.True(t, config.Timer.BellOnComplete)
	assert.Equal(t, "debug", config.Logging.Level)
}

func TestSaveToPath(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	config := DefaultConfig()
	config.Timers["work"] = "30m"

	err := SaveToPath(config, configPath)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Load and verify
	loaded, err := LoadFromPath(configPath)
	require.NoError(t, err)
	assert.Equal(t, "30m", loaded.Timers["work"])
}

func TestConfigPath(t *testing.T) {
	path := ConfigPath()
	assert.Contains(t, path, "pomodux")
	assert.Contains(t, path, "config.yaml")
}

func TestStatePath(t *testing.T) {
	path := StatePath()
	assert.Contains(t, path, "pomodux")
}

func TestHistoryPath(t *testing.T) {
	path := HistoryPath()
	assert.Contains(t, path, "history.json")
}

func TestTimerStatePath(t *testing.T) {
	path := TimerStatePath()
	assert.Contains(t, path, "timer_state.json")
}

