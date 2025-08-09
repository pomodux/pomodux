---
status: approved
type: technical
---

# ADR 011: TUI-First Command Architecture

## 1. Context / Background

The original Pomodux command architecture used a dual approach where CLI commands (`start`, `pause`, `resume`, `stop`) operated on a background timer process, and the TUI was launched separately for visual feedback. This approach had several issues:

- **Cross-Process Complexity**: Synchronization between CLI and TUI processes
- **State Consistency**: Potential race conditions between multiple command invocations
- **User Experience**: Users had to choose between CLI or TUI modes
- **Process Management**: Complex coordination between timer process and display process

The TUI timer feature implementation requires a unified approach where the timer runs within the TUI process, eliminating cross-process synchronization while providing immediate visual feedback.

## 2. Decision

**Pomodux will implement a TUI-first architecture where `pomodux start` immediately launches the TUI with the timer running inside the same process, eliminating CLI+TUI dual modes.**

### Architecture Changes

- **Single Process**: Timer and TUI run in the same process
- **Immediate TUI Launch**: `pomodux start` launches TUI immediately (blocking)
- **Unified State**: No cross-process state synchronization needed
- **Simplified Commands**: Remove separate CLI timer process commands
- **Lock Integration**: File-based locking prevents multiple timer instances

### Command Behavior Changes
```bash
# New behavior - immediate TUI launch
pomodux start 25m "work"        # Launches TUI immediately
pomodux start 5m "break"        # Launches TUI immediately  
pomodux start 90m "deep focus"  # Launches TUI immediately

# Other commands work with lock validation
pomodux status    # Shows status or conflicts with running timer
pomodux stop      # Stops timer if running in same process
pomodux pause     # Pauses timer if running in same process
pomodux resume    # Resumes timer if running in same process
```

## 3. Rationale

### **Eliminates Cross-Process Synchronization**
- Single process contains both timer logic and display
- No IPC needed between timer and TUI components
- Eliminates race conditions and state inconsistencies

### **Improved User Experience**
- Immediate visual feedback when starting timer
- Consistent interface - always shows TUI when timer is active
- Simplified mental model - one process, one interface

### **Simplified Architecture**
- Reduces complexity by removing dual CLI/TUI modes  
- File-based locking handles multi-process conflicts
- Single source of truth for timer state

### **Better Plugin Integration**
- Plugins run in same process as timer and TUI
- Direct access to UI components for modal dialogs
- Simplified event handling without process boundaries

## 4. Alternatives Considered

### **Keep Dual CLI/TUI Architecture**
- **Description**: Maintain separate CLI timer process with optional TUI
- **Rejected**: Complex synchronization, poor user experience, race conditions

### **Background Daemon Approach**  
- **Description**: Timer daemon with CLI clients and TUI frontend
- **Rejected**: Over-engineered for desktop timer application, complex setup

### **Hybrid Command Mode**
- **Description**: Some commands CLI-only, others TUI-only
- **Rejected**: Confusing user experience, inconsistent interface

### **Optional TUI Flag**
- **Description**: `pomodux start --tui` for TUI mode, default CLI
- **Rejected**: Adds complexity, users would forget flag, dual maintenance

## 5. Consequences

### **Positive Consequences**
- **Simplified Architecture**: Single process eliminates IPC complexity
- **Better UX**: Immediate visual feedback and consistent interface
- **Reduced Bugs**: No cross-process race conditions or sync issues
- **Enhanced Plugins**: Direct UI access for plugin modal dialogs
- **Easier Testing**: Single process easier to test than multi-process

### **Negative Consequences**
- **Breaking Change**: Changes command behavior significantly  
- **Always Visual**: Cannot run "headless" timer without TUI
- **Terminal Required**: Requires terminal that supports TUI mode
- **Process Coupling**: Timer and display are tightly coupled

### **Mitigation Strategies**
- **Clear Documentation**: Document new command behavior prominently
- **Graceful Degradation**: TUI handles small terminal sizes appropriately
- **Lock-Based Feedback**: Other commands provide clear feedback about running timers
- **Plugin Support**: Rich plugin API enables customization for different use cases

## 6. Breaking Changes

This is a **BREAKING CHANGE** with the following command behavior changes:

### **Changed Commands**
- `pomodux start` → Now launches TUI immediately (blocking)
- `pomodux pause/resume/stop` → Only work within same process or show conflict errors

### **Removed Commands**  
- `pomodux break` → Use `pomodux start 5m "break"` 
- `pomodux long-break` → Use `pomodux start 15m "long break"`

### **New Behavior**
- All timer operations happen within TUI interface
- Cross-process operations handled via lock conflict detection
- Status commands show information about running timers in other processes

## 7. Migration Guide

### **Old Workflow**
```bash
# Old: CLI-based timer management
pomodux start 25m        # Background timer
pomodux status           # Check status  
pomodux pause            # Pause timer
pomodux tui              # Optional TUI view
```

### **New Workflow**
```bash
# New: TUI-first with immediate launch
pomodux start 25m "work" # Launches TUI immediately
# All controls available within TUI:
# [P]ause, [R]esume, [S]top, [Q]uit
```

### **For Scripts/Automation**
- Use plugins for automation instead of CLI commands
- Lock file provides timer status for external scripts
- `pomodux status` provides machine-readable JSON output

## 8. Implementation Details

### **Command Flow**
1. User runs `pomodux start 25m "work"`
2. CLI acquires file-based lock
3. CLI starts timer with session name  
4. CLI launches TUI with timer reference
5. TUI provides interactive controls
6. Lock released when TUI exits

### **Error Handling**
- Lock conflicts result in clear error messages
- TUI launch failures fall back to error display
- Terminal compatibility issues handled gracefully

## 9. Implementation Status

- **Approved** (2025-01-09)
- **Implemented** (2025-01-09)  
- **Tested**: Full integration testing with TUI components

## 10. References

- Implementation: `internal/cli/start.go`
- TUI Integration: `internal/tui/tui.go`
- Lock Management: `internal/timer/lock.go`
- [ADR 009: File-Based Timer Locking](009-file-based-timer-locking.md)

## 11. Related Components

### **Affected Components**
- `internal/cli/start.go` - TUI launch integration
- `internal/cli/pause.go` - Lock conflict detection  
- `internal/cli/resume.go` - Lock conflict detection
- `internal/cli/stop.go` - Lock conflict detection
- `internal/cli/status.go` - Cross-process status reporting
- `internal/tui/tui.go` - Primary user interface

### **Integration Points**
- File-based locking coordinates multiple process attempts
- CLI commands provide feedback about running timers
- Plugin system operates within TUI process context
- Configuration hot-reload supports TUI-first workflow