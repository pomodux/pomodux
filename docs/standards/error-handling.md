# Error Handling Standards for Pomodux

## Overview

This document defines error handling patterns and standards for the Pomodux project, following Go best practices and industry standards for CLI applications.

## Core Principles

### 1. Error Wrapping

**Always wrap errors with context** using `fmt.Errorf` with `%w` verb:

```go
// Good: Wrapped error with context
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
    }
    
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
    }
    
    return &config, nil
}
```

**Why:**
- Preserves original error for debugging
- Adds context about what operation failed
- Enables error unwrapping with `errors.Unwrap()`
- Better error messages for users

### 2. Error Types

**Use sentinel errors for expected conditions:**

```go
package config

import "errors"

var (
    ErrConfigNotFound = errors.New("config file not found")
    ErrInvalidConfig  = errors.New("invalid config file")
    ErrMissingTimers  = errors.New("no timer presets defined")
)
```

**Usage:**
```go
func LoadConfig(path string) (*Config, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, ErrConfigNotFound
    }
    // ...
}

// Caller can check for specific errors
config, err := LoadConfig(path)
if errors.Is(err, config.ErrConfigNotFound) {
    // Create default config
}
```

**Custom error types for structured errors:**

```go
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error in field %q: %s (value: %v)", 
        e.Field, e.Message, e.Value)
}
```

### 3. Exit Codes

**Follow Unix/Linux exit code conventions:**

| Exit Code | Meaning | Usage |
|-----------|---------|-------|
| 0 | Success | Normal completion |
| 1 | General error | Unspecified error, operation failed |
| 2 | Usage error | Invalid command-line arguments |
| 130 | Interrupted | SIGINT (Ctrl+C) received |

**Implementation:**
```go
package main

import (
    "os"
    "syscall"
)

func main() {
    if err := run(); err != nil {
        handleError(err)
        os.Exit(1)
    }
}

func handleError(err error) {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    
    // Check for specific error types
    if errors.Is(err, cli.ErrInvalidCommand) {
        fmt.Fprintf(os.Stderr, "Run 'pomodux --help' for usage\n")
        os.Exit(2) // Usage error
    }
}
```

**Signal Handling:**
```go
func main() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        sig := <-sigChan
        // Cleanup and save state
        handleInterrupt(sig)
        os.Exit(130) // SIGINT exit code
    }()
    
    // Main application logic
}
```

## Error Handling Patterns

### 1. Function-Level Error Handling

**Pattern: Return errors, don't panic (except for programming errors):**

```go
// Good: Return error
func ParseDuration(s string) (time.Duration, error) {
    if s == "" {
        return 0, fmt.Errorf("duration cannot be empty")
    }
    
    d, err := time.ParseDuration(s)
    if err != nil {
        return 0, fmt.Errorf("invalid duration format %q: %w", s, err)
    }
    
    if d <= 0 {
        return 0, fmt.Errorf("duration must be positive, got %v", d)
    }
    
    if d > 24*time.Hour {
        return 0, fmt.Errorf("duration exceeds maximum (24h), got %v", d)
    }
    
    return d, nil
}
```

**Never panic for user input errors:**
```go
// Bad: Panic on user error
func ParseDuration(s string) time.Duration {
    d, err := time.ParseDuration(s)
    if err != nil {
        panic(err) // Don't do this!
    }
    return d
}
```

### 2. Error Messages

**User-Facing Errors:**
- Clear and actionable
- Suggest how to fix the problem
- Include relevant context (file paths, values)
- Avoid technical jargon

```go
// Good: Clear, actionable error
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf(
                "config file not found: %s\n"+
                "Create a config file at %s or run 'pomodux --help' for setup instructions",
                path, path)
        }
        return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
    }
    // ...
}
```

**Developer-Facing Errors (Logging):**
- Include stack traces for debugging
- Use structured logging with context
- Log at appropriate levels (DEBUG, INFO, WARN, ERROR)

```go
import "github.com/sirupsen/logrus"

func LoadConfig(path string) (*Config, error) {
    logger := logrus.WithField("config_path", path)
    
    data, err := os.ReadFile(path)
    if err != nil {
        logger.WithError(err).Error("Failed to read config file")
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    logger.Debug("Config file loaded successfully")
    // ...
}
```

### 3. Error Propagation

**Don't swallow errors:**

```go
// Bad: Error ignored
func SaveHistory(history *History, path string) {
    data, _ := json.Marshal(history) // Error ignored!
    os.WriteFile(path, data, 0644)    // Error ignored!
}
```

**Good: Propagate errors up the call stack:**

```go
func SaveHistory(history *History, path string) error {
    data, err := json.Marshal(history)
    if err != nil {
        return fmt.Errorf("failed to marshal history: %w", err)
    }
    
    if err := os.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("failed to write history file %s: %w", path, err)
    }
    
    return nil
}
```

### 4. Graceful Degradation

**Continue operation with defaults when non-critical errors occur:**

```go
func LoadConfig(path string) (*Config, error) {
    config, err := loadConfigFile(path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            // Non-critical: use defaults
            logrus.WithError(err).Warn("Config file not found, using defaults")
            return DefaultConfig(), nil
        }
        // Critical: return error
        return nil, err
    }
    
    // Validate and fix config
    if err := validateConfig(config); err != nil {
        logrus.WithError(err).Warn("Invalid config, using defaults for invalid fields")
        applyDefaults(config)
    }
    
    return config, nil
}
```

## Error Handling by Component

### 1. CLI Commands

**Pattern: Handle errors at command level, provide user-friendly messages:**

