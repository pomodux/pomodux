# Pomodux Configuration Guide

## Overview

Pomodux uses a YAML configuration file to manage application settings. The configuration file is automatically created on first run and follows XDG standards for file location.

## Configuration File Location

### Default Location
- **Linux/macOS**: `~/.config/pomodux/config.yaml`
- **Windows**: `%APPDATA%\pomodux\config.yaml`

### Custom Location
You can specify a custom configuration file using the `--config` flag:
```bash
pomodux --config /path/to/custom/config.yaml start 25m
```

## Configuration File Structure

The configuration file is organized into several sections, each controlling different aspects of the application:

```yaml
# Timer Configuration
timer:
  default_work_duration: 25m
  default_break_duration: 5m
  default_long_break_duration: 15m
  default_session_name: "work"
  auto_start_breaks: false

# TUI Configuration
tui:
  theme: "default"
  key_bindings:
    start: "s"
    stop: "q"
    pause: "p"
    resume: "r"

# Notification Configuration
notifications:
  enabled: true
  sound: false
  message: "Timer completed!"

# Plugin Configuration
plugins:
  directory: "~/.config/pomodux/plugins"
  enabled:
    example_plugin: true
    another_plugin: false

# Logging Configuration
logging:
  level: "info"
  format: "text"
  output: "file"
  log_file: ""
  show_caller: false
```

## Configuration Sections

### Timer Configuration

Controls the default timer durations and behavior.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `default_work_duration` | Duration | `25m` | Default duration for work sessions |
| `default_break_duration` | Duration | `5m` | Default duration for break sessions |
| `default_long_break_duration` | Duration | `15m` | Default duration for long break sessions |
| `default_session_name` | String | `"work"` | Default session name when none specified |
| `auto_start_breaks` | Boolean | `false` | Automatically start breaks after work sessions |

**Duration Format**: Use Go duration format (e.g., `25m`, `1h30m`, `90s`)

**Example**:
```yaml
timer:
  default_work_duration: 45m
  default_break_duration: 10m
  default_long_break_duration: 20m
  default_session_name: "work"
  auto_start_breaks: true
```

### TUI Configuration

Controls the Terminal User Interface appearance and key bindings.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `theme` | String | `"default"` | TUI theme (currently only "default" supported) |
| `key_bindings` | Map | See defaults | Custom key bindings for TUI controls |

**Default Key Bindings**:
```yaml
key_bindings:
  start: "s"
  stop: "q"
  pause: "p"
  resume: "r"
```

**Example**:
```yaml
tui:
  theme: "default"
  key_bindings:
    start: "space"
    stop: "x"
    pause: "p"
    resume: "r"
```

### Notification Configuration

Controls notification behavior when timers complete.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | Boolean | `true` | Enable/disable notifications |
| `sound` | Boolean | `false` | Enable/disable sound notifications |
| `message` | String | `"Timer completed!"` | Default notification message |

**Example**:
```yaml
notifications:
  enabled: true
  sound: true
  message: "Time to take a break!"
```

### Plugin Configuration

Controls plugin system behavior and individual plugin settings.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `directory` | String | XDG plugins dir | Directory containing plugin files |
| `enabled` | Map | `{}` | Enable/disable specific plugins |

**Plugin Directory**: 
- Default: `~/.config/pomodux/plugins`
- Supports tilde expansion (`~`) and environment variables
- Can be absolute or relative path

**Plugin Enable/Disable**:
```yaml
plugins:
  directory: "~/.config/pomodux/plugins"
  enabled:
    my_plugin: true
    debug_plugin: false
    test_plugin: true
```

### Logging Configuration

Controls application logging behavior and output.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `level` | String | `"info"` | Logging level (debug, info, warn, error) |
| `format` | String | `"text"` | Log format (text, json) |
| `output` | String | `"file"` | Log output (console, file, both) |
| `log_file` | String | `""` | Custom log file path (empty = auto) |
| `show_caller` | Boolean | `false` | Show function caller in log messages |

