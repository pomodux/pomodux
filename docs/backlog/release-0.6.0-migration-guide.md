# Release 0.6.0 Migration Guide

## Overview

Release 0.6.0 introduces significant architectural changes to Pomodux, transitioning from a CLI+TUI dual interface to a unified TUI-only system. This migration guide provides step-by-step instructions for adapting to the new architecture.

## Breaking Changes Summary

### **Major Changes**
- **Session Types**: Enum-based session types replaced with generic session names
- **Plugin API**: Complete redesign of plugin system (existing plugins will not work)
- **CLI Commands**: `pomodux break` and `pomodux long-break` commands removed
- **Session History**: Existing session history will be cleared
- **TUI Architecture**: New global stage pattern for component management

### **Impact Assessment**
- **High Impact**: Plugin developers (complete API redesign)
- **Medium Impact**: Users with custom configurations
- **Low Impact**: Basic timer usage (session names instead of types)

## Migration Steps

### **Step 1: Backup Current Data**

Before upgrading, backup your current Pomodux data:

```bash
# Backup configuration
cp ~/.config/pomodux/config.yaml ~/.config/pomodux/config.yaml.backup

# Backup session history (if you want to preserve it)
cp ~/.config/pomodux/session_history.json ~/.config/pomodux/session_history.json.backup

# Backup plugins (if you have custom plugins)
cp -r ~/.config/pomodux/plugins ~/.config/pomodux/plugins.backup
```

### **Step 2: Install Release 0.6.0**

Follow the standard installation process for Release 0.6.0:

```bash
# Download and install the new release
# (Follow installation instructions for your platform)
```

### **Step 3: Update Session Usage**

#### **Before (Release 0.5.x and earlier)**
```bash
# Start work session
pomodux start 25m

# Start break session
pomodux break

# Start long break session
pomodux long-break
```

#### **After (Release 0.6.0)**
```bash
# Start work session
pomodux start 25m "work"

# Start break session
pomodux start 5m "break"

# Start long break session
pomodux start 15m "long break"

# Start custom session
pomodux start 45m "deep work"
pomodux start 30m "meeting"
pomodux start 10m "quick break"
```

### **Step 4: Update Configuration (Optional)**

The configuration file structure remains compatible, but you may want to update session durations:

#### **Before**
```yaml
timer:
  default_work_duration: 25m
  default_break_duration: 5m
  default_long_break_duration: 15m
```

#### **After**
```yaml
timer:
  default_work_duration: 25m
  # Note: break and long-break durations are now used as defaults
  # when no session name is provided
  default_break_duration: 5m
  default_long_break_duration: 15m
```

### **Step 5: Plugin Migration (For Plugin Developers)**

#### **Plugin API Changes**

**Before (tview-based dialogs)**
```lua
-- Show notification
pomodux.show_notification("Timer completed!")

-- Show list selection
local choice = pomodux.show_list_selection("Select option", {"Option 1", "Option 2"})

-- Show input prompt
local input = pomodux.show_input_prompt("Enter name", "", "Your name")
```

**After (Bubbletea-based components)**
```lua
-- Show notification (auto-dismissing)
pomodux.show_notification("Timer completed!", 5) -- 5 seconds

-- Show modal dialog
pomodux.show_modal("Information", "This is a modal dialog")

-- Update status information
pomodux.update_status("Current task: Deep work session")
```

#### **Plugin Development Guide**

For detailed plugin development information, see:
- [Plugin Development Guide](../plugin-development.md) - Updated for new API
- [Test Plugin Examples](../examples/test_plugins/) - Reference implementations

### **Step 6: Verify Migration**

Test the new functionality:

```bash
# Test basic timer functionality
pomodux start 1m "test session"

# Test TUI interface
pomodux start 2m "tui test"

# Test configuration loading
pomodux --config ~/.config/pomodux/config.yaml start 1m "config test"
```

## Configuration Changes

### **New Configuration Options**

Release 0.6.0 introduces new configuration options for the TUI:

