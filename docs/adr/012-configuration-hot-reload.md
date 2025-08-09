---
status: approved
type: technical
---

# ADR 012: Configuration Hot-Reload Pattern

## 1. Context / Background

The TUI timer feature introduces flexible session naming where users can specify any custom session name. To support seamless user customization, the system needs a way to reload configuration changes without requiring application restart.

Previous configuration loading happened once at application startup, meaning users had to:
- Restart the application to see configuration changes
- Remember their preferred default session names
- Manually specify session names even for common patterns

With the new generic session architecture, users need the ability to:
- Set their preferred `default_session_name` in configuration  
- See changes take effect immediately on next timer start
- Customize session defaults for different contexts or time periods

## 2. Decision  

**Pomodux will implement configuration hot-reload where configuration is reloaded on every `pomodux start` command, allowing users to customize default session names without application restart.**

### Implementation Approach

- **Reload Trigger**: Every invocation of `pomodux start`
- **Reload Timing**: During lock acquisition phase, before timer startup
- **Scope**: Full configuration reload including all sections
- **Error Handling**: Graceful fallback to previous config on reload failures
- **Performance**: Acceptable overhead (~1ms) for improved UX

### Configuration Lifecycle
1. User modifies `~/.config/pomodux/config.yaml`
2. User runs `pomodux start` (with or without session name)
3. System reloads configuration automatically
4. New `default_session_name` used if no session specified
5. Timer starts with current configuration values

## 3. Rationale

### **Seamless User Experience**
- No application restart needed for configuration changes
- Immediate feedback for configuration modifications
- Encourages experimentation with different session naming patterns

### **Session Customization Support**
- Users can quickly adjust `default_session_name` for different contexts
- Supports workflow changes without friction (work vs. study vs. personal)
- Enables dynamic session naming patterns

### **Simplified Operations**
- Eliminates restart requirement for configuration changes
- Reduces user cognitive load - changes "just work"
- Consistent with modern application behavior expectations

### **Development Workflow**
- Easier testing of different configuration values
- Faster development cycle for configuration-related features
- Better debugging experience for configuration issues

## 4. Alternatives Considered

### **File System Watching**
- **Description**: Monitor config file with fsnotify and reload automatically
- **Rejected**: Complex implementation, resource overhead, timing issues

### **Reload Command**
- **Description**: Add `pomodux reload-config` command
- **Rejected**: Extra step for users, easy to forget, inconsistent UX

### **Startup-Only Loading**  
- **Description**: Keep existing single load at startup (status quo)
- **Rejected**: Poor UX for session name customization, requires restarts

### **Cache with TTL**
- **Description**: Cache configuration with time-based expiration
- **Rejected**: Arbitrary timing, potential for stale configuration

### **Manual Reload Flag**
- **Description**: `pomodux start --reload-config` flag option
- **Rejected**: Users would forget flag, inconsistent behavior

## 5. Consequences

### **Positive Consequences**
- **Enhanced UX**: Configuration changes take effect immediately
- **Reduced Friction**: No restart required for session customization  
- **Better Adoption**: Users more likely to customize session names
- **Improved Workflow**: Supports dynamic session naming patterns
- **Development Friendly**: Easier testing and debugging

### **Negative Consequences**
- **Performance Overhead**: Small I/O cost on every timer start (~1ms)
- **Error Complexity**: Need to handle configuration reload failures
- **Validation Timing**: Configuration errors discovered at runtime instead of startup
- **State Consistency**: Need to ensure configuration changes don't break running timers

### **Mitigation Strategies**
- **Performance**: Configuration reload is fast (~1ms) and only on timer start
- **Error Handling**: Graceful fallback to previous configuration on errors
- **Validation**: Clear error messages for invalid configuration
- **Documentation**: Examples of effective configuration patterns

## 6. Implementation Details

### **Reload Sequence**
1. `pomodux start` command invoked
2. Initialize lock manager
3. **Reload configuration** (`config.Load()`)
4. Use reloaded config for default session name
5. Acquire lock with session information
6. Start timer with current configuration

### **Error Handling**
```go
cfg, err := config.Load()
if err != nil {
    // Graceful fallback - use previous config or sensible defaults
    logger.Warn("Configuration reload failed, using defaults", 
                map[string]interface{}{"error": err.Error()})
    cfg = getDefaultConfig()
}
```

### **Performance Considerations**
- YAML parsing: ~0.5ms for typical config file
- File I/O: ~0.5ms for local filesystem read
- Total overhead: ~1ms - acceptable for user-initiated operation

## 7. Configuration Scope

### **Reloaded Settings**
- `timer.default_session_name` - Primary use case
- `timer.default_work_duration` - Also benefits from hot-reload  
- `timer.default_break_duration` - Consistency
- All other configuration sections - Comprehensive reload

### **Not Affected by Reload**
- Running timer state - Current session continues with original config
- Plugin system - Requires restart for plugin directory changes
- Log configuration - Requires restart for logger reconfiguration

## 8. Usage Examples

### **Session Name Customization**
```yaml
# Morning work session
timer:
  default_session_name: "morning focus"

# Later, edit config for afternoon
timer:  
  default_session_name: "afternoon tasks"

# Changes take effect on next 'pomodux start'
```

### **Context Switching**
```bash
# Edit config for study context
vim ~/.config/pomodux/config.yaml
# Set: default_session_name: "study"

pomodux start 25m  # Uses "study" as session name

# Later, switch to work context  
vim ~/.config/pomodux/config.yaml
# Set: default_session_name: "work"

pomodux start 45m  # Uses "work" as session name
```

## 9. Implementation Status

- **Approved** (2025-01-09)
- **Implemented** (2025-01-09)
- **Tested**: Configuration reload functionality verified

## 10. References

- Implementation: `internal/cli/start.go` (config reload)
- Configuration: `internal/config/config.go` (Load function)
- [ADR 010: Generic Session Architecture](010-generic-session-architecture.md)

## 11. Related Components

### **Affected Components**
- `internal/cli/start.go` - Configuration reload on timer start
- `internal/config/config.go` - Configuration loading and validation
- `internal/timer/lock.go` - Uses reloaded config during lock acquisition

### **Integration Points**  
- Timer startup sequence includes configuration reload
- Default session name resolution uses current configuration
- Error handling provides fallbacks for configuration issues
- Performance optimization ensures minimal impact on timer start