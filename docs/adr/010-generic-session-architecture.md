---
status: approved
type: technical
---

# ADR 010: Generic Session Architecture

## 1. Context / Background

The original Pomodux timer implementation used hardcoded session types through a `SessionType` enum with predefined values (`work`, `break`, `long_break`). This approach had several limitations:

- **Inflexibility**: Users could only use predefined session types
- **Poor Extensibility**: Adding new session types required code changes
- **Limited Customization**: No support for user-defined session names
- **Rigid UX**: Fixed categories didn't match diverse user workflows

Modern productivity workflows require flexible session naming to support various activities like "deep work", "meetings", "research", "code review", etc. The enum-based approach prevented users from customizing their timer experience.

## 2. Decision

**Pomodux will replace the `SessionType` enum with string-based session names, allowing users to specify any custom session name while maintaining backward compatibility through configuration defaults.**

### Architecture Changes

- **Remove `SessionType` enum** completely from codebase
- **Replace with `sessionName string`** in all timer operations
- **Update API**: `StartWithSessionName(duration, sessionName)` 
- **Plugin Events**: Pass `session_name` instead of `session_type`
- **Configuration**: Add `default_session_name` field
- **History Records**: Store session names as strings

### Session Name Rules
- Any non-empty string is valid
- Default to configured `default_session_name` if none specified  
- Preserve user input exactly (case-sensitive)
- No predefined restrictions or validation

## 3. Rationale

### **User Flexibility**
- Unlimited custom session types ("deep work", "standup", "research")
- Personalized timer experience matching individual workflows
- Support for different languages and naming conventions

### **Simplified Architecture** 
- Eliminates enum maintenance and predefined categories
- Reduces code complexity by removing type restrictions
- More intuitive API with direct string parameters

### **Future-Proof Design**
- No code changes needed for new session types
- Easy integration with external productivity systems
- Flexible foundation for advanced features

### **Plugin Compatibility**
- Plugins receive meaningful session names instead of generic types
- Better plugin integration with user's actual workflows
- Enhanced event data for plugin decision-making

## 4. Alternatives Considered

### **Extensible Enum Pattern**
- **Description**: Keep enum but make it extensible with user-defined values
- **Rejected**: Complex implementation, maintains unnecessary abstraction layer

### **Session Categories with Names**
- **Description**: Hybrid approach with categories and custom names
- **Rejected**: Added complexity without clear benefit, confusing UX

### **Predefined Set with Custom Option**
- **Description**: Default types plus "custom" option with user input
- **Rejected**: Artificial limitation, inconsistent user experience

### **Configuration-Based Enum**
- **Description**: Define session types in configuration file
- **Rejected**: Over-engineered solution for simple string requirements

## 5. Consequences

### **Positive Consequences**
- **Enhanced UX**: Users can name sessions according to their workflows
- **Simplified Code**: Removes enum complexity and type restrictions
- **Better Integration**: More meaningful data for plugins and external tools
- **Flexibility**: Supports unlimited session types without code changes
- **Internationalization**: Session names can be in any language

### **Negative Consequences** 
- **Breaking Change**: Existing configurations and plugins need updates
- **No Type Safety**: String-based approach loses compile-time validation
- **Potential Confusion**: No guidance on session name conventions

### **Mitigation Strategies**
- **Migration Support**: Clear upgrade documentation and examples
- **Sensible Defaults**: Preconfigured `default_session_name = "work"`
- **Validation**: Runtime validation prevents empty session names
- **Documentation**: Examples of effective session naming patterns

## 6. Breaking Changes

This is a **BREAKING CHANGE** with the following impacts:

### **Removed Commands**
- `pomodux break` → Use `pomodux start 5m "break"`
- `pomodux long-break` → Use `pomodux start 15m "long break"`

### **Configuration Changes**
- Added: `default_session_name` field in timer configuration
- Breaking: State and history file format changed

### **Plugin API Changes**
- Events now include `session_name` instead of `session_type`
- All existing plugins will need updates

## 7. Migration Guide

### **For Users**
```bash
# Old commands
pomodux break
pomodux long-break

# New equivalent commands  
pomodux start 5m "break"
pomodux start 15m "long break"

# New possibilities
pomodux start 25m "deep work"
pomodux start 30m "code review"
pomodux start 45m "research"
```

### **For Plugin Developers**
```lua
-- Old event handling
pomodux.register_hook("timer_started", function(event)
    if event.data.session_type == "work" then
        -- Handle work session
    end
end)

-- New event handling
pomodux.register_hook("timer_started", function(event) 
    if event.data.session_name == "work" then
        -- Handle work session
    end
    -- Now supports any session name!
end)
```

## 8. Implementation Status

- **Approved** (2025-01-09)
- **Implemented** (2025-01-09)
- **Tested**: Full test coverage with string-based sessions

## 9. References

- Implementation: `internal/timer/timer.go`, `internal/timer/engine.go`
- Configuration: `internal/config/config.go`
- Plugin Integration: `internal/plugin/manager.go`

## 10. Related Components

### **Affected Components**
- `internal/timer/timer.go` - Core timer with session name support
- `internal/timer/engine.go` - Engine API updated for session names
- `internal/timer/state.go` - State persistence with session names
- `internal/timer/history.go` - History records with session names
- `internal/config/config.go` - New `default_session_name` field
- `internal/cli/start.go` - New session name parameter
- `internal/plugin/manager.go` - Event data with session names

### **Integration Points**
- CLI commands accept optional session name parameter
- Configuration provides default session name
- Plugin events include meaningful session names
- History and state management preserve session information