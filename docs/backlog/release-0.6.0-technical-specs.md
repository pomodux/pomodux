# Release 0.6.0 Technical Specifications

> **STATUS: PLANNING**
> 
> **Document Type:** Technical Specifications
> **Release:** [Release 0.6.0](release-0.6.0.md)
> **Last Updated:** 2025-01-27

---

## Overview

This document provides detailed technical specifications for Release 0.6.0 - TUI-Only Timer with File-Based Locking. It defines the implementation approach, architecture changes, and technical requirements for each component.

## Current Codebase Analysis

### Existing Implementation
- **TUI**: Basic Bubbletea implementation exists in `internal/tui/tui.go` with Lipgloss styling
- **Plugin System**: Lua-based system with tview dialogs (needs migration per ADR 007)
- **Timer Core**: SessionType enum-based system with persistent state management
- **CLI Commands**: Complete CLI with `start`, `break`, `long-break`, `pause`, `resume`, `stop`, `status`
- **Configuration**: Comprehensive config system with XDG compliance
- **Logging**: Structured logging with logrus
- **Testing**: Unit tests, integration tests, and UAT with bats

### Key Issues Identified
1. **Mixed TUI Libraries**: Plugin system uses tview (violates ADR 007)
2. **SessionType Architecture**: Hardcoded enum system (needs replacement)
3. **Cross-Process Architecture**: Current system has cross-process synchronization
4. **No File Locking**: No existing lock manager implementation
5. **No Global Stage**: No existing stage management system
6. **No teatest Integration**: TUI tests not using teatest framework

### Documentation Requirements
- **ADR for Global Stage**: ✅ **COMPLETED** - ADR 008 created for global stage pattern
- **Plugin Development Guide**: ✅ **COMPLETED** - Updated for Release 0.6.0 API
- **Migration Guide**: ✅ **COMPLETED** - Comprehensive migration guide created
- **Configuration File Specifications**: ✅ **COMPLETED** - Comprehensive configuration documentation created

## Architecture Changes

### 1. File-Based Lock Manager

#### 1.1 Lock File Structure
```go
// internal/timer/lock.go
type TimerLockManager struct {
    lockFile     string
    lockFd       *os.File
    locked       bool
    timerPID     int
    sessionName  string
    mu           sync.Mutex
}

type LockFileState struct {
    PID         int       `json:"pid"`
    SessionName string    `json:"session_name"`
    StartTime   time.Time `json:"start_time"`
    Duration    int       `json:"duration_seconds"`
    Locked      time.Time `json:"locked_at"`
}

// Integration with existing timer manager
func (t *Timer) StartWithSessionName(duration time.Duration, sessionName string) error {
    // Acquire lock first
    if err := t.lockManager.AcquireLock(sessionName, duration); err != nil {
        return err
    }
    
    // Continue with existing timer logic
    return t.StartWithType(duration, SessionType(sessionName))
}
```

#### 1.2 Lock File Location (XDG Compliant)
```go
func getLockDir() (string, error) {
    runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
    if runtimeDir == "" {
        // Fallback to state directory
        stateDir, err := getStateDir()
        if err != nil {
            return "", err
        }
        return filepath.Join(stateDir, "runtime"), nil
    }
    return filepath.Join(runtimeDir, "pomodux"), nil
}

// Lock file: ~/.local/state/pomodux/runtime/timer.lock
// Or: /run/user/1000/pomodux/timer.lock (if XDG_RUNTIME_DIR exists)
```

#### 1.3 Process Validation & Recovery
```go
func (lm *TimerLockManager) validateProcess(pid int) bool {
    // Check if process exists and is actually pomodux
    process, err := os.FindProcess(pid)
    if err != nil {
        return false
    }
    
    // On Unix systems, check if process is still alive
    if err := process.Signal(syscall.Signal(0)); err != nil {
        return false
    }
    
    return lm.validateProcessName(pid)
}

func (lm *TimerLockManager) recoverOrphanedLock() error {
    state, err := lm.ReadLockState()
    if err != nil {
        return err
    }
    
    if !lm.validateProcess(state.PID) {
        return lm.forceReleaseLock()
    }
    
    return ErrTimerAlreadyRunning
}
```

### 2. Generic Timer Architecture