```yaml
tui:
  theme: "default"
  key_bindings:
    start: "s"
    stop: "q"
    pause: "p"
    resume: "r"
  # New options for global stage
  window_management:
    modal_timeout: 30  # seconds
    notification_duration: 5  # seconds
```

### **Deprecated Configuration**

The following configuration options are deprecated but remain functional:

```yaml
# These still work but are no longer the primary interface
timer:
  auto_start_breaks: false  # Now handled by session names
```

## Troubleshooting

### **Common Issues**

#### **"Command not found" for break/long-break**
**Issue**: `pomodux break` and `pomodux long-break` commands no longer exist.

**Solution**: Use session names instead:
```bash
# Instead of: pomodux break
pomodux start 5m "break"

# Instead of: pomodux long-break  
pomodux start 15m "long break"
```

#### **Plugin errors after upgrade**
**Issue**: Existing plugins fail to load or function incorrectly.

**Solution**: Plugins need to be updated for the new API:
- Check plugin compatibility with Release 0.6.0
- Update plugins to use new Bubbletea-based API
- Test plugins with new TUI architecture

#### **Session history missing**
**Issue**: Previous session history is not visible.

**Solution**: Session history format has changed:
- Old history is cleared during upgrade
- New sessions will use the updated format
- History is now stored with session names instead of types

#### **TUI not displaying correctly**
**Issue**: TUI interface appears broken or incomplete.

**Solution**: Check terminal compatibility:
- Ensure terminal supports UTF-8
- Verify terminal size is adequate (minimum 80x24)
- Check for conflicting terminal configurations

### **Debug Mode**

Enable debug logging to troubleshoot issues:

```yaml
logging:
  level: "debug"
  format: "text"
  output: "console"
```

```bash
# Run with debug logging
pomodux start 1m "debug test"
```

## Rollback Plan

If you need to rollback to a previous version:

### **Step 1: Uninstall Release 0.6.0**
```bash
# Remove the new version
# (Follow uninstall instructions for your platform)
```

### **Step 2: Restore Previous Version**
```bash
# Install previous version (e.g., 0.5.x)
# (Follow installation instructions for previous version)
```

### **Step 3: Restore Configuration**
```bash
# Restore backed up configuration
cp ~/.config/pomodux/config.yaml.backup ~/.config/pomodux/config.yaml

# Restore session history (if desired)
cp ~/.config/pomodux/session_history.json.backup ~/.config/pomodux/session_history.json
```

### **Step 4: Restore Plugins**
```bash
# Restore backed up plugins
cp -r ~/.config/pomodux/plugins.backup ~/.config/pomodux/plugins
```

## Support

### **Documentation**
- [Configuration File Specifications](../configuration_file_specifications.md)
- [Plugin Development Guide](../plugin-development.md)
- [Release Notes](../releases/release-0.6.0.md)

### **Testing**
- Run UAT tests to verify functionality
- Test plugin integration with new API
- Validate configuration changes

### **Reporting Issues**
- Check existing issues in the project repository
- Create new issues with detailed information
- Include configuration files and error logs

## Migration Checklist

- [ ] Backup current configuration and data
- [ ] Install Release 0.6.0
- [ ] Update session usage patterns
- [ ] Test basic timer functionality
- [ ] Test TUI interface
- [ ] Update plugin code (if applicable)
- [ ] Verify configuration changes
- [ ] Test plugin integration
- [ ] Update documentation and scripts
- [ ] Remove old CLI command references

## Summary

Release 0.6.0 represents a significant architectural improvement that simplifies the codebase while expanding capabilities. The migration process is straightforward for most users, with the main changes being:

1. **Session naming**: Use descriptive session names instead of predefined types
2. **TUI-first**: All interactions now go through the TUI interface
3. **Plugin updates**: Plugins need to be updated for the new API
4. **Configuration**: Minimal configuration changes required

The new architecture provides a more consistent and maintainable foundation for future development while delivering an improved user experience. 