```go
func startTimer(c *cli.Context) error {
    args := c.Args()
    if args.Len() == 0 {
        return cli.Exit("Error: duration or preset required\n"+
            "Usage: pomodux start <duration|preset> [label]", 2)
    }
    
    durationOrPreset := args.Get(0)
    label := args.Get(1)
    
    timer, err := createTimer(durationOrPreset, label)
    if err != nil {
        return cli.Exit(fmt.Sprintf("Error: %v\n"+
            "Run 'pomodux --help' for usage", err), 1)
    }
    
    if err := runTimer(timer); err != nil {
        return fmt.Errorf("timer failed: %w", err)
    }
    
    return nil
}
```

### 2. Configuration Loading

**Pattern: Validate and provide helpful errors:**

```go
func LoadConfig(path string) (*Config, error) {
    // Check file exists
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file not found: %s\n"+
            "Create a config file or use defaults", path)
    }
    
    // Read file
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    // Parse YAML
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("invalid YAML in config file: %w", err)
    }
    
    // Validate
    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    
    return &config, nil
}
```

### 3. File I/O Operations

**Pattern: Handle filesystem errors gracefully:**

```go
func SaveState(state *TimerState, path string) error {
    // Create directory if needed
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create state directory: %w", err)
    }
    
    // Atomic write: write to temp file, then rename
    tmpPath := path + ".tmp"
    data, err := json.Marshal(state)
    if err != nil {
        return fmt.Errorf("failed to marshal state: %w", err)
    }
    
    if err := os.WriteFile(tmpPath, data, 0600); err != nil {
        return fmt.Errorf("failed to write state file: %w", err)
    }
    
    if err := os.Rename(tmpPath, path); err != nil {
        return fmt.Errorf("failed to save state file: %w", err)
    }
    
    return nil
}
```

### 4. Timer Operations

**Pattern: Validate inputs, return clear errors:**

```go
func NewTimer(duration time.Duration, label string) (*Timer, error) {
    if duration <= 0 {
        return nil, fmt.Errorf("duration must be positive, got %v", duration)
    }
    
    if duration > 24*time.Hour {
        return nil, fmt.Errorf("duration exceeds maximum (24h), got %v", duration)
    }
    
    if len(label) > 200 {
        return nil, fmt.Errorf("label too long (max 200 chars), got %d", len(label))
    }
    
    return &Timer{
        duration: duration,
        label:    label,
        state:    StateIdle,
    }, nil
}
```

### 5. TUI Error Handling

**Pattern: Display errors in TUI, don't crash:**

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case errorMsg:
        // Display error in TUI
        m.error = msg.Error()
        return m, nil
        
    case tea.KeyMsg:
        if msg.String() == "p" {
            if err := m.timer.Pause(); err != nil {
                // Return error message, don't panic
                return m, func() tea.Msg {
                    return errorMsg{err}
                }
            }
        }
    }
    
    return m, nil
}
```

## Logging Errors

### Structured Logging

**Use logrus with structured fields:**

```go
import "github.com/sirupsen/logrus"

func LoadConfig(path string) (*Config, error) {
    logger := logrus.WithFields(logrus.Fields{
        "component": "config",
        "path":      path,
    })
    
    data, err := os.ReadFile(path)
    if err != nil {
        logger.WithError(err).Error("Failed to read config file")
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    logger.Debug("Config file loaded")
    // ...
}
```

### Log Levels

**Use appropriate log levels:**

- **DEBUG**: Detailed information for debugging
- **INFO**: General informational messages
- **WARN**: Warning messages (non-critical errors)
- **ERROR**: Error messages (operation failed)

```go
// Non-critical: Warn and continue
if err := loadOptionalConfig(); err != nil {
    logrus.WithError(err).Warn("Optional config not loaded, using defaults")
}

// Critical: Error and return
if err := loadRequiredConfig(); err != nil {
    logrus.WithError(err).Error("Failed to load required config")
    return err
}
```

## Error Testing

### Test Error Cases

**Always test error conditions:**

```go
func TestParseDuration_Errors(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"empty string", "", true},
        {"invalid format", "25 minutes", true},
        {"negative duration", "-5m", true},
        {"zero duration", "0s", true},
        {"too large", "25h", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := ParseDuration(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Test Error Wrapping

**Verify error wrapping preserves original error:**

```go
func TestLoadConfig_WrapsErrors(t *testing.T) {
    _, err := LoadConfig("/nonexistent/path")
    assert.Error(t, err)
    
    // Verify original error is preserved
    assert.True(t, errors.Is(err, os.ErrNotExist))
    
    // Verify error message includes context
    assert.Contains(t, err.Error(), "config file")
    assert.Contains(t, err.Error(), "/nonexistent/path")
}
```

## Checklist

### Before Committing

- [ ] All errors are wrapped with context
- [ ] Exit codes follow Unix conventions
- [ ] User-facing errors are clear and actionable
- [ ] Errors are logged at appropriate levels
- [ ] No errors are silently ignored
- [ ] Error cases are tested

### Code Review

- [ ] Error messages are user-friendly
- [ ] Errors are properly wrapped
- [ ] Exit codes are correct
- [ ] Graceful degradation is implemented where appropriate
- [ ] Error handling doesn't leak implementation details

## References

- [Go Error Handling Best Practices](https://go.dev/blog/error-handling-and-go)
- [Effective Go - Errors](https://go.dev/doc/effective_go#errors)
- [Go 1.13 Error Wrapping](https://go.dev/blog/go1.13-errors)
- [Unix Exit Codes](https://tldp.org/LDP/abs/html/exitcodes.html)

---

**Last Updated:** 2025-01-15  
**Next Review:** 2025-02-15

