# TUI-Only Timer with Simple Plugin Information 📋 PLANNED

> **Note**: This backlog item follows the 4-gate approval process. Issues can be created from this backlog using the GitHub issue templates in `.github/ISSUE_TEMPLATE/`.

## Feature Status
- **Status**: 🔄 PLANNING (Part of Release 0.6.0)
- **Priority**: High
- **Component**: User Interface / Timer Engine
- **Dependencies**: ADR 007 (TUI Standardization), ADR 002 (Persistent Timer Design), ADR 004 (Plugin System)
- **Release**: [Release 0.6.0](release-0.6.0.md)

## Feature Description

Refactor the current CLI+TUI dual interface to a unified TUI-only architecture with a global stage manager, simplified generic timer (session name + duration), simple plugin information display, and robust file-based state locking. This eliminates cross-process synchronization complexity while simplifying the timer core to be session-agnostic and ensuring only one timer instance can exist system-wide.

Key constraints:
- **Single Timer Instance**: Only 1 timer can exist and run at any given time - enforced via file-based locking
- **Process-Safe Locking**: File-based lock prevents multiple process instances with automatic recovery
- **Generic Timer**: Timer accepts duration + session name (string), no hardcoded session types
- **Default Session Name**: `pomodux start 25m` creates session named "work" by default
- **Bubbletea-Native**: All implementations must work within Bubbletea's capabilities
- **Lipgloss Styling**: All UI components use Lipgloss for consistent theming and positioning
- **Single Process**: Eliminate cross-process synchronization issues
- **Event-Driven**: Single timer emits 6 events, everything else reacts
- **Window Spawning Rules**: Complex plugin windows only allowed during specific timer events

## User Story

As a user, I want a unified TUI interface with a simplified timer that accepts any session name so that I can time any type of activity without being restricted to predefined session types, and I want assurance that only one timer can run at any time across all processes.

## Enhanced User Stories

- **As a user, I want `pomodux start 25m` to create a "work" session by default** so that I get sensible defaults without specifying session names.
- **As a user, I want `pomodux start 15m "break"` to create a custom-named session** so that I can organize my timing sessions however I want.
- **As a user, I want the TUI to display the session name prominently** so that I can see what type of session I'm currently timing.
- **As a user, I want all sessions to have the same controls** so that pause, resume, and stop work identically regardless of session name.
- **As a user, I want clear error messages when a timer is already running** so that I know exactly what to do next.
- **As a user, I want automatic recovery from crashed timer processes** so that I never get stuck with orphaned locks.
- **As a plugin developer, I want to display simple information during active timing** so that I can provide contextual data without interrupting user focus.
- **As a plugin developer, I want to create complex dialogs during appropriate timer states** so that I can provide rich interactions when users are mentally available.
- **As a user, I want plugins to respect my focus time** so that I'm only interrupted by complex dialogs when I'm not actively timing.
- **As a developer, I want a single global stage for window management** so that all UI components work together seamlessly.

## Acceptance Criteria

### Single Timer Instance Architecture with File-Based Locking
- [ ] Only one timer instance can exist and run at any given time
- [ ] Timer instance serves as single source of truth for all timer state
- [ ] **File-based lock prevents multiple process instances**
- [ ] **Lock file contains current timer state (PID, session, duration)**
- [ ] **Atomic lock acquisition with process validation**
- [ ] **Automatic lock cleanup on process exit or crash**
- [ ] **Lock file validation checks for stale processes**
- [ ] **Helpful error handling**: Second timer attempt shows current status + actionable guidance
- [ ] **Process-aware messages**: "Timer running in process 1234 (22m 30s remaining in 'work' session)"
- [ ] **Clear next steps**: Suggest `pomodux status`, `pomodux stop`, or `pomodux pause` commands
- [ ] **Lock file recovery**: Handle orphaned locks from crashed processes
- [ ] **Comprehensive logging**: All timer operations, lock operations, and state changes logged
- [ ] **Structured log data**: Session names, PIDs, timestamps, and operation results in log entries
- [ ] **Error logging**: All error conditions logged with sufficient context for debugging
- [ ] Timer state validation and consistency checks
- [ ] `pomodux start [duration]` launches TUI immediately (no CLI process)
- [ ] TUI initializes timer after launch (TUI-first, then timer init)
- [ ] Single process architecture eliminates cross-process synchronization
- [ ] All existing timer functionality preserved (pause, resume, stop, progress)
- [ ] Timer completion and exit behavior matches current experience
- [ ] Terminal state properly restored on exit

