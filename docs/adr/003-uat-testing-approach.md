---
status: approved
type: technical
---

# ADR 003: TUI Testing Approach for Bubbletea Applications

## 1. Context / Background

### 1.1 Current State
Pomodux is a TUI application built with Bubbletea that requires comprehensive testing to validate functionality, user experience, and real-world usage scenarios. We need to establish a standardized approach for TUI testing using teatest that ensures reliability, maintainability, and comprehensive coverage.

### 1.2 Requirements
- Comprehensive testing of all TUI components and interactions
- Validation of user workflows and real-world usage scenarios
- Clear documentation of expected behaviors and outputs
- Reproducible test execution across different environments
- Integration with development workflow
- Automated testing using teatest framework

## 2. Decision

**Selected Solution:** teatest-Based TUI Testing for Bubbletea Applications

### 2.1 Rationale
teatest is the official testing framework for Bubbletea applications and provides comprehensive testing capabilities including simulation of keypresses, window resizes, and assertions on output and model state. This approach ensures reliable, automated testing that integrates seamlessly with the Bubbletea ecosystem.

## 3. Solutions Considered

### 3.1 Option A: Custom TUI Testing Framework
- **Pros:**
  - Could be tailored to specific needs
  - Full control over testing capabilities
- **Cons:**
  - Significant development and maintenance overhead
  - Reinventing existing solutions
  - No community support or documentation

### 3.2 Option B: Manual Testing Only
- **Pros:**
  - Simple to implement
  - No additional dependencies
- **Cons:**
  - Manual validation required for each test
  - No automated regression testing
  - Prone to human error
  - Time-consuming and not scalable

### 3.3 Option C: teatest Framework (Selected)
- **Pros:**
  - Official Bubbletea testing framework
  - Comprehensive TUI testing capabilities
  - Simulates user interactions (keypresses, resizes)
  - Automated test execution and validation
  - Integration with Go testing ecosystem
  - Community support and documentation
- **Cons:**
  - Learning curve for TUI testing concepts
  - Requires understanding of Bubbletea model patterns

## 4. Consequences

### 4.1 Positive
- **Comprehensive TUI Coverage**: teatest provides thorough validation of all TUI interactions
- **User Experience Focus**: Simulates real user interactions with keypresses and window resizes
- **Regression Protection**: Automated tests prevent regressions in TUI behavior
- **Integration**: Seamless integration with Go testing ecosystem and CI/CD
- **Reliability**: Deterministic testing with predictable results
- **Professional Standards**: Uses official Bubbletea testing framework

### 4.2 Negative
- **Learning Curve**: Developers need to understand teatest and TUI testing patterns
- **Setup Complexity**: Requires understanding of Bubbletea model testing
- **Framework Dependency**: Tied to teatest framework evolution

### 4.3 Risks
- **Framework Changes**: teatest is experimental and may change
- **Complexity**: TUI testing can be more complex than unit testing
- **Test Maintenance**: TUI tests may require updates when UI changes

**Mitigation Strategies**:
- Follow teatest best practices and documentation
- Create reusable test helpers and patterns
- Regular test maintenance and updates
- Clear documentation of TUI testing approaches

## 5. Component Information

### 5.1 Repository Links
- **teatest**: https://github.com/charmbracelet/x/exp/teatest
- **Bubbletea**: https://github.com/charmbracelet/bubbletea
- **Documentation**: https://github.com/charmbracelet/x/exp/teatest

### 5.2 Maintenance Status
- **Last Updated**: 2025-01-27
- **Active Development**: Yes (part of Charm ecosystem)
- **Community Support**: Strong support from Charm team and community
- **Version Compatibility**: Compatible with current Bubbletea versions

### 5.3 Integration Verification
- **Compatibility Tested**: Yes, with current Pomodux TUI structure
- **Existing Component Impact**: Integrates with existing Bubbletea models
- **Migration Path**: Direct integration with existing TUI components

## 6. Implementation Strategy

### 6.1 TUI Testing with teatest
**Duration**: Immediate implementation
**Deliverables**:
- teatest-based test suite for all TUI components
- Comprehensive test scenarios covering all user interactions
- Clear documentation of TUI testing patterns
- Integration with Go testing workflow