#### 2.1 Timer Core Changes
```go
// internal/timer/timer.go
type Timer struct {
    mu             sync.Mutex
    status         TimerStatus
    sessionName    string        // Changed from SessionType
    startTime      time.Time
    duration       time.Duration
    elapsed        time.Duration
    stateManager   *StateManager
    historyManager *HistoryManager
    pluginManager  *plugin.PluginManager
    lockManager    *TimerLockManager // New field
}

// New method signature
func (t *Timer) StartWithSessionName(duration time.Duration, sessionName string) error {
    // Implementation with lock acquisition
}

// Backward compatibility method
func (t *Timer) StartWithType(duration time.Duration, sessionType SessionType) error {
    // Convert SessionType to sessionName for backward compatibility
    sessionName := string(sessionType)
    return t.StartWithSessionName(duration, sessionName)
}
```

#### 2.2 State Management Updates
```go
// internal/timer/state.go
type State struct {
    Status      TimerStatus   `json:"status"`
    SessionName string        `json:"session_name"` // Changed from SessionType
    Duration    time.Duration `json:"duration"`
    StartTime   time.Time     `json:"start_time"`
    Elapsed     time.Duration `json:"elapsed"`
}
```

### 3. Global Stage Architecture

#### 3.1 Stage Manager Structure
```go
// internal/tui/stage.go
type GlobalStage struct {
    timerWindow     *TimerWindow
    pluginModals    []*PluginModal
    notifications   []*Notification
    statusPanel     *StatusPanel
    mu              sync.RWMutex
    eventBus        chan StageEvent
}

type StageEvent struct {
    Type      string
    Data      interface{}
    Timestamp time.Time
}

type TimerWindow struct {
    sessionName string
    duration    time.Duration
    remaining   time.Duration
    progress    float64
    paused      bool
    style       lipgloss.Style
}

type PluginModal struct {
    title       string
    content     string
    visible     bool
    style       lipgloss.Style
}

type Notification struct {
    message     string
    visible     bool
    expiresAt   time.Time
    style       lipgloss.Style
}
```

#### 3.2 Stage Management
```go
func (s *GlobalStage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case StageEvent:
        return s.handleStageEvent(msg)
    case tea.WindowSizeMsg:
        return s.handleResize(msg)
    case tea.KeyMsg:
        return s.handleKeyPress(msg)
    }
    return s, nil
}

func (s *GlobalStage) View() string {
    // Render timer window as primary
    // Overlay plugin modals when visible
    // Show notifications at top
    // Display status panel at bottom
}
```

### 4. Enhanced Event System

#### 4.1 Timer Events (6 Events)
```go
// internal/timer/events.go
const (
    EventTimerSetup     = "timer_setup"
    EventTimerStarted   = "timer_started"
    EventTimerPaused    = "timer_paused"
    EventTimerResumed   = "timer_resumed"
    EventTimerCompleted = "timer_completed"
    EventTimerStopped   = "timer_stopped"
)

type TimerEvent struct {
    Type        string                 `json:"event_type"`
    SessionName string                 `json:"session_name"`
    Timestamp   time.Time              `json:"timestamp"`
    Data        map[string]interface{} `json:"data"`
}
```

#### 4.2 Event Logging
```go
// Example: timer_started event
logger.Info("Timer started", map[string]interface{}{
    "event_type":    "timer_started",
    "session_name":  "deep work",
    "duration":      "45m0s",
    "start_time":    time.Now().Unix(),
    "process_id":    os.Getpid(),
})

// Example: timer_paused event
logger.Info("Timer paused", map[string]interface{}{
    "event_type":    "timer_paused",
    "session_name":  "deep work", 
    "elapsed":       "23m15s",
    "remaining":     "21m45s",
    "pause_time":    time.Now().Unix(),
})
```

### 5. Plugin API Redesign

#### 5.1 New Plugin API Structure
```go
// internal/plugin/api_v2.go
type PluginAPIv2 struct {
    stage       *GlobalStage
    eventBus    chan PluginEvent
    logger      *logger.Logger
}

type PluginEvent struct {
    Type      string                 `json:"type"`
    Data      map[string]interface{} `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
}

// Plugin window spawning rules
func (api *PluginAPIv2) CanShowModal(eventType string) bool {
    allowedEvents := []string{
        EventTimerSetup,
        EventTimerPaused,
        EventTimerCompleted,
        EventTimerStopped,
    }
    
    for _, allowed := range allowedEvents {
        if eventType == allowed {
            return true
        }
    }
    return false
}