### Lock Manager Implementation
- [ ] **XDG-compliant lock file location** (`$XDG_RUNTIME_DIR/pomodux/timer.lock` or fallback)
- [ ] **JSON-based lock file format** with timer state information
- [ ] **Process validation** to detect orphaned locks from crashed processes
- [ ] **Atomic lock operations** using file system primitives
- [ ] **Lock timeout handling** for edge cases and system recovery
- [ ] **Graceful error messages** for all lock conflict scenarios
- [ ] **Cross-platform compatibility** (Linux, macOS, Windows)

### Generic Timer Architecture
- [ ] Timer accepts `duration` + `sessionName` (string) instead of hardcoded SessionType enum
- [ ] Remove SessionTypeWork, SessionTypeBreak, SessionTypeLongBreak constants
- [ ] Default session name "work" when no name specified
- [ ] Support custom session names: `pomodux start 30m "deep work"`
- [ ] Session name included in all timer events and history records
- [ ] Backward compatibility: existing history shows session names instead of types
- [ ] **Global Stage Architecture**: Singleton stage manages all UI components in defined hierarchy
- [ ] **Timer Window (Primary)**: Always centered, main display, never hidden
- [ ] **Plugin Modal Window**: Overlays timer during allowed events only
- [ ] **Plugin Notification**: Auto-dismissing overlays allowed during all events
- [ ] **Plugin Status Panel**: Integrated zone within timer window for real-time metrics
- [ ] **Single Process Flow**: TUI process → Timer instance → Events → Stage → Components
- [ ] **Event Distribution**: Stage coordinates all timer events to appropriate components
- [ ] **Z-order Management**: Proper layering of timer, modals, and notifications
- [ ] **Consistent Lipgloss theming and styling** across all stage components (ADR 007)
- [ ] **Unified color palette and style definitions** using Lipgloss for all UI elements
- [ ] **Responsive layout system** using Lipgloss positioning and alignment

### Timer Event System (6 Events)
- [ ] `timer_setup` - Before timer starts (blocking, can modify/cancel)
- [ ] `timer_started` - After timer begins (non-blocking)
- [ ] `timer_paused` - When user pauses (non-blocking)
- [ ] `timer_resumed` - When user resumes (non-blocking)
- [ ] `timer_completed` - Natural completion (non-blocking)
- [ ] `timer_stopped` - Manual stop (non-blocking)
- [ ] Timer events include `session_name` instead of `session_type`
- [ ] All events properly trigger plugin hooks and TUI updates with session context
- [ ] **Comprehensive event logging**: All 6 timer events logged through standard logger with structured data
- [ ] **Event metadata logging**: Session name, duration, timestamps, and event-specific data logged
- [ ] **Plugin interaction logging**: Plugin responses and errors logged for debugging

### Plugin Window Spawning Rules
- [ ] **Modal windows allowed**: `timer_setup`, `timer_paused`, `timer_completed`, `timer_stopped`
- [ ] **Modal windows blocked**: `timer_started`, `timer_resumed`
- [ ] **Simple info allowed**: All 6 events can display non-modal information
- [ ] **User cancellation handling**: Plugin dialogs can be canceled by user without error
- [ ] **Graceful cancellation**: User canceling timer setup returns to shell, not error state
- [ ] API validation prevents complex UI during active timing periods

### Simple Plugin Information Display
- [ ] Status line updates (max 50 characters, single line)
- [ ] Auto-dismissing notifications (3-5 second timeout)
- [ ] Timer title enhancements (session context)
- [ ] Progress metrics display (key-value pairs)
- [ ] Real-time updates via Bubbletea message system
- [ ] Information zones integrated into timer panel layout

### Enhanced Error Handling & User Experience
- [ ] **Timer conflict detection**: Clear identification of existing timer process
- [ ] **Process information display**: Show PID, session name, remaining time
- [ ] **Actionable guidance**: Suggest specific commands (status, stop, pause)
- [ ] **Automatic recovery**: Handle crashed processes and orphaned locks
- [ ] **Graceful degradation**: Fallback behavior when lock system unavailable
- [ ] **Comprehensive logging**: All lock operations logged for troubleshooting

