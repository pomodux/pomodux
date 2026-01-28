# Testing Strategy for Pomodux

## Overview

This document defines the testing standards and practices for the Pomodux project. We follow a **Test-Driven Development (TDD)** approach with a goal of **70% code coverage** for the MVP release.

## Testing Philosophy

### Test-Driven Development (TDD)

We follow the TDD cycle:
1. **Red**: Write a failing test
2. **Green**: Write minimal code to pass the test
3. **Refactor**: Improve code while keeping tests passing

**Benefits:**
- Ensures code is testable from the start
- Documents expected behavior
- Prevents regressions
- Guides design toward better architecture

### Coverage Goal

**Target: 65% - 70% code coverage**

**Rationale:**
- Balances thoroughness with development velocity
- Focuses on critical paths (timer logic, state persistence, config)
- Allows pragmatic exceptions for hard-to-test code (TUI rendering, terminal I/O)

**Coverage Exclusions:**
- Main entry points (`main()` functions)
- TUI rendering internals (low-level lipgloss styling, covered partially by `teatest` integration tests)
- Terminal-specific I/O (handled by Bubbletea framework)
- Error message formatting (low risk, high maintenance)

## Testing Framework

### Core Testing

**Standard Library:** Use Go's built-in `testing` package

```go
package timer

import "testing"

func TestTimerStart(t *testing.T) {
    // Test implementation
}
```

### Assertions

**Library:** `github.com/stretchr/testify/assert`

**Rationale:**
- Cleaner test code than manual `if` statements
- Better error messages
- Widely used in Go community

**Usage:**
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestTimerDuration(t *testing.T) {
    timer := NewTimer(25 * time.Minute)
    assert.Equal(t, 25*time.Minute, timer.Duration())
    assert.True(t, timer.IsRunning())
}
```

### Test Utilities

**Time Mocking:** Use `github.com/benbjohnson/clock` for time-dependent tests

**Rationale:**
- Timer accuracy is critical
- Need to test time-based logic without waiting
- Prevents flaky tests from system clock changes

**Usage:**
```go
import "github.com/benbjohnson/clock"