// Migration from tview to Bubbletea
func (api *PluginAPIv2) ShowModal(title, content string) error {
    // Replace tview.NewModal() with Bubbletea modal component
    return api.stage.ShowModal(title, content)
}

func (api *PluginAPIv2) ShowListSelection(title string, options []string) (int, bool) {
    // Replace tview.NewList() with Bubbletea list component
    return api.stage.ShowListSelection(title, options)
}
```

#### 5.2 Plugin Information Display
```go
// Simple status updates (max 50 characters)
func (api *PluginAPIv2) UpdateStatus(message string) error {
    if len(message) > 50 {
        message = message[:47] + "..."
    }
    return api.stage.UpdateStatus(message)
}

// Auto-dismissing notifications (3-5 second timeout)
func (api *PluginAPIv2) ShowNotification(message string, duration time.Duration) error {
    if duration < 3*time.Second || duration > 5*time.Second {
        duration = 3 * time.Second
    }
    return api.stage.ShowNotification(message, duration)
}
```

## Implementation Phases

### Phase 1: Core Timer Simplification + Lock Manager

#### 1.1 Lock Manager Implementation
**Files to Create/Modify:**
- `internal/timer/lock.go` (new)
- `internal/timer/timer.go` (modify)
- `internal/timer/state.go` (modify)
- `internal/timer/lock_test.go` (new)

**Key Implementation Tasks:**
1. Implement `TimerLockManager` with file-based locking
2. Add process validation and orphaned lock recovery
3. Create XDG-compliant lock file location
4. Integrate lock acquisition into timer startup
5. Add comprehensive error handling and logging
6. **NEW**: Create mock file system for testing

**Testing Requirements:**
- Unit tests for file locking operations with mocked file system
- Integration tests for process validation
- End-to-end tests for concurrent timer prevention
- Cross-platform lock behavior validation
- **NEW**: Test lock manager integration with existing timer

#### 1.2 Timer Core Refactoring
**Files to Modify:**
- `internal/timer/timer.go`
- `internal/timer/state.go`
- `internal/cli/start.go`
- `internal/cli/break.go` (remove)
- `internal/cli/long_break.go` (remove)

**Key Implementation Tasks:**
1. Replace `SessionType` enum with `sessionName string`
2. Update timer events to use `session_name`
3. Remove break and long-break commands
4. Update `pomodux start` to accept optional session name
5. Update state persistence with session names

### Phase 2: TUI-Only Refactor + Global Stage

#### 2.1 Global Stage Architecture
**Files to Create/Modify:**
- `internal/tui/stage.go` (new)
- `internal/tui/timer_window.go` (new)
- `internal/tui/tui.go` (modify)
- `internal/tui/stage_test.go` (new)

**Key Implementation Tasks:**
1. Create global stage singleton
2. Implement timer window component
3. Add plugin modal support
4. Create notification system
5. Implement responsive layout with Lipgloss
6. **NEW**: Migrate existing TUI to use global stage
7. **NEW**: Add teatest framework for TUI testing

#### 2.2 TUI Integration
**Files to Modify:**
- `internal/cli/start.go`
- `internal/tui/tui.go`
- `cmd/pomodux/main.go`

**Key Implementation Tasks:**
1. Convert `pomodux start` to launch TUI immediately
2. Implement TUI-first, timer-init-second pattern
3. Add session name display in timer window
4. Integrate lock manager with TUI lifecycle

### Phase 3: Plugin API Redesign

#### 3.1 Plugin API Design
**Files to Create/Modify:**
- `internal/plugin/api_v2.go` (new)
- `internal/plugin/manager.go` (modify)
- `internal/plugin/compatibility.go` (new)
- `internal/plugin/bubbletea_components.go` (new)

**Key Implementation Tasks:**
1. Design new plugin API for TUI architecture
2. Create plugin window spawning rules
3. Implement simple information display system
4. Add backward compatibility layer
5. **NEW**: Migrate tview dialogs to Bubbletea components
6. **NEW**: Remove tview dependency from go.mod

#### 3.2 Plugin System Implementation
**Files to Create/Modify:**
- `internal/plugin/window_manager.go` (new)
- `internal/plugin/sdk.go` (new)
- `examples/plugin_integration_v2.go` (new)

**Key Implementation Tasks:**
1. Implement plugin window management
2. Add plugin information zones
3. Create real-time update system
4. Develop plugin SDK with Lipgloss helpers

## Technical Requirements

### Performance Requirements
- **Timer Accuracy**: Equivalent to current implementation
- **Lock Operations**: Complete within 1 second
- **UI Updates**: Responsive within 100ms
- **Memory Usage**: < 100MB total
- **Startup Time**: < 2 seconds

### Cross-Platform Compatibility
- **Linux**: Primary target with full XDG compliance
- **macOS**: Full compatibility with fallback paths
- **Windows**: Basic compatibility with simplified paths

### Error Handling
- **Lock Conflicts**: Clear error messages with actionable guidance
- **Process Recovery**: Automatic recovery from crashed processes
- **File System Issues**: Graceful degradation with clear error messages
- **Plugin Errors**: Non-blocking error handling with logging

### Logging Requirements
- **Structured Logging**: All operations logged with structured data
- **Event Logging**: All 6 timer events logged with metadata
- **Plugin Logging**: Plugin interactions and errors logged
- **Debug Logging**: Comprehensive debug information for troubleshooting

## Testing Strategy

### Unit Testing
- **Lock Manager**: Mock file system and process operations
- **Timer Core**: Session name integration and state management
- **Global Stage**: Component interactions and event handling
- **Plugin API**: API validation and error handling

### Integration Testing
- **Timer + Lock Manager**: End-to-end timer operations
- **TUI + Timer**: User interface and timer integration
- **Plugin + Stage**: Plugin window management and events
- **Cross-Component**: Full system integration

### End-to-End Testing
- **User Workflows**: Complete user scenarios
- **Error Scenarios**: Lock conflicts, process crashes, plugin errors
- **Cross-Platform**: Platform-specific behavior validation
- **Performance**: Load testing and performance validation

### Test Coverage Requirements
- **Overall Coverage**: 80% minimum
- **Critical Paths**: 95% minimum (lock manager, timer core)
- **Public APIs**: 100% coverage
- **Error Handling**: 100% coverage for error paths

## Migration Strategy

### Breaking Changes
- **Existing Plugins**: Will not work with new API (single user environment)
- **Configuration**: No configuration changes required
- **Session History**: Existing history will be cleared (breaking change)
- **CLI Commands**: `pomodux break` and `pomodux long-break` removed

### Test Plugin Development
- **Custom Configuration**: Create test plugins to validate functionality
- **Plugin Testing**: Develop test plugins for window spawning and events
- **Documentation**: Create plugin development guide for new API

### Documentation Updates
- **User Guide**: Updated with new session naming
- **Plugin Development**: New API documentation and examples
- **Migration Guide**: Step-by-step migration instructions
- **Troubleshooting**: Common issues and solutions

## Risk Mitigation

### Medium Risk Items
1. **Plugin API Redesign**
   - **Mitigation**: Create test plugins to validate functionality
   - **Fallback**: Single user environment simplifies testing

2. **File Locking Implementation**
   - **Mitigation**: Use proven patterns and extensive testing
   - **Fallback**: In-memory state when file system unavailable

3. **Global Stage Pattern**
   - **Mitigation**: Start with minimal implementation and iterate
   - **Fallback**: Simplified window management if complexity too high

### Performance Monitoring
- **Baseline Establishment**: Measure current performance before changes
- **Continuous Monitoring**: Monitor performance throughout development
- **Optimization Targets**: Define specific performance targets
- **Regression Detection**: Automated performance regression testing

## Success Metrics

### Functional Metrics
- **Single Timer Enforcement**: 100% success rate
- **Lock Recovery**: Automatic recovery in < 5 seconds
- **Session Support**: Support for any session name
- **Test Plugin Functionality**: 100% test plugin validation

### Performance Metrics
- **Timer Accuracy**: ±1 second over 1 hour
- **Startup Time**: < 2 seconds
- **Memory Usage**: < 100MB
- **UI Responsiveness**: < 100ms for user interactions

### Quality Metrics
- **Test Coverage**: 80% overall, 95% critical paths
- **Error Rate**: < 1% for normal operations
- **User Satisfaction**: Measured through user acceptance testing
- **Documentation Completeness**: 100% API documentation

---

**Note**: These technical specifications provide the foundation for implementing Release 0.6.0. All implementation must follow the TDD approach and meet the quality standards defined in the release process rules. 