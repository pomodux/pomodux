package timer

import (
	"sync"

	"github.com/rsmacapinlac/pomodux/internal/config"
	"github.com/rsmacapinlac/pomodux/internal/logger"
	"github.com/rsmacapinlac/pomodux/internal/plugin"
)

// Global timer instance
var (
	globalTimer          *Timer
	timerOnce            sync.Once
	globalStateManager   *StateManager
	globalHistoryManager *HistoryManager
	globalPluginManager  *plugin.PluginManager
)

// GetGlobalTimer returns the global timer instance.
// This ensures all CLI commands use the same timer.
func GetGlobalTimer() *Timer {
	timerOnce.Do(func() {
		// Create state manager
		stateManager, err := NewStateManager()
		if err != nil {
			// If we can't create state manager, create timer without persistence
			globalTimer = NewTimer()
			return
		}
		globalStateManager = stateManager

		// Create history manager
		historyManager, err := NewHistoryManager()
		if err != nil {
			// If we can't create history manager, create timer without history
			globalTimer = NewTimerWithManagers(stateManager, nil)
			return
		}
		globalHistoryManager = historyManager

		// Initialize plugin manager
		cfg, err := config.Load()
		if err != nil {
			logger.Warn("Failed to load configuration for plugin system", map[string]interface{}{"error": err.Error()})
			// Continue without plugin system
			globalTimer = NewTimerWithManagers(stateManager, historyManager)
			return
		}

		// Create plugin manager
		pluginManager := plugin.NewPluginManager(cfg.Plugins.Directory, cfg)

		// Load plugins
		if err := pluginManager.LoadPlugins(); err != nil {
			logger.Warn("Failed to load plugins", map[string]interface{}{"error": err.Error()})
			// Continue without plugins
			globalTimer = NewTimerWithManagers(stateManager, historyManager)
			return
		}

		globalPluginManager = pluginManager

		// List loaded plugins
		plugins := pluginManager.ListPlugins()
		logger.Info("Plugin system initialized", map[string]interface{}{
			"plugins_directory": cfg.Plugins.Directory,
			"plugins_loaded":    len(plugins),
		})

		for _, p := range plugins {
			logger.Info("Plugin loaded", map[string]interface{}{
				"name":        p.Name,
				"version":     p.Version,
				"description": p.Description,
				"author":      p.Author,
			})
		}

		// Create timer with all managers
		globalTimer = NewTimerWithManagers(stateManager, historyManager)
		globalTimer.SetPluginManager(pluginManager)
	})

	return globalTimer
}

// ShutdownGlobalTimer gracefully shuts down the global timer.
// This should be called when the application exits.
func ShutdownGlobalTimer() {
	if globalPluginManager != nil {
		globalPluginManager.Shutdown()
	}
}
