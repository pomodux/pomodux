# Pomodux UAT Test Suite

This directory contains User Acceptance Testing (UAT) scripts for the Pomodux application, implementing the testing strategy defined in [ADR 003: UAT Testing Approach](docs/adr/003-uat-testing-approach.md).

## Overview

The UAT test suite provides manual testing approaches for user experience validation and acceptance testing of the single-process TUI application.

## Directory Structure

```
tests/uat/
├── manual/
│   └── uat-script.sh              # Enhanced manual UAT script
└── fixtures/
    └── configs/                   # Test configuration files
```

## Testing Strategy

Since Pomodux has been refactored to a single-process TUI-only application, the testing approach focuses on:

1. **Go Unit Tests** - Comprehensive automated testing (`make test`)
2. **Manual UAT Scripts** - User experience validation and acceptance testing

## Manual UAT Testing

### Enhanced UAT Script

The enhanced manual UAT script provides comprehensive testing with automated validation:

```bash
# Run the enhanced UAT script
./tests/uat/manual/uat-script.sh
```

**Features:**
- ✅ Automated test execution and validation
- ✅ Comprehensive error handling and logging
- ✅ Test result reporting and metrics
- ✅ Proper test environment setup and cleanup
- ✅ Structured test organization
- ✅ Integration with release management process

**Test Categories:**
- Basic functionality testing (help, version)
- Configuration management
- TUI timer operations (start, pause, resume, stop via keyboard)
- Session types (work, break, long break)
- Session history and recording
- Error handling and edge cases
- Single timer enforcement (file-based locking)

## Go Unit Testing

The primary automated testing is handled by the Go test suite:

```bash
# Run all unit tests
make test

# Run tests with coverage
make test-coverage
```

**Test Coverage:**
- **Config Tests**: Configuration loading, validation, and management
- **Logger Tests**: Logging functionality and levels
- **Plugin Tests**: Plugin system, loading, and event emission
- **Timer Tests**: Core timer functionality, state management, and locking
- **Lock Tests**: File-based single timer enforcement

## Test Configuration

### Test Environment Setup

All tests automatically:
- Build the application (`make build`)
- Set up test configuration files
- Backup existing configuration
- Clean up after completion

### Test Configuration Files

Test configurations use short durations for quick execution:

```json
{
    "timer": {
        "default_work_duration": "1m",
        "default_break_duration": "30s",
        "default_long_break_duration": "2m",
        "default_session_name": "work"
    }
}
```

### Test Isolation

Each test runs in isolation:
- Configuration is backed up and restored
- Session history is cleaned between tests
- No interference between test runs
- Single timer enforcement prevents conflicts

## Integration with Release Process

### Gate 3 Integration

The UAT test suite integrates with the 4-gate approval process:

1. **Unit Tests**: Automated Go tests run on every build
2. **Manual UAT**: Required for Gate 3 approval
3. **Test Results**: Documented in release documentation
4. **Issue Tracking**: Test failures tracked as GitHub issues

### Release Testing Workflow

1. **Development Phase**: Unit tests run on every commit (`make test`)
2. **UAT Phase**: Manual UAT scripts executed by stakeholders
3. **Release Phase**: Both test suites run for final validation
4. **Post-Release**: Automated tests run for regression protection

## Quality Standards

### Test Coverage Requirements

- **Feature Coverage**: 100% of user-facing features
- **TUI Functionality**: All keyboard controls and interactions
- **Error Scenario Coverage**: All documented error conditions
- **Workflow Coverage**: All documented user workflows

### Test Quality Standards

- **Clarity**: Tests must be clear and understandable
- **Reliability**: Tests must be deterministic and repeatable
- **Maintainability**: Tests must be easy to update and extend
- **Documentation**: All tests must be well-documented

### Performance Standards

- **Execution Time**: Manual UAT should complete within 30 minutes
- **Resource Usage**: Tests should not consume excessive resources
- **Reliability**: Tests should have 95%+ pass rate in stable environments

## Troubleshooting

### Common Issues

#### Build Failures
```bash
# Ensure application builds successfully
make build
```

#### Test Failures
1. Check if application is built: `make build`
2. Check configuration: `./bin/pomodux --help`
3. Check for running timers: Look for lock files in XDG runtime directory
4. Use TUI controls (P/R/S/Q) to pause/resume/stop/quit timers

#### Permission Issues
```bash
# Make scripts executable
chmod +x tests/uat/manual/uat-script.sh
```

### Debug Mode

Run tests with verbose output:

```bash
# Manual UAT with debug output
./tests/uat/manual/uat-script.sh 2>&1 | tee uat-debug.log

# Go tests with verbose output
go test -v ./...
```

## Single-Process TUI Architecture

The Pomodux application now uses a single-process TUI-first architecture:

- **TUI-Only Interface**: `pomodux start` immediately launches interactive TUI
- **Event-Driven**: Timer uses event notifications instead of polling
- **Single Timer Enforcement**: File-based locking prevents multiple instances
- **No CLI Commands**: All timer control happens within TUI interface
- **Backwards Compatibility**: `pomodux start 25m "work"` syntax still supported

## Testing the TUI

### Manual TUI Testing

Test the interactive TUI interface:

1. **Start Timer**: `./bin/pomodux 1m test`
2. **Keyboard Controls**: Test P(ause), R(esume), S(top), Q(uit)
3. **Progress Display**: Verify real-time progress updates
4. **Session Completion**: Let timer complete naturally
5. **Multiple Instance**: Try starting second timer (should fail)

### TUI Test Scenarios

- **Normal Completion**: Start timer and let it complete
- **User Interruption**: Start timer and stop with 'S' or 'Q'
- **Pause/Resume**: Test pause with 'P' and resume with 'R'
- **Lock Enforcement**: Try starting multiple timers simultaneously
- **Error Handling**: Test invalid durations, missing arguments

## Future Enhancements

### Planned Improvements

1. **Visual Testing**: Screenshot comparison for TUI components
2. **Performance Testing**: Automated performance benchmarking
3. **Cross-Platform Testing**: Multi-platform test execution
4. **Continuous Integration**: GitHub Actions integration
5. **TUI Automation**: Automated TUI interaction testing

### Advanced Features

- **Accessibility Testing**: TUI accessibility validation
- **Security Testing**: File permission and lock security
- **Load Testing**: Multiple process coordination testing
- **Advanced Reporting**: Enhanced test analytics and metrics

## 🔗 Related Documentation

- **[Release Management](docs/release-management.md)** - Release process and approval gates
- **[Requirements](../../docs/requirements.md)** - Project requirements and specifications
- **[ADR 011](docs/adr/011-tui-first-command-architecture.md)** - TUI-first architecture decision

## Contributing

When adding new tests:

1. **Follow Naming Conventions**: Use descriptive test names
2. **Maintain Isolation**: Ensure tests don't interfere with each other
3. **Add Documentation**: Document test purpose and expected behavior
4. **Update Coverage**: Ensure new features are covered by tests
5. **Run All Tests**: Verify no regressions (`make test`)

## Support

For issues with the test suite:

1. Check the troubleshooting section above
2. Review the Go test output for specific failure details
3. Check the [ADR 003](docs/adr/003-uat-testing-approach.md) for testing strategy
4. Create a GitHub issue with detailed error information