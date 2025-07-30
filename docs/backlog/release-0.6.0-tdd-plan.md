# Release 0.6.0 TDD Plan

> **STATUS: PLANNING**
> 
> **Document Type:** Test-Driven Development Plan
> **Release:** [Release 0.6.0](release-0.6.0.md)
> **Last Updated:** 2025-01-27

---

## Overview

This document defines the Test-Driven Development (TDD) approach for Release 0.6.0 - TUI-Only Timer with File-Based Locking. All implementation must follow the TDD process: Red (write failing test), Green (write minimal code to pass), Refactor (clean up while keeping tests passing).

## Current Testing Infrastructure

### Existing Test Coverage
- **Timer Tests**: Comprehensive unit tests in `internal/timer/timer_test.go`
- **Plugin Tests**: Unit and integration tests in `internal/plugin/manager_test.go`
- **Config Tests**: Configuration validation tests in `internal/config/config_test.go`
- **UAT Tests**: Automated tests using bats framework in `tests/uat/automated/`
- **Integration Tests**: End-to-end tests in `tests/uat/automated/integration/`

### Missing Test Infrastructure
- **TUI Tests**: No existing TUI tests using teatest framework
- **Lock Manager Tests**: No existing lock manager (new component)
- **Global Stage Tests**: No existing stage management tests
- **Cross-Platform Tests**: Limited cross-platform testing

## TDD Process Requirements

### Core TDD Principles
1. **Write failing test first (Red)**: Define desired behavior before implementation
2. **Write minimal code to pass (Green)**: Implement only what's needed to make test pass
3. **Refactor while keeping tests passing**: Clean up code without breaking tests
4. **Repeat**: Continue cycle for each new feature or behavior

### Test Coverage Requirements
- **Overall Coverage**: Minimum 80% for all code
- **Critical Paths**: Minimum 95% for lock manager and timer core
- **Public APIs**: 100% coverage for all public interfaces
- **Error Handling**: 100% coverage for error paths
- **Integration Tests**: Comprehensive component interaction testing

## Phase 1: Core Timer Simplification + Lock Manager

### 1.1 Lock Manager TDD

#### Test Suite: `internal/timer/lock_test.go`

**Test Categories:**

1. **Lock Acquisition Tests**
```go
func TestTimerLockManager_AcquireLock_Success(t *testing.T)
func TestTimerLockManager_AcquireLock_AlreadyLocked(t *testing.T)
func TestTimerLockManager_AcquireLock_InvalidSessionName(t *testing.T)
func TestTimerLockManager_AcquireLock_FileSystemError(t *testing.T)
```

2. **Process Validation Tests**
```go
func TestTimerLockManager_ValidateProcess_ValidProcess(t *testing.T)
func TestTimerLockManager_ValidateProcess_InvalidProcess(t *testing.T)
func TestTimerLockManager_ValidateProcess_ProcessNotFound(t *testing.T)
func TestTimerLockManager_ValidateProcess_CrossPlatform(t *testing.T)
```

3. **Orphaned Lock Recovery Tests**
```go
func TestTimerLockManager_RecoverOrphanedLock_Success(t *testing.T)
func TestTimerLockManager_RecoverOrphanedLock_ValidProcess(t *testing.T)
func TestTimerLockManager_RecoverOrphanedLock_CorruptedLockFile(t *testing.T)
func TestTimerLockManager_RecoverOrphanedLock_PermissionError(t *testing.T)
```

4. **Lock File Location Tests**
```go
func TestGetLockDir_XDGCompliant(t *testing.T)
func TestGetLockDir_FallbackToStateDir(t *testing.T)
func TestGetLockDir_CreateDirectoryIfNotExists(t *testing.T)
func TestGetLockDir_PermissionHandling(t *testing.T)
```

**Mock Dependencies:**
- File system operations (using `testing/fstest` or custom mocks)
- Process operations (using `os/exec` mocks)
- System calls (using platform-specific mocks)