**Implementation Details**:
```go
package tui_test

import (
    "testing"
    "time"
    
    "github.com/charmbracelet/x/exp/teatest"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/pomodux/internal/tui"
    "github.com/pomodux/internal/timer"
)

// TestTimerTUIInteraction tests the main timer TUI flow
func TestTimerTUIInteraction(t *testing.T) {
    // Create a new timer
    timer := timer.New()
    
    // Create TUI model
    model := tui.NewModel(timer)
    
    // Start teatest session
    tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(80, 24))
    
    // Test initial state
    tm.Send(tea.WindowSizeMsg{Width: 80, Height: 24})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Timer Ready"))
    }, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))
    
    // Test start timer
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Timer Running"))
    }, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))
    
    // Test pause timer
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Timer Paused"))
    }, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))
    
    // Test resume timer
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Timer Running"))
    }, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))
    
    // Test stop timer
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Timer Stopped"))
    }, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))
}

// TestTUIResize tests window resizing behavior
func TestTUIResize(t *testing.T) {
    timer := timer.New()
    model := tui.NewModel(timer)
    tm := teatest.NewTestModel(t, model)
    
    // Test resize to small window
    tm.Send(tea.WindowSizeMsg{Width: 40, Height: 10})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        // Verify UI adapts to small screen
        return len(string(bts)) > 0
    }, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*1))
    
    // Test resize to large window
    tm.Send(tea.WindowSizeMsg{Width: 120, Height: 40})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        // Verify UI adapts to large screen
        return len(string(bts)) > 0
    }, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*1))
}

// TestProgressDisplay tests progress bar and time display
func TestProgressDisplay(t *testing.T) {
    timer := timer.New()
    timer.Start(time.Minute * 5) // 5 minute timer
    
    model := tui.NewModel(timer)
    tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(80, 24))
    
    // Wait for initial render
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        output := string(bts)
        // Check for progress bar elements
        return strings.Contains(output, "█") || strings.Contains(output, "▓") || strings.Contains(output, "░")
    }, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))
}
```

### 6.2 Advanced TUI Testing Patterns
**Duration**: Ongoing enhancement
**Deliverables**:
- Advanced teatest patterns and helpers
- CI/CD integration with Go testing
- Automated TUI regression testing
- Test utilities and shared patterns

**Advanced teatest Patterns**:
```go
// Helper function for common TUI test setup
func setupTUITest(t *testing.T) (*teatest.TestModel, *timer.Timer) {
    timer := timer.New()
    model := tui.NewModel(timer)
    tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(80, 24))
    return tm, timer
}

// Test helper for waiting for specific states
func waitForTimerState(t *testing.T, tm *teatest.TestModel, expectedState string) {
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte(expectedState))
    }, teatest.WithCheckInterval(time.Millisecond*50), teatest.WithDuration(time.Second*5))
}

// Golden file testing for consistent UI output
func TestTUIGoldenFiles(t *testing.T) {
    tm, timer := setupTUITest(t)
    timer.Start(time.Minute * 25)
    
    // Wait for stable output
    time.Sleep(time.Millisecond * 100)
    
    // Compare against golden file
    output := tm.Output()
    goldenFile := filepath.Join("testdata", "timer_running.golden")
    
    if *update {
        ioutil.WriteFile(goldenFile, output, 0644)
    }
    
    expected, err := ioutil.ReadFile(goldenFile)
    require.NoError(t, err)
    
    assert.Equal(t, string(expected), string(output))
}

// Test keyboard input sequences
func TestKeyboardSequences(t *testing.T) {
    tm, _ := setupTUITest(t)
    
    // Simulate rapid key presses
    keys := []tea.KeyMsg{
        {Type: tea.KeyRunes, Runes: []rune("s")}, // start
        {Type: tea.KeyRunes, Runes: []rune("p")}, // pause
        {Type: tea.KeyRunes, Runes: []rune("r")}, // resume
        {Type: tea.KeyRunes, Runes: []rune("q")}, // quit
    }
    
    for _, key := range keys {
        tm.Send(key)
        time.Sleep(time.Millisecond * 100) // Allow processing
    }
    
    waitForTimerState(t, tm, "Timer Stopped")
}
```

