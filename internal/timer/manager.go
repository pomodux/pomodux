package timer

import (
	"sync"

	"github.com/pomodux/pomodux/internal/config"
	"github.com/pomodux/pomodux/internal/logger"
)

// Global timer instance - now using EventDrivenTimer and Application
var (
	globalApplication *Application
	globalTimer       *EventDrivenTimer
	globalConfig      *config.Config
	timerOnce         sync.Once
)

// GetGlobalApplication returns the global application instance.
// This ensures single-process timer coordination.
func GetGlobalApplication() *Application {
	timerOnce.Do(func() {
		var err error
		globalApplication, err = NewApplication()
		if err != nil {
			logger.Error("Failed to create global application", err, map[string]interface{}{})
			// Create minimal application without full initialization
			globalApplication = &Application{}
		}
		globalTimer = globalApplication.GetTimer()
	})

	return globalApplication
}

// GetGlobalTimer returns the global event-driven timer instance.
func GetGlobalTimer() *EventDrivenTimer {
	app := GetGlobalApplication()
	return app.GetTimer()
}

// SetGlobalConfig sets the global configuration to be used by the timer manager.
// This must be called before GetGlobalTimer() to ensure the correct config is used.
func SetGlobalConfig(cfg *config.Config) {
	globalConfig = cfg
}

// ShutdownGlobalTimer gracefully shuts down the global application.
// This should be called when the application exits.
func ShutdownGlobalTimer() {
	if globalApplication != nil {
		globalApplication.Close()
	}
}