**Integration Tests:**
```go
func TestTimerLockManager_EndToEnd_ConcurrentAccess(t *testing.T)
func TestTimerLockManager_EndToEnd_ProcessCrash(t *testing.T)
func TestTimerLockManager_EndToEnd_CrossProcess(t *testing.T)
```

### 1.2 Timer Core Refactoring TDD

#### Test Suite: `internal/timer/timer_test.go`

**Test Categories:**

1. **Session Name Integration Tests**
```go
func TestTimer_StartWithSessionName_Success(t *testing.T)
func TestTimer_StartWithSessionName_DefaultSessionName(t *testing.T)
func TestTimer_StartWithSessionName_EmptySessionName(t *testing.T)
func TestTimer_StartWithSessionName_SpecialCharacters(t *testing.T)
```

2. **State Management Tests**
```go
func TestTimer_SaveState_WithSessionName(t *testing.T)
func TestTimer_LoadState_WithSessionName(t *testing.T)
func TestTimer_StateMigration_FromSessionType(t *testing.T)
func TestTimer_StateConsistency_WithSessionName(t *testing.T)
```

3. **Event System Tests**
```go
func TestTimer_EmitEvent_WithSessionName(t *testing.T)
func TestTimer_EventData_IncludesSessionName(t *testing.T)
func TestTimer_EventLogging_StructuredData(t *testing.T)
func TestTimer_EventTiming_AccurateTimestamps(t *testing.T)
```

4. **Backward Compatibility Tests**
```go
func TestTimer_BackwardCompatibility_ExistingState(t *testing.T)
func TestTimer_BackwardCompatibility_ExistingHistory(t *testing.T)
func TestTimer_BackwardCompatibility_ExistingPlugins(t *testing.T)
```

**Mock Dependencies:**
- Plugin manager (using interface mocks)
- State manager (using file system mocks)
- Logger (using test logger implementation)

### 1.3 CLI Command Refactoring TDD

#### Test Suite: `internal/cli/start_test.go`

**Test Categories:**

1. **Command Interface Tests**
```go
func TestStartCmd_WithSessionName(t *testing.T)
func TestStartCmd_WithoutSessionName(t *testing.T)
func TestStartCmd_InvalidDuration(t *testing.T)
func TestStartCmd_InvalidSessionName(t *testing.T)
```

2. **Integration Tests**
```go
func TestStartCmd_Integration_TimerStart(t *testing.T)
func TestStartCmd_Integration_LockAcquisition(t *testing.T)
func TestStartCmd_Integration_ErrorHandling(t *testing.T)
```

## Phase 2: TUI-Only Refactor + Global Stage

### 2.1 Global Stage TDD

#### Test Suite: `internal/tui/stage_test.go`

**Test Categories:**

1. **Stage Management Tests**
```go
func TestGlobalStage_New_Singleton(t *testing.T)
func TestGlobalStage_Update_StageEvents(t *testing.T)
func TestGlobalStage_Update_WindowResize(t *testing.T)
func TestGlobalStage_Update_KeyPress(t *testing.T)
```

2. **Migration from Existing TUI Tests**
```go
func TestGlobalStage_MigrateExistingTUI(t *testing.T)
func TestGlobalStage_BackwardCompatibility(t *testing.T)
func TestGlobalStage_SessionNameDisplay(t *testing.T)
```

2. **Component Management Tests**
```go
func TestGlobalStage_AddTimerWindow(t *testing.T)
func TestGlobalStage_AddPluginModal(t *testing.T)
func TestGlobalStage_AddNotification(t *testing.T)
func TestGlobalStage_RemoveComponent(t *testing.T)
```

3. **Event Distribution Tests**
```go
func TestGlobalStage_DistributeEvent_TimerEvent(t *testing.T)
func TestGlobalStage_DistributeEvent_PluginEvent(t *testing.T)
func TestGlobalStage_DistributeEvent_SystemEvent(t *testing.T)
```

**Mock Dependencies:**
- Bubbletea framework (using teatest)
- Terminal operations (using test terminal)
- Event system (using channel mocks)