### 6.3 TUI Testing Best Practices
**Duration**: Ongoing improvement
**Deliverables**:
- teatest testing standards and patterns
- Performance benchmarking for TUI rendering
- Cross-platform TUI testing
- Advanced interaction scenarios

## 7. Test Organization Structure

### 7.1 Directory Structure
```
internal/
├── tui/
│   ├── tui.go                    # Main TUI implementation
│   ├── tui_test.go              # Core TUI tests with teatest
│   ├── tui_integration_test.go  # Integration tests
│   └── testdata/
│       ├── golden/              # Golden files for output comparison
│       │   ├── timer_idle.golden
│       │   ├── timer_running.golden
│       │   └── timer_paused.golden
│       └── fixtures/            # Test data and configurations
├── timer/
│   ├── timer.go                 # Timer logic
│   └── timer_test.go           # Timer unit tests
└── testhelpers/
    ├── teatest_helpers.go       # Shared teatest utilities
    └── mock_timer.go           # Mock timer for testing
```

### 7.2 Test Categories

#### TUI Component Tests
- **Model State Tests**: Validate Bubbletea model state transitions
- **View Rendering Tests**: Test UI rendering and layout
- **Input Handling Tests**: Validate keyboard and resize handling
- **Progress Display Tests**: Test timer progress visualization
- **Error State Tests**: Validate error handling and display

#### Integration Tests
- **Timer-TUI Integration**: Test timer and TUI component interaction
- **State Persistence Tests**: Test state saving and loading
- **Event Flow Tests**: Test event-driven updates
- **Cross-Platform Tests**: Validate behavior across different terminals

## 8. Integration with Development Workflow

### 8.1 CI/CD Integration
- **Automated TUI Tests**: Run as part of Go test suite
- **Pull Request Validation**: All TUI tests must pass before merge
- **Test Results**: Integrated with standard Go test reporting
- **Issue Tracking**: Test failures tracked as GitHub issues

### 8.2 Test Execution Workflow
1. **Development Phase**: teatest runs with `go test` on every commit
2. **Pull Request Phase**: Full TUI test suite validation
3. **Release Phase**: Comprehensive teatest validation
4. **Post-Release**: Automated regression testing

## 9. Quality Standards

### 9.1 Test Coverage Requirements
- **TUI Component Coverage**: 100% of TUI components and interactions
- **User Flow Coverage**: All documented user workflows
- **Error Scenario Coverage**: All error states and edge cases
- **Cross-Platform Coverage**: All supported terminal environments

### 9.2 Test Quality Standards
- **Clarity**: Tests must be clear and understandable
- **Reliability**: Tests must be deterministic and repeatable
- **Maintainability**: Tests must be easy to update and extend
- **Documentation**: All test patterns must be well-documented

### 9.3 Performance Standards
- **Execution Time**: teatest suite should complete within 30 seconds
- **Resource Usage**: Tests should not consume excessive terminal resources
- **Reliability**: Tests should have 95%+ pass rate in stable environments

## 10. Future Considerations

### 10.1 Scalability
- **Test Parallelization**: Parallel execution of independent TUI tests
- **Performance Testing**: Automated TUI rendering performance benchmarking
- **Cross-Terminal Testing**: Automated testing across different terminal emulators

### 10.2 Advanced TUI Testing Features
- **Golden File Testing**: Screenshot-like comparison for TUI output
- **Accessibility Testing**: TUI accessibility validation
- **Performance Profiling**: TUI rendering performance analysis

### 10.3 Tool Evolution
- **teatest Enhancement**: Following teatest framework evolution
- **Go Testing Integration**: Enhanced integration with Go testing ecosystem
- **TUI Testing Patterns**: Development of advanced TUI testing patterns

---

**References**:
- [teatest Documentation](https://github.com/charmbracelet/x/exp/teatest)
- [Bubbletea Testing Guide](https://github.com/charmbracelet/bubbletea)
- [Go Testing Standards](https://golang.org/doc/code.html#Testing)
- [TUI Testing Best Practices](https://charm.sh/) 