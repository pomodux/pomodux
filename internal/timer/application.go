package timer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pomodux/pomodux/internal/config"
	"github.com/pomodux/pomodux/internal/logger"
	"github.com/pomodux/pomodux/internal/plugin"
)

// Application represents the unified single-process Pomodux application
type Application struct {
	// Core Components
	timer       *EventDrivenTimer
	coordinator *SingleTimerCoordinator

	// Support Components
	config         *config.Config
	pluginManager  *plugin.PluginManager
	historyManager *HistoryManager
	stateManager   *StateManager

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex
}

// NewApplication creates a new unified Pomodux application
func NewApplication() (*Application, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Create single timer coordinator
	coordinator, err := NewSingleTimerCoordinator()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create timer coordinator: %w", err)
	}

	// Initialize managers
	stateManager, err := NewStateManager()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create state manager: %w", err)
	}

	historyManager, err := NewHistoryManager()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create history manager: %w", err)
	}

	// Create event-driven timer
	timer := NewEventDrivenTimerWithManagers(stateManager, historyManager)

	// Create plugin manager
	pluginManager := plugin.NewPluginManager(cfg.Plugins.Directory, cfg)

	app := &Application{
		timer:          timer,
		coordinator:    coordinator,
		config:         cfg,
		pluginManager:  pluginManager,
		historyManager: historyManager,
		stateManager:   stateManager,
		ctx:            ctx,
		cancel:         cancel,
	}

	// Initialize the application
	if err := app.initialize(); err != nil {
		app.Close()
		return nil, fmt.Errorf("failed to initialize application: %w", err)
	}

	logger.Info("Pomodux application created successfully", map[string]interface{}{
		"process_id": os.Getpid(),
		"plugins_enabled": cfg.Plugins.Enabled,
	})

	return app, nil
}

// initialize sets up all application components
func (app *Application) initialize() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	// Load plugins (plugin manager will check individual plugin enabled status)
	if err := app.pluginManager.LoadPlugins(); err != nil {
		logger.Warn("Failed to load plugins", map[string]interface{}{
			"error": err.Error(),
		})
		// Don't fail application startup for plugin errors
	}

	// Wire timer to managers
	// Note: EventDrivenTimer embeds Timer, so we can access its methods
	app.timer.SetHistoryManager(app.historyManager)
	app.timer.SetPluginManager(app.pluginManager)

	// Setup graceful shutdown
	app.setupShutdownHandlers()

	logger.Debug("Application initialized successfully")
	return nil
}

// StartTimer starts a timer session with the given duration and session name
func (app *Application) StartTimer(duration time.Duration, sessionName string) error {
	app.mu.Lock()
	defer app.mu.Unlock()

	// Acquire timer lock
	if err := app.coordinator.AcquireTimerLock(app.timer); err != nil {
		return fmt.Errorf("failed to acquire timer lock: %w", err)
	}

	// Start the timer
	if err := app.timer.StartWithSessionName(duration, sessionName); err != nil {
		// Release lock on failure
		app.coordinator.ReleaseLock()
		return fmt.Errorf("failed to start timer: %w", err)
	}

	logger.Info("Timer started successfully", map[string]interface{}{
		"duration":     duration,
		"session_name": sessionName,
		"process_id":   os.Getpid(),
	})

	return nil
}

// StopTimer stops the currently running timer
func (app *Application) StopTimer() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	if !app.coordinator.IsTimerRunning() {
		return fmt.Errorf("no timer is currently running")
	}

	// Stop the timer
	if err := app.timer.Stop(); err != nil {
		return fmt.Errorf("failed to stop timer: %w", err)
	}

	// Release the lock
	app.coordinator.ReleaseLock()

	logger.Info("Timer stopped successfully", map[string]interface{}{
		"process_id": os.Getpid(),
	})

	return nil
}

// PauseTimer pauses the currently running timer
func (app *Application) PauseTimer() error {
	app.mu.RLock()
	defer app.mu.RUnlock()

	if !app.coordinator.IsTimerRunning() {
		return fmt.Errorf("no timer is currently running")
	}

	return app.timer.Pause()
}

// ResumeTimer resumes a paused timer
func (app *Application) ResumeTimer() error {
	app.mu.RLock()
	defer app.mu.RUnlock()

	if !app.coordinator.IsTimerRunning() {
		return fmt.Errorf("no timer is currently running")
	}

	return app.timer.Resume()
}

// GetTimer returns the application's timer instance
func (app *Application) GetTimer() *EventDrivenTimer {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.timer
}

// GetConfig returns the application's configuration
func (app *Application) GetConfig() *config.Config {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.config
}

// IsTimerRunning returns true if a timer is currently running
func (app *Application) IsTimerRunning() bool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.coordinator.IsTimerRunning()
}

// setupShutdownHandlers sets up graceful shutdown for the application
func (app *Application) setupShutdownHandlers() {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-sigChan:
			logger.Info("Shutdown signal received", map[string]interface{}{
				"signal": sig.String(),
			})
			app.shutdown()
		case <-app.ctx.Done():
			logger.Debug("Application context canceled")
		}
	}()
}

// shutdown performs graceful shutdown of the application
func (app *Application) shutdown() {
	logger.Info("Starting graceful shutdown")

	// Stop any running timer
	if app.coordinator.IsTimerRunning() {
		logger.Info("Stopping timer during shutdown")
		if err := app.StopTimer(); err != nil {
			logger.Warn("Failed to stop timer during shutdown", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Cancel context to stop background operations
	app.cancel()

	logger.Info("Graceful shutdown completed")
}

// Close cleans up all application resources
func (app *Application) Close() {
	logger.Debug("Closing application")

	// Shutdown gracefully
	app.shutdown()

	// Wait for background goroutines to finish
	app.wg.Wait()

	// Close timer resources
	if app.timer != nil {
		app.timer.Close()
	}

	// Release any remaining locks
	if app.coordinator != nil {
		app.coordinator.ReleaseLock()
	}

	logger.Info("Application closed successfully")
}

// Wait blocks until the application context is canceled
func (app *Application) Wait() {
	<-app.ctx.Done()
	app.wg.Wait()
}
