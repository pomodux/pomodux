---
status: approved
type: technical
---

# ADR 009: File-Based Timer Locking Strategy

## 1. Context / Background

The TUI timer feature requires ensuring that only one timer instance can run system-wide, preventing conflicts when multiple processes attempt to start timers simultaneously. Without proper synchronization, users could accidentally start multiple timers, leading to:

- Inconsistent state management
- Conflicting plugin event emissions
- Confusing user experience with multiple timer displays
- Race conditions in state file writes

Previous timer implementation relied on process-level state management without cross-process coordination, making it possible to start multiple timer instances.

## 2. Decision

**Pomodux will implement file-based locking using XDG-compliant lock files with process validation to ensure single timer instance enforcement across all processes.**

### Implementation Details

- **Lock File Location**: `$XDG_RUNTIME_DIR/pomodux/timer.lock` (XDG-compliant)
- **Lock File Format**: JSON containing PID, session name, start time, and duration
- **Atomic Operations**: Use `O_CREATE|O_EXCL` for atomic lock acquisition
- **Process Validation**: Validate lock owner process existence before respecting locks
- **Orphan Recovery**: Automatically detect and recover from orphaned locks
- **Corruption Handling**: Delete corrupted lock files and retry operation

### Lock File Structure
```json
{
  "pid": 12345,
  "session_name": "deep work",
  "start_time": "2025-01-09T10:30:00Z",
  "duration_seconds": 1500,
  "locked_at": "2025-01-09T10:30:00Z"
}
```

## 3. Rationale

### **Cross-Process Coordination**
- File-based locks work across all process boundaries
- No dependency on shared memory or IPC mechanisms
- Platform-independent solution for process coordination

### **Robustness and Recovery**
- Process validation prevents stale lock issues
- Automatic orphan detection handles crashed processes
- Corruption recovery ensures system reliability

### **XDG Compliance**
- Follows Linux desktop standards for runtime files
- Proper cleanup when user session ends
- Consistent with other desktop applications

### **Performance and Simplicity**
- Minimal overhead for lock operations (< 1ms typical)
- Simple implementation without external dependencies
- Atomic file operations prevent race conditions

## 4. Alternatives Considered

### **Process Signals (POSIX)**
- **Rejected**: Platform-specific, complex signal handling, not available on Windows

### **Shared Memory Segments**
- **Rejected**: Complex cleanup, permission issues, platform-specific implementation

### **TCP/UDP Sockets**
- **Rejected**: Overkill for local coordination, firewall/permission complications

### **Database-Based Locking**
- **Rejected**: Heavy dependency, unnecessary complexity for simple use case

### **PID Files Only**
- **Rejected**: No session information, difficult to provide user feedback

## 5. Consequences

### **Positive Consequences**
- **Single Instance Guarantee**: Robust prevention of timer conflicts
- **User Experience**: Clear error messages when conflicts occur
- **System Reliability**: Automatic recovery from process crashes
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Standards Compliant**: Follows XDG Base Directory Specification

### **Negative Consequences**
- **Complexity**: Additional file system operations and error handling
- **Disk I/O**: Small overhead for lock file operations
- **Edge Cases**: Rare scenarios with corrupted locks or permission issues

### **Mitigation Strategies**
- **Comprehensive Testing**: Multi-process conflict scenarios
- **Timeout Handling**: Operations complete within 1 second
- **Error Recovery**: Automatic corruption detection and cleanup
- **User Feedback**: Clear error messages for conflict scenarios

## 6. Implementation Status

- **Approved** (2025-01-09)  
- **Implemented** (2025-01-09)
- **Tested**: Full test coverage including edge cases

## 7. References

- [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
- Implementation: `internal/timer/lock.go`
- Tests: `internal/timer/lock_test.go`

## 8. Related Components

### **Affected Components**
- `internal/timer/lock.go` - Lock manager implementation  
- `internal/cli/start.go` - Lock acquisition during timer start
- `internal/cli/status.go` - Lock conflict detection
- `internal/cli/pause.go` - Lock validation for operations
- `internal/cli/resume.go` - Lock validation for operations
- `internal/cli/stop.go` - Lock validation for operations

### **Integration Points**
- Timer startup sequence with lock acquisition
- CLI command error handling for lock conflicts
- State management coordination with lock lifecycle