### 2.2 Timer Window TDD

#### Test Suite: `internal/tui/timer_window_test.go`

**Test Categories:**

1. **Rendering Tests**
```go
func TestTimerWindow_Render_BasicDisplay(t *testing.T)
func TestTimerWindow_Render_WithSessionName(t *testing.T)
func TestTimerWindow_Render_PausedState(t *testing.T)
func TestTimerWindow_Render_ProgressBar(t *testing.T)
```

2. **Layout Tests**
```go
func TestTimerWindow_Layout_Centered(t *testing.T)
func TestTimerWindow_Layout_Responsive(t *testing.T)
func TestTimerWindow_Layout_MinimumSize(t *testing.T)
func TestTimerWindow_Layout_MaximumSize(t *testing.T)
```

3. **Interaction Tests**
```go
func TestTimerWindow_HandleKeyPress_Pause(t *testing.T)
func TestTimerWindow_HandleKeyPress_Resume(t *testing.T)
func TestTimerWindow_HandleKeyPress_Stop(t *testing.T)
func TestTimerWindow_HandleKeyPress_Quit(t *testing.T)
```

**Integration Tests:**
```go
func TestTimerWindow_Integration_WithStage(t *testing.T)
func TestTimerWindow_Integration_WithTimer(t *testing.T)
func TestTimerWindow_Integration_ResponsiveLayout(t *testing.T)
```

### 2.3 TUI Integration TDD

#### Test Suite: `internal/tui/tui_test.go`

**Test Categories:**

1. **TUI Launch Tests**
```go
func TestRunTUI_Launch_Success(t *testing.T)
func TestRunTUI_Launch_WithSessionName(t *testing.T)
func TestRunTUI_Launch_ErrorHandling(t *testing.T)
func TestRunTUI_Launch_TerminalState(t *testing.T)
```

2. **Integration Tests**
```go
func TestRunTUI_Integration_TimerStart(t *testing.T)
func TestRunTUI_Integration_UserInteraction(t *testing.T)
func TestRunTUI_Integration_PluginIntegration(t *testing.T)
```

## Phase 3: Plugin API Redesign

### 3.1 Plugin API v2 TDD

#### Test Suite: `internal/plugin/api_v2_test.go`

**Test Categories:**

1. **API Interface Tests**
```go
func TestPluginAPIv2_New_Initialization(t *testing.T)
func TestPluginAPIv2_RegisterPlugin(t *testing.T)
func TestPluginAPIv2_RegisterHook(t *testing.T)
func TestPluginAPIv2_EmitEvent(t *testing.T)
```

2. **Test Plugin Development Tests**
```go
func TestPluginAPIv2_TestPluginConfiguration(t *testing.T)
func TestPluginAPIv2_TestPluginFunctionality(t *testing.T)
func TestPluginAPIv2_RemoveTviewDependency(t *testing.T)
```

2. **Window Management Tests**
```go
func TestPluginAPIv2_CanShowModal_AllowedEvents(t *testing.T)
func TestPluginAPIv2_CanShowModal_BlockedEvents(t *testing.T)
func TestPluginAPIv2_ShowModal_Success(t *testing.T)
func TestPluginAPIv2_ShowModal_Validation(t *testing.T)
```

3. **Information Display Tests**
```go
func TestPluginAPIv2_UpdateStatus_Success(t *testing.T)
func TestPluginAPIv2_UpdateStatus_LengthLimit(t *testing.T)
func TestPluginAPIv2_ShowNotification_Success(t *testing.T)
func TestPluginAPIv2_ShowNotification_DurationLimit(t *testing.T)
```

### 3.2 Test Plugin Development TDD

#### Test Suite: `internal/plugin/test_plugins_test.go`

**Test Categories:**

1. **Test Plugin Configuration Tests**
```go
func TestTestPlugin_Configuration_Success(t *testing.T)
func TestTestPlugin_Configuration_Validation(t *testing.T)
func TestTestPlugin_Configuration_ErrorHandling(t *testing.T)
```

