# Release 0.4.0 - TUI API for Lua Plugins

## Overview
This release introduces a Go-based TUI API for Lua plugins, allowing plugins to display dialogs, notifications, and other UI elements within the Pomodux terminal interface.

## Features
- Expose TUI functions to Lua plugins
- Example plugin demonstrating TUI usage
- Documentation for plugin authors

## Implementation Plan
- [ ] Design TUI API
- [ ] Select Go TUI library
- [ ] Implement Go TUI API
- [ ] Expose TUI API to Lua plugins
- [ ] Update plugin manager for TUI calls
- [ ] Write example plugin
- [ ] Update documentation
- [ ] Add tests

## Testing
- [ ] Unit tests for TUI API
- [ ] Integration tests with plugins

## Documentation
- [ ] Update technical specifications
- [ ] Update development setup
- [ ] Document plugin API

## Retrospective
- [ ] Lessons learned
- [ ] User feedback
- [ ] Documentation and rules audit
- [ ] Improvement proposals 

## Plugin Configuration Structure Update

This release introduces a more flexible way to enable or disable plugins in your `config.yaml` file.

### New Style: Plugin-Specific Configuration Blocks

You can now enable or disable plugins using plugin-specific blocks:

```yaml
plugins:
    kimai_integration:
        enabled: false
    directory: /home/ritchie/.config/pomodux/plugins
```

### Old Style (Still Supported)

The previous style is still supported for backward compatibility:

```yaml
plugins:
    enabled:
        kimai_integration: false
    directory: /home/ritchie/.config/pomodux/plugins
```

### Benefits
- More flexible and extensible for future plugin options
- Cleaner, more organized config file
- Easier to manage per-plugin settings
- Both styles are recognized by the application 