func TestTimerCompletion(t *testing.T) {
    mockClock := clock.NewMock()
    timer := NewTimerWithClock(25*time.Minute, mockClock)
    
    timer.Start()
    mockClock.Add(25 * time.Minute)
    
    assert.True(t, timer.IsCompleted())
}
```

## Test Organization

### File Structure

**Naming Convention:** `*_test.go` files alongside source files

```
internal/
├── timer/
│   ├── timer.go
│   ├── timer_test.go
│   ├── state.go
│   └── state_test.go
├── config/
│   ├── config.go
│   └── config_test.go
```

### Test Categories

#### 1. Unit Tests

**Scope:** Test individual functions and methods in isolation

**Requirements:**
- Fast execution (<100ms per test)
- No external dependencies (filesystem, network)
- Deterministic (same input = same output)
- Mock external dependencies

**Example:**
```go
func TestParseDuration(t *testing.T) {
    tests := []struct {
        input    string
        expected time.Duration
        hasError bool
    }{
        {"25m", 25 * time.Minute, false},
        {"1h30m", 90 * time.Minute, false},
        {"invalid", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            result, err := ParseDuration(tt.input)
            if tt.hasError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}
```

#### 2. Integration Tests

**Scope:** Test component interactions and external systems

**Requirements:**
- Test file I/O operations
- Test config loading/saving
- Test history persistence
- Use temporary directories/files

**Example:**
```go
func TestConfigLoad(t *testing.T) {
    tmpDir := t.TempDir()
    configPath := filepath.Join(tmpDir, "config.yaml")
    
    // Write test config
    err := os.WriteFile(configPath, []byte(testConfigYAML), 0644)
    require.NoError(t, err)
    
    // Load config
    config, err := LoadConfig(configPath)
    assert.NoError(t, err)
    assert.Equal(t, 25*time.Minute, config.Timers["work"])
}
```

#### 3. TUI Tests

**Scope:** Test Bubbletea model logic and rendered output

**Approach:**
- Test model state transitions via `Update()` calls
- Test message handling and command generation
- Test rendered output using `teatest` for integration-level TUI verification

##### 3a. Unit Tests (Model Logic)

Test state transitions and message handling without rendering:

```go
func TestModelPause(t *testing.T) {
    m := initialModel()
    m.timer.Start()

    // Simulate 'p' keypress
    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}
    newModel, cmd := m.Update(msg)

    assert.True(t, newModel.timer.IsPaused())
    assert.Nil(t, cmd) // No command needed
}
```

##### 3b. Integration Tests (Rendered Output via teatest)

**Library:** `github.com/charmbracelet/x/exp/teatest`

**Rationale:**
- TUI rendering bugs (layout, theming, component visibility) cannot be caught by model-logic unit tests alone
- `teatest` runs a real Bubbletea program in a virtual terminal, enabling assertions on rendered output
- Prevents regressions in visual output across refactors and theme changes
- Avoids reliance on manual verification for every change

**What to test with teatest:**
- Correct components render in each state (running, paused, completed)
- Time display shows expected format (mm:ss)
- Status indicator reflects current state
- Control legend visibility
- Confirmation dialog appears on stop key
- Terminal size warning renders when terminal is too small
- Theme colors apply correctly (spot-check via golden files)

**Example:**
```go
func TestRunningView(t *testing.T) {
    m := initialModel()
    tm := teatest.NewModel(t, m, teatest.WithInitialTermSize(80, 24))

    // Send a start key or init message
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

    // Wait for output and assert on rendered content
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "RUNNING")
    }, teatest.WithDuration(3*time.Second))
}
```

**Golden file testing** can be used for full-screen snapshots:
```go
func TestRunningViewGolden(t *testing.T) {
    m := initialModel()
    tm := teatest.NewModel(t, m, teatest.WithInitialTermSize(80, 24))

    // Allow the model to render
    time.Sleep(500 * time.Millisecond)

    out := tm.FinalOutput(t)
    teatest.RequireEqualOutput(t, out) // Compares against testdata golden file
}
```

**Note:** Golden file tests are more brittle and should be used sparingly for critical layouts. Prefer string-contains assertions for most cases.

## Test Coverage

### Coverage Measurement

**Tool:** `go test -cover`

**Target:** 70% overall coverage

**Coverage by Package:**
- `internal/timer`: 85%+ (critical logic)
- `internal/config`: 80%+ (user-facing)
- `internal/history`: 75%+ (data integrity)
- `internal/tui`: 60%+ (model logic only, rendering excluded)
- `internal/theme`: 70%+ (theme application)
- `internal/plugin`: 65%+ (MVP: event emission only)

### Coverage Exclusions

**Justified Exclusions:**
```go
// main.go - entry point, tested via integration
//go:build !test

// TUI rendering - visual output, tested manually
func (m model) View() string {
    // ... rendering code ...
}

// Error message formatting - low risk
func formatError(err error) string {
    // ... formatting ...
}
```

**Marking Exclusions:**
```go
// Exclude from coverage: TUI rendering is visual-only
// Coverage: off
func (m model) View() string {
    // ...
}
```

## Test Data Management

### Fixtures

**Location:** `testdata/` directory

**Structure:**
```
testdata/
├── config/
│   ├── valid.yaml
│   ├── invalid.yaml
│   └── missing-timers.yaml
├── history/
│   ├── empty.json
│   ├── single-session.json
│   └── corrupted.json
└── state/
    ├── running.json
    └── paused.json