2. **Test Plugin Functionality Tests**
```go
func TestTestPlugin_Functionality_Validation(t *testing.T)
func TestTestPlugin_Functionality_EventHooks(t *testing.T)
func TestTestPlugin_Functionality_Configuration(t *testing.T)
```

## Phase 4: Enhanced Plugin Information System

### 4.1 Information Display TDD

#### Test Suite: `internal/tui/information_display_test.go`

**Test Categories:**

1. **Status Panel Tests**
```go
func TestStatusPanel_Update_Success(t *testing.T)
func TestStatusPanel_Update_LengthLimit(t *testing.T)
func TestStatusPanel_Update_RealTime(t *testing.T)
func TestStatusPanel_Clear_Success(t *testing.T)
```

2. **Notification System Tests**
```go
func TestNotificationSystem_Show_Success(t *testing.T)
func TestNotificationSystem_AutoDismiss(t *testing.T)
func TestNotificationSystem_MultipleNotifications(t *testing.T)
func TestNotificationSystem_DurationLimits(t *testing.T)
```

### 4.2 Plugin SDK TDD

#### Test Suite: `internal/plugin/sdk_test.go`

**Test Categories:**

1. **SDK Helper Tests**
```go
func TestSDKHelpers_Styling(t *testing.T)
func TestSDKHelpers_Layout(t *testing.T)
func TestSDKHelpers_EventHandling(t *testing.T)
func TestSDKHelpers_ErrorHandling(t *testing.T)
```

## Phase 5: Advanced Plugin Windows & Finalization

### 5.1 Advanced Plugin Windows TDD

#### Test Suite: `internal/plugin/window_manager_test.go`

**Test Categories:**

1. **Window Management Tests**
```go
func TestWindowManager_CreateWindow_Success(t *testing.T)
func TestWindowManager_ShowWindow_Validation(t *testing.T)
func TestWindowManager_HideWindow_Success(t *testing.T)
func TestWindowManager_ZOrderManagement(t *testing.T)
```

2. **Event Integration Tests**
```go
func TestWindowManager_EventIntegration_AllowedEvents(t *testing.T)
func TestWindowManager_EventIntegration_BlockedEvents(t *testing.T)
func TestWindowManager_EventIntegration_UserCancellation(t *testing.T)
```

## End-to-End Testing Strategy

### User Workflow Tests

#### Test Suite: `tests/e2e/user_workflows_test.go`

**Test Categories:**

1. **Basic Timer Workflows**
```go
func TestUserWorkflow_StartTimer_Success(t *testing.T)
func TestUserWorkflow_PauseResume_Success(t *testing.T)
func TestUserWorkflow_StopTimer_Success(t *testing.T)
func TestUserWorkflow_CustomSessionName(t *testing.T)
```

2. **Plugin Integration Workflows**
```go
func TestUserWorkflow_PluginNotification(t *testing.T)
func TestUserWorkflow_PluginModal(t *testing.T)
func TestUserWorkflow_PluginStatus(t *testing.T)
```

3. **Error Handling Workflows**
```go
func TestUserWorkflow_TimerConflict(t *testing.T)
func TestUserWorkflow_ProcessCrash(t *testing.T)
func TestUserWorkflow_PluginError(t *testing.T)
```

### Performance Tests

#### Test Suite: `tests/performance/performance_test.go`

**Test Categories:**

1. **Timer Accuracy Tests**
```go
func TestPerformance_TimerAccuracy_OneHour(t *testing.T)
func TestPerformance_TimerAccuracy_CrossPlatform(t *testing.T)
func TestPerformance_TimerAccuracy_UnderLoad(t *testing.T)
```

2. **Lock Performance Tests**
```go
func TestPerformance_LockOperations_Speed(t *testing.T)
func TestPerformance_LockOperations_Concurrent(t *testing.T)
func TestPerformance_LockOperations_Recovery(t *testing.T)
```

