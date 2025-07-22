package timer

import (
	"sync"
	"testing"

	"github.com/pomodux/pomodux/internal/config"
)

func TestSetGlobalConfig(t *testing.T) {
	// Create a test config
	testConfig := &config.Config{}
	testConfig.Plugins.Directory = "/test/plugins"
	testConfig.Plugins.Enabled = map[string]bool{
		"test-plugin": true,
	}

	// Set the global config
	SetGlobalConfig(testConfig)

	// Verify that the global config was set
	if globalConfig == nil {
		t.Fatal("Global config was not set")
	}

	if globalConfig.Plugins.Directory != "/test/plugins" {
		t.Errorf("Expected plugins directory to be '/test/plugins', got '%s'", globalConfig.Plugins.Directory)
	}

	if !globalConfig.Plugins.Enabled["test-plugin"] {
		t.Error("Expected test-plugin to be enabled")
	}
}

func TestGetGlobalTimerWithConfig(t *testing.T) {
	// Create a test config
	testConfig := &config.Config{}
	testConfig.Plugins.Directory = "/test/plugins"
	testConfig.Plugins.Enabled = map[string]bool{
		"test-plugin": true,
	}

	// Reset global variables for clean test
	globalTimer = nil
	timerOnce = sync.Once{}
	globalConfig = testConfig

	// Get the global timer
	timer := GetGlobalTimer()

	if timer == nil {
		t.Fatal("Expected timer to be created")
	}

	// Verify that the timer was created with the correct config
	if globalConfig == nil {
		t.Fatal("Global config should not be nil after timer creation")
	}

	if globalConfig.Plugins.Directory != "/test/plugins" {
		t.Errorf("Expected plugins directory to be '/test/plugins', got '%s'", globalConfig.Plugins.Directory)
	}
}