```

**Usage:**
```go
func TestLoadConfig(t *testing.T) {
    configPath := "testdata/config/valid.yaml"
    config, err := LoadConfig(configPath)
    // ...
}
```

### Temporary Files

**Use `t.TempDir()` for test-specific files:**

```go
func TestSaveHistory(t *testing.T) {
    tmpDir := t.TempDir()
    historyPath := filepath.Join(tmpDir, "history.json")
    
    history := NewHistory()
    err := history.Save(historyPath)
    assert.NoError(t, err)
    
    // Verify file exists and is valid
    data, err := os.ReadFile(historyPath)
    assert.NoError(t, err)
    assert.ValidJSON(t, data)
}
```

## Test Execution

### Running Tests

**All Tests:**
```bash
go test ./...
```

**With Coverage:**
```bash
go test -cover ./...
```

**Coverage Report:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Verbose Output:**
```bash
go test -v ./...
```

**Specific Package:**
```bash
go test ./internal/timer
```

### CI Integration

**Coverage Threshold:**
- Fail CI if coverage drops below 70%
- Report coverage percentage in CI output
- Generate coverage badge for README

**Example CI Step:**
```yaml
- name: Run tests with coverage
  run: |
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out | grep total | awk '{print $3}'
```

## Test Quality Standards

### Test Naming

**Convention:** `Test<FunctionName>_<Scenario>`

**Examples:**
```go
func TestTimerStart_ValidDuration(t *testing.T)
func TestTimerStart_InvalidDuration(t *testing.T)
func TestTimerPause_WhileRunning(t *testing.T)
func TestConfigLoad_MissingFile(t *testing.T)
```

### Test Structure

**AAA Pattern:** Arrange, Act, Assert

```go
func TestTimerResume(t *testing.T) {
    // Arrange
    timer := NewTimer(25 * time.Minute)
    timer.Start()
    timer.Pause()
    
    // Act
    timer.Resume()
    
    // Assert
    assert.False(t, timer.IsPaused())
    assert.True(t, timer.IsRunning())
}
```

### Test Independence

**Requirements:**
- Tests must not depend on execution order
- Tests must not share state
- Each test should be runnable in isolation
- Use `t.Parallel()` for independent tests

```go
func TestTimerIsolation(t *testing.T) {
    t.Parallel() // Safe to run in parallel
    
    timer := NewTimer(25 * time.Minute)
    // ... test implementation ...
}
```

### Error Testing

**Test Both Success and Failure Cases:**

```go
func TestParseDuration(t *testing.T) {
    t.Run("valid duration", func(t *testing.T) {
        d, err := ParseDuration("25m")
        assert.NoError(t, err)
        assert.Equal(t, 25*time.Minute, d)
    })
    
    t.Run("invalid duration", func(t *testing.T) {
        d, err := ParseDuration("invalid")
        assert.Error(t, err)
        assert.Zero(t, d)
    })
    
    t.Run("empty string", func(t *testing.T) {
        d, err := ParseDuration("")
        assert.Error(t, err)
        assert.Zero(t, d)
    })
}
```

## TDD Workflow

### Development Cycle

1. **Write Test First:**
   ```go
   func TestTimerStart(t *testing.T) {
       timer := NewTimer(25 * time.Minute)
       timer.Start()
       assert.True(t, timer.IsRunning())
   }
   ```

2. **Run Test (Should Fail):**
   ```bash
   go test ./internal/timer
   # Expected: compilation error or test failure
   ```

3. **Write Minimal Implementation:**
   ```go
   func (t *Timer) Start() {
       t.running = true
   }
   ```

4. **Run Test (Should Pass):**
   ```bash
   go test ./internal/timer
   # Expected: test passes
   ```

5. **Refactor:**
   - Improve code quality
   - Extract common logic
   - Ensure tests still pass

### When to Skip TDD

**Exceptions (Justified):**
- Exploratory coding (spike solutions)
- Prototyping TUI layout
- Learning new framework APIs
- **But:** Write tests before committing

## Testing Checklist

### Before Committing

- [ ] All tests pass: `go test ./...`
- [ ] Coverage meets 70% threshold
- [ ] No skipped tests (unless justified)
- [ ] Tests are independent and parallelizable
- [ ] Test names are descriptive
- [ ] Error cases are tested
- [ ] Edge cases are tested
- [ ] No test data pollution (use temp files)

### Code Review

- [ ] Tests cover new functionality
- [ ] Tests are readable and maintainable
- [ ] Test failures provide clear error messages
- [ ] No hardcoded values (use constants or fixtures)
- [ ] Tests follow AAA pattern
- [ ] Coverage exclusions are justified

## Tools and Resources

### Required Tools

- `go test` - Standard Go testing tool
- `github.com/stretchr/testify` - Assertions
- `github.com/benbjohnson/clock` - Time mocking
- `github.com/charmbracelet/x/exp/teatest` - TUI integration testing

### Optional Tools

- `golang.org/x/tools/cmd/cover` - Coverage analysis
- `github.com/axw/gocov` - Coverage reporting
- `github.com/fatih/gomodifytags` - Test tag management

### Documentation

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [TDD Best Practices](https://golang.org/doc/effective_go#testing)

## Continuous Improvement

### Regular Reviews

- Review test coverage monthly
- Identify untested critical paths
- Refactor tests for maintainability
- Update strategy based on learnings

### Metrics to Track

- Overall coverage percentage
- Coverage by package
- Test execution time
- Flaky test frequency
- Test maintenance burden

---

**Last Updated:** 2025-01-15  
**Next Review:** 2025-02-15