3. **UI Performance Tests**
```go
func TestPerformance_UIResponsiveness(t *testing.T)
func TestPerformance_MemoryUsage(t *testing.T)
func TestPerformance_StartupTime(t *testing.T)
```

## Test Infrastructure

### Test Utilities

#### Mock Implementations
```go
// internal/timer/mocks/file_system_mock.go
type MockFileSystem struct {
    files map[string][]byte
    locks map[string]bool
}

// internal/timer/mocks/process_mock.go
type MockProcessManager struct {
    processes map[int]bool
}

// internal/plugin/mocks/plugin_mock.go
type MockPlugin struct {
    name    string
    hooks   map[string]func(interface{}) error
}
```

#### Test Helpers
```go
// internal/timer/test_helpers.go
func CreateTestTimer() *Timer
func CreateTestLockManager() *TimerLockManager
func SimulateProcessCrash(pid int) error

// internal/tui/test_helpers.go
func CreateTestStage() *GlobalStage
func SimulateKeyPress(key string) tea.KeyMsg
func SimulateWindowResize(width, height int) tea.WindowSizeMsg
```

### Test Configuration

#### Test Environment Setup
```go
// tests/setup.go
func SetupTestEnvironment() error
func TeardownTestEnvironment() error
func CreateTestConfig() *config.Config
func CreateTestLogger() *logger.Logger
```

#### Cross-Platform Testing
```go
// tests/cross_platform_test.go
func TestCrossPlatform_LockBehavior(t *testing.T)
func TestCrossPlatform_ProcessValidation(t *testing.T)
func TestCrossPlatform_FileSystem(t *testing.T)
```

## Continuous Integration

### Test Automation

#### GitHub Actions Workflow
```yaml
# .github/workflows/test.yml
name: TDD Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: [1.21, 1.22]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}
      - run: make test
      - run: make test-coverage
      - run: make test-e2e
```

#### Coverage Reporting
```yaml
# .github/workflows/coverage.yml
name: Coverage Report
on: [push]
jobs:
  coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
      - run: make test-coverage
      - uses: codecov/codecov-action@v3
```

### Quality Gates

#### Pre-commit Hooks
```bash
#!/bin/bash
# .githooks/pre-commit
make test
make lint
make format-check
```

#### Pull Request Requirements
- All tests must pass
- Coverage must meet requirements (80% overall, 95% critical paths)
- No linting errors
- Code must be properly formatted
- Documentation must be updated

## Test Data Management

### Test Fixtures

#### Lock File Fixtures
```json
// tests/fixtures/lock_files/valid_lock.json
{
  "pid": 1234,
  "session_name": "work",
  "start_time": "2025-01-27T10:00:00Z",
  "duration_seconds": 1500,
  "locked_at": "2025-01-27T10:00:00Z"
}
```

#### Timer State Fixtures
```json
// tests/fixtures/timer_states/running_session.json
{
  "status": "running",
  "session_name": "deep work",
  "duration": "45m0s",
  "start_time": "2025-01-27T10:00:00Z",
  "elapsed": "23m15s"
}
```

### Test Utilities

#### Test Data Generators
```go
// tests/utils/generators.go
func GenerateTestSessionName() string
func GenerateTestDuration() time.Duration
func GenerateTestLockState() *LockFileState
func GenerateTestTimerState() *State
```

## Success Metrics

### Test Quality Metrics
- **Test Coverage**: 80% overall, 95% critical paths
- **Test Execution Time**: < 30 seconds for unit tests, < 5 minutes for e2e tests
- **Test Reliability**: 99% pass rate in CI
- **Test Maintainability**: Clear test structure and documentation

### Code Quality Metrics
- **TDD Compliance**: 100% of new code written with TDD
- **Test First Development**: All features start with failing tests
- **Refactoring Safety**: All refactoring maintains test coverage
- **Documentation**: All test suites properly documented

---

**Note**: This TDD plan ensures that Release 0.6.0 is developed with high quality and reliability. All implementation must follow the TDD process and meet the coverage requirements defined in this plan. 