**Log Levels**:
- `debug`: Detailed debug information
- `info`: General information messages
- `warn`: Warning messages
- `error`: Error messages only

**Log Formats**:
- `text`: Human-readable text format
- `json`: Structured JSON format

**Log Output**:
- `console`: Output to terminal
- `file`: Output to log file
- `both`: Output to both console and file

**Example**:
```yaml
logging:
  level: "debug"
  format: "json"
  output: "both"
  log_file: "~/.local/share/pomodux/debug.log"
  show_caller: true
```

## Environment Variables

### XDG Configuration
Pomodux respects XDG standards for configuration location:

- `XDG_CONFIG_HOME`: Override default config directory
- `XDG_DATA_HOME`: Override default data directory
- `XDG_RUNTIME_DIR`: Override default runtime directory

### Example Environment Setup
```bash
export XDG_CONFIG_HOME="$HOME/.config"
export XDG_DATA_HOME="$HOME/.local/share"
export XDG_RUNTIME_DIR="/tmp/pomodux-$USER"
```

## Path Expansion

Configuration supports path expansion for user convenience:

### Tilde Expansion
```yaml
plugins:
  directory: "~/my-plugins"  # Expands to /home/user/my-plugins
```

### Environment Variables
```yaml
logging:
  log_file: "$HOME/.local/share/pomodux/app.log"
```

## Configuration Validation

Pomodux validates configuration on load and reports errors for:

- Invalid duration values (must be positive)
- Invalid log levels (debug, info, warn, error only)
- Invalid log formats (text, json only)
- Invalid log outputs (console, file, both only)
- Invalid file paths (no path traversal allowed)

## Configuration Examples

### Minimal Configuration
```yaml
timer:
  default_work_duration: 25m
```

### Development Configuration
```yaml
timer:
  default_work_duration: 5m
  default_break_duration: 1m
  default_long_break_duration: 3m
  default_session_name: "dev"

logging:
  level: "debug"
  format: "json"
  output: "both"
  show_caller: true

plugins:
  enabled:
    debug_plugin: true
    test_plugin: true
```

### Production Configuration
```yaml
timer:
  default_work_duration: 45m
  default_break_duration: 15m
  default_long_break_duration: 30m
  default_session_name: "focus"
  auto_start_breaks: true

notifications:
  enabled: true
  sound: true
  message: "Time for a break!"

logging:
  level: "info"
  format: "text"
  output: "file"
  log_file: "/var/log/pomodux/app.log"

plugins:
  enabled:
    productivity_plugin: true
    analytics_plugin: true
```

## Troubleshooting

### Common Issues

**Configuration Not Found**
```bash
# Check if config file exists
ls -la ~/.config/pomodux/config.yaml

# Create default config
pomodux --help  # This will create default config
```

**Invalid Configuration**
```bash
# Check configuration syntax
pomodux --config /path/to/config.yaml start 25m

# Look for error messages in output
```

**Plugin Directory Issues**
```bash
# Check plugin directory exists
ls -la ~/.config/pomodux/plugins/

# Create plugin directory if missing
mkdir -p ~/.config/pomodux/plugins/
```

### Debug Configuration Loading
```bash
# Enable debug logging
pomodux --config debug-config.yaml start 25m

# Where debug-config.yaml contains:
logging:
  level: "debug"
  format: "text"
  output: "console"
```

## Migration Notes

### From Previous Versions
- Configuration format remains compatible
- New options are added with sensible defaults
- Deprecated options are logged as warnings
- Configuration validation is enhanced

### Breaking Changes
- None in current version
- Future versions may deprecate specific options
- Deprecation warnings will be shown in logs

## Related Documentation

- [Plugin Development Guide](plugin-development.md) - Plugin configuration details
- [Development Setup](development-setup.md) - Development environment configuration
- [ADR](adr/) - Architecture decisions affecting configuration 