### ADR Compliance
- [ ] **ADR 002**: Maintains persistent timer with keypress controls and real-time display
- [ ] **ADR 004**: Plugin system integration preserved and enhanced with window management
- [ ] **ADR 007**: All UI components use Bubbletea with unified Lipgloss theming
- [ ] Single process eliminates previous cross-process synchronization problems

### Non-functional Requirements
- [ ] Performance equivalent to current timer (no measurable impact on accuracy)
- [ ] Cross-platform compatibility (Linux, macOS, Windows)
- [ ] Backward compatibility with existing plugin API where possible
- [ ] Memory usage within reasonable bounds (< 100MB)
- [ ] **Timer window centering**: Window always centered horizontally and vertically in terminal
- [ ] **Dynamic repositioning**: Real-time centering adaptation on terminal resize
- [ ] **Responsive sizing**: Timer window scales appropriately to terminal dimensions
- [ ] **Minimum size handling**: Graceful degradation for small terminals (40x12 minimum)
- [ ] **Maximum size constraints**: Timer window doesn't become unnecessarily large
- [ ] **Lipgloss positioning**: Use lipgloss.Place() for precise center alignment
- [ ] Responsive UI updates (< 100ms for user interactions)
- [ ] Lock file operations complete within 1 second
- [ ] Graceful handling of file system permission issues
- [ ] Lock file operations complete within 1 second
- [ ] Graceful handling of file system permission issues

## Timer Window UX Design

### ASCII Representation

#### Basic Timer Window (Core Experience)
```
┌─────────────────────────────────────────────────────────────────────┐
│                              WORK SESSION                           │
│                                                                     │
│                            12m 30s remaining                        │
│                                                                     │
│    [████████████████████████████████░░░░░░░░░░░░░░]  65%            │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  [P]ause  [R]esume  [S]top  [Q]uit                    Ctrl+C Exit   │
└─────────────────────────────────────────────────────────────────────┘
```

#### Custom Session Example
```
┌─────────────────────────────────────────────────────────────────────┐
│                           DEEP FOCUS SESSION                        │
│                                                                     │
│                            45m 15s remaining                        │
│                                                                     │
│    [██████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  25%            │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  [P]ause  [R]esume  [S]top  [Q]uit                    Ctrl+C Exit   │
└─────────────────────────────────────────────────────────────────────┘
```

#### Paused State
```
┌─────────────────────────────────────────────────────────────────────┐
│                        BREAK SESSION (PAUSED)                       │
│                                                                     │
│                            3m 45s remaining                         │
│                                                                     │
│    [████████████████████████████████████████░░░░░░░░]  80%         │
│                                                                     │
│                    ⏸️  Timer paused. Press [R] to resume.           │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  [P]ause  [R]esume  [S]top  [Q]uit                    Ctrl+C Exit   │
└─────────────────────────────────────────────────────────────────────┘
```

### Plugin Information Display (Separate Examples)

#### Simple Status Integration
```
┌─────────────────────────────────────────────────────────────────────┐
│                              WORK SESSION                           │
│                                                                     │
│                            12m 30s remaining                        │
│                                                                     │
│    [████████████████████████████████░░░░░░░░░░░░░░]  65%            │
│                                                                     │
│                     📊 Focus Score: 95% | 🔥 Streak: 4             │ ← Plugin status
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  [P]ause  [R]esume  [S]top  [Q]uit                    Ctrl+C Exit   │
└─────────────────────────────────────────────────────────────────────┘
```

#### Auto-Dismissing Notification
```
┌─────────────────────────────────────────────────────────────────────┐
│                           DEEP FOCUS SESSION                        │
│                                                                     │
│                            45m 15s remaining                        │
│                                                                     │
│    [██████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  25%            │
│                                                                     │
│                        ✅ Logged to Kimai (2m ago)                  │ ← Auto-dismiss
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  [P]ause  [R]esume  [S]top  [Q]uit                    Ctrl+C Exit   │
└─────────────────────────────────────────────────────────────────────┘
```

## Proposed Technical Design
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

// Example: plugin interaction
logger.Debug("Plugin hook executed", map[string]interface{}{
    "event_type":    "timer_completed",
    "plugin_name":   "kimai_integration",
    "hook_result":   "success",
    "execution_time": "150ms",
    "session_name":  "work",
})

// Example: plugin error
logger.Error("Plugin hook failed", map[string]interface{}{
    "event_type":    "timer_setup", 
    "plugin_name":   "todoist_sync",
    "error":         "API authentication failed",
    "session_name":  "work",
    "retry_count":   2,
})
```
- **Custom Messages**: Define TimerEvent, PluginInfo, WindowManagement message types
- **Command Coordination**: Use tea.Batch for multiple commands, tea.Sequence for ordered operations
- **Component Models**: Each component implements tea.Model interface with Lipgloss styling
- **State Management**: Centralized state with component-specific sub-models
- **Lipgloss Rendering**: All View() methods use Lipgloss styles for consistent presentation
- **Dynamic Styling**: Lipgloss styles that adapt based on timer state and terminal dimensions
- **Z-order Rendering**: Stage manages layering of timer window, modals, and notifications

### File-Based Lock Manager Architecture

#### Lock File Structure
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
```

#### Enhanced Error Types
```go
type TimerConflictError struct {
    PID         int
    SessionName string
    StartTime   time.Time
    Duration    time.Duration
}

func (e *TimerConflictError) Error() string {
    elapsed := time.Since(e.StartTime)
    remaining := e.Duration - elapsed
    
    if remaining <= 0 {
        return fmt.Sprintf("Timer process %d may have completed. Try: pomodux status", e.PID)
    }
    
    return fmt.Sprintf(
        "Timer already running in process %d (%s remaining in '%s' session)\n" +
        "Try: pomodux status | pomodux stop | pomodux pause",
        e.PID,
        formatDuration(remaining),
        e.SessionName,
    )
}
```

#### Lock File Location (XDG Compliant)
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

### Enhanced Timer Startup Sequence with Comprehensive Logging
```go
func (t *Timer) StartWithSessionName(duration time.Duration, sessionName string) error {
    logger.Info("Timer startup initiated", map[string]interface{}{
        "session_name": sessionName,
        "duration":     duration.String(),
        "process_id":   os.Getpid(),
    })
    
    // 1. Attempt to acquire timer lock
    lockManager, err := NewTimerLockManager()
    if err != nil {
        logger.Error("Failed to initialize lock manager", map[string]interface{}{
            "error": err.Error(),
            "session_name": sessionName,
        })
        return fmt.Errorf("failed to initialize lock manager: %w", err)
    }

    // 2. Try to acquire exclusive lock
    logger.Debug("Attempting to acquire timer lock", map[string]interface{}{
        "session_name": sessionName,
        "duration": duration.String(),
    })
    
    if err := lockManager.AcquireLock(sessionName, duration); err != nil {
        if errors.Is(err, ErrTimerAlreadyRunning) {
            // Read current state from lock file
            state, _ := lockManager.ReadLockState()
            logger.Warn("Timer already running in another process", map[string]interface{}{
                "running_pid":         state.PID,
                "running_session":     state.SessionName,
                "running_start_time":  state.StartTime,
                "attempted_session":   sessionName,
                "attempted_duration":  duration.String(),
            })
            return &TimerConflictError{
                PID:         state.PID,
                SessionName: state.SessionName,
                StartTime:   state.StartTime,
                Duration:    time.Duration(state.Duration) * time.Second,
            }
        }
        logger.Error("Failed to acquire timer lock", map[string]interface{}{
            "error": err.Error(),
            "session_name": sessionName,
        })
        return fmt.Errorf("failed to acquire timer lock: %w", err)
    }

    logger.Info("Timer lock acquired successfully", map[string]interface{}{
        "session_name": sessionName,
        "lock_file": lockManager.lockFile,
        "process_id": os.Getpid(),
    })

    // 3. Set up cleanup on exit
    t.lockManager = lockManager
    defer func() {
        logger.Debug("Releasing timer lock on exit", map[string]interface{}{
            "session_name": sessionName,
        })
        t.lockManager.ReleaseLock()
    }()

    logger.Info("Timer startup complete", map[string]interface{}{
        "session_name": sessionName,
        "duration": duration.String(),
        "status": "ready_to_start",
    })

    // 4. Continue with timer startup...
}
```

### Process Validation & Recovery with Logging
```go
func (lm *TimerLockManager) validateProcess(pid int) bool {
    logger.Debug("Validating process for lock ownership", map[string]interface{}{
        "pid": pid,
        "validator_pid": os.Getpid(),
    })
    
    // Check if process exists and is actually pomodux
    process, err := os.FindProcess(pid)
    if err != nil {
        logger.Debug("Process not found during validation", map[string]interface{}{
            "pid": pid,
            "error": err.Error(),
        })
        return false
    }
    
    // On Unix systems, check if process is still alive
    if err := process.Signal(syscall.Signal(0)); err != nil {
        logger.Debug("Process not responding to signal", map[string]interface{}{
            "pid": pid,
            "error": err.Error(),
        })
        return false
    }
    
    // Additional validation: check process name/cmdline if needed
    isValid := lm.validateProcessName(pid)
    logger.Debug("Process validation complete", map[string]interface{}{
        "pid": pid,
        "valid": isValid,
    })
    return isValid
}

func (lm *TimerLockManager) recoverOrphanedLock() error {
    logger.Info("Attempting to recover orphaned lock", map[string]interface{}{
        "lock_file": lm.lockFile,
    })
    
    state, err := lm.ReadLockState()
    if err != nil {
        logger.Error("Failed to read lock state during recovery", map[string]interface{}{
            "lock_file": lm.lockFile,
            "error": err.Error(),
        })
        return err
    }
    
    if !lm.validateProcess(state.PID) {
        logger.Warn("Recovering orphaned timer lock", map[string]interface{}{
            "orphaned_pid": state.PID,
            "session_name": state.SessionName,
            "start_time": state.StartTime,
            "lock_age": time.Since(state.Locked),
        })
        
        if err := lm.forceReleaseLock(); err != nil {
            logger.Error("Failed to force release orphaned lock", map[string]interface{}{
                "orphaned_pid": state.PID,
                "error": err.Error(),
            })
            return err
        }
        
        logger.Info("Orphaned lock recovered successfully", map[string]interface{}{
            "orphaned_pid": state.PID,
            "session_name": state.SessionName,
        })
        return nil
    }
    
    logger.Debug("Lock owner process is still valid", map[string]interface{}{
        "owner_pid": state.PID,
        "session_name": state.SessionName,
    })
    return ErrTimerAlreadyRunning
}
```

### High Risk
- **Lock File Implementation Complexity**: File-based locking across platforms adds significant complexity
  - **Mitigation**: Use proven file locking patterns with comprehensive testing
  - **Mitigation**: Implement atomic operations with proper error handling
  - **Mitigation**: Add extensive cross-platform testing in CI/CD
- **Process Validation Edge Cases**: Detecting orphaned locks reliably across platforms
  - **Mitigation**: Conservative validation with timeout fallbacks
  - **Mitigation**: Manual recovery commands for edge cases
- **Architecture Complexity**: Global stage pattern adds significant complexity
  - **Mitigation**: Start with minimal viable stage, iterate incrementally
  - **Mitigation**: Comprehensive testing with teatest framework
- **Plugin API Redesign**: Complete restructuring of plugin system for new architecture
  - **Mitigation**: Clean slate approach - design optimal API without legacy constraints
  - **Mitigation**: Comprehensive documentation and examples for new plugin patterns

### Medium Risk
- **Performance Impact**: Single process handling all UI and timer logic
  - **Mitigation**: Profile performance vs current implementation
  - **Mitigation**: Optimize hot paths and minimize unnecessary updates
- **Event System Complexity**: 6 events with different behaviors and user interaction patterns
  - **Mitigation**: Clear documentation and examples for each event type
  - **Mitigation**: Runtime validation and helpful error messages for actual errors
  - **Mitigation**: Distinguish user cancellations from system errors
- **File System Permissions**: Lock file creation may fail in restricted environments
  - **Mitigation**: Graceful degradation with clear error messages
  - **Mitigation**: Fallback to in-memory state when file system unavailable

### Low Risk
- **Cross-Platform Issues**: Bubbletea handles platform differences
  - **Mitigation**: Test on all target platforms during development
- **User Learning Curve**: TUI-only may confuse CLI-only users
  - **Mitigation**: Clear documentation and smooth transition

## Implementation Notes

### Phase 1: Core Timer Simplification + Lock Manager + Logging
1. Replace SessionType enum with `sessionName string` in timer core
2. **Implement TimerLockManager with file-based locking and comprehensive logging**
3. **Add process validation and orphaned lock recovery with detailed logging**
4. Update timer events to use `session_name` instead of `session_type`
5. Remove `pomodux break` and `pomodux long-break` commands  
6. Update `pomodux start` to accept optional session name parameter
7. **Integrate lock acquisition into timer startup sequence with full logging**
8. **Add comprehensive error messages for timer conflicts**
9. **Implement structured logging for all timer operations using standard logger**

### Phase 1.5: Lock Manager Testing & Validation
1. **Test concurrent timer prevention across processes**
2. **Test orphaned lock recovery scenarios**
3. **Test lock file corruption handling**
4. **Test process crash scenarios**
5. **Cross-platform lock behavior validation**
6. **Performance testing of lock operations**

### Phase 2: TUI-Only Refactor + Global Stage + Event Logging
1. Create global stage singleton with basic window management
2. Convert `pomodux start` to launch TUI immediately with session name display
3. Implement TUI-first, timer-init-second pattern
4. **Integrate lock manager with TUI lifecycle and logging**
5. **Implement comprehensive event logging for all 6 timer events**
6. **Add structured logging for stage and component interactions**
7. Test across all platforms with custom session names

### Phase 3: Event System & Plugin API Redesign + Plugin Logging
1. Implement all 6 timer events with proper timing and logging
2. Create event distribution system from timer to stage with comprehensive logging
3. **Design new plugin API from scratch** - optimized for TUI architecture
4. **Implement plugin window spawning rules and validation with logging**
5. **Create plugin development examples and documentation**
6. **Add structured logging for all plugin interactions, successes, and failures**
7. Test new plugin system with comprehensive scenarios

### Phase 4: Enhanced Plugin Information System
1. Design and implement information zones in timer panel
2. **Create new plugin API for simple information updates**
3. Add real-time update system via Bubbletea messages
4. Implement auto-dismissing notifications and status updates
5. **Develop plugin SDK with Lipgloss styling helpers**

### Phase 5: Advanced Plugin Windows & API Finalization
1. **Implement complex plugin window support during allowed events**
2. **Create comprehensive plugin API for window creation and management**
3. **Add Lipgloss-based plugin styling system**
4. **Develop plugin development tools and hot-reload capabilities**
5. Comprehensive testing of all plugin interaction patterns
6. **Create plugin marketplace/discovery system foundation**

### Success Criteria
- **Single Timer Enforcement**: Only one timer instance can exist at any time (process-safe)
- **Robust Lock Management**: Automatic recovery from crashed processes and orphaned locks
- **Generic Session Support**: Timer accepts any string as session name
- **Simplified Command Interface**: Single `pomodux start` command replaces multiple session-specific commands
- **Zero cross-process synchronization**: Single process architecture with file-based state locking
- **Maintained functionality**: All existing timer features work identically
- **Enhanced plugin capabilities**: Modern plugin API designed specifically for TUI architecture
- **Plugin development experience**: Comprehensive SDK with Lipgloss styling helpers and development tools
- **Improved user experience**: Immediate visual feedback with session context and clear error messages
- **ADR compliance**: Aligns with all existing architectural decisions

### Key Implementation Decisions
1. **File-Based Lock Manager**: Prevents multiple timer processes with automatic recovery
2. **Process Validation**: Detect and recover from crashed timer processes
3. **Global Stage Pattern**: Singleton stage manages all TUI components
4. **Event-Driven Architecture**: Timer emits events, everything else reacts
5. **Bubbletea-Native**: No workarounds, work within framework capabilities
6. **Plugin API Redesign**: Clean slate approach for optimal TUI integration
7. **Simple Information System**: Non-modal updates during active timing
8. **XDG Compliance**: Lock files follow proper system standards
9. **Plugin Development SDK**: Comprehensive tools and helpers for plugin creators

This feature represents a significant architectural improvement that simplifies the codebase while expanding capabilities and ensuring robust single-timer enforcement across processes, aligning with existing ADRs and maintaining backward compatibility where possible.