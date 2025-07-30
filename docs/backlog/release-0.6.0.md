# Release 0.6.0 - TUI-Only Timer with File-Based Locking

> **STATUS: PLANNING**
> 
> **Release Type:** Major Feature Release
> **Target Date:** TBD (After comprehensive planning and risk assessment)
> **Priority:** High

---

## Release Overview

**Primary Feature:** TUI-Only Timer with File-Based Locking and Generic Session Support
**Goal:** Refactor the current CLI+TUI dual interface to a unified TUI-only architecture with robust single-timer enforcement and simplified session management.

## Feature Summary

This release implements a major architectural improvement that:
- **Eliminates cross-process synchronization** through single-process TUI architecture
- **Enforces single timer instance** via file-based locking with automatic recovery
- **Simplifies session management** by replacing enum-based session types with generic session names
- **Enhances plugin system** with modern TUI-native API and window management
- **Improves user experience** with immediate visual feedback and clear error messages

## Release Scope

### Core Features
1. **File-Based Lock Manager** - Prevents multiple timer processes with automatic recovery
2. **Generic Timer Architecture** - Session name + duration instead of hardcoded session types
3. **TUI-Only Interface** - Single process architecture eliminating cross-process sync
4. **Enhanced Plugin System** - Modern API designed for TUI architecture
5. **Comprehensive Logging** - Structured logging for all timer operations and events

### Architectural Changes
- Replace SessionType enum with `sessionName string`
- Implement TimerLockManager with XDG-compliant file locking
- Create global stage pattern for TUI component management
- Redesign plugin API for TUI-native integration
- Add 6-event timer system with comprehensive logging

## Release Phases & Milestones

### Phase 1: Core Timer Simplification + Lock Manager (Medium Risk)
**Duration:** 2-3 weeks
**Focus:** Foundation and risk mitigation

#### 1.1 Lock Manager Implementation
- [ ] Implement TimerLockManager with file-based locking
- [ ] Add process validation and orphaned lock recovery
- [ ] Create XDG-compliant lock file location
- [ ] Implement atomic lock operations with error handling
- [ ] Add comprehensive logging for all lock operations
- [ ] **NEW**: Create lock manager tests with mocked file system

#### 1.2 Timer Core Refactoring
- [ ] Replace SessionType enum with `sessionName string`
- [ ] Update timer events to use `session_name` instead of `session_type`
- [ ] Remove `pomodux break` and `pomodux long-break` commands
- [ ] Update `pomodux start` to accept optional session name parameter
- [ ] Integrate lock acquisition into timer startup sequence
- [ ] **NEW**: Update existing timer tests to use session names
- [ ] **NEW**: Clear existing session history data (breaking change)

#### 1.3 Testing & Validation
- [ ] Test concurrent timer prevention across processes
- [ ] Test orphaned lock recovery scenarios
- [ ] Test lock file corruption handling
- [ ] Test process crash scenarios
- [ ] Cross-platform lock behavior validation
- [ ] Performance testing of lock operations

**Success Criteria:**
- Only one timer instance can exist at any time (process-safe)
- Automatic recovery from crashed processes and orphaned locks
- Generic session support with custom session names
- Comprehensive logging for all timer operations

### Phase 2: TUI-Only Refactor + Global Stage (Medium Risk)
**Duration:** 2-3 weeks
**Focus:** User interface and experience

#### 2.1 Global Stage Architecture
- [ ] Create global stage singleton with basic window management
- [ ] Implement timer window as primary component
- [ ] Add plugin modal window support
- [ ] Create plugin notification system
- [ ] Implement plugin status panel integration
- [ ] **NEW**: Migrate existing TUI to use global stage pattern

#### 2.2 TUI Integration
- [ ] Convert `pomodux start` to launch TUI immediately
- [ ] Implement TUI-first, timer-init-second pattern
- [ ] Add session name display in timer window
- [ ] Integrate lock manager with TUI lifecycle
- [ ] Implement responsive layout with Lipgloss positioning
- [ ] **NEW**: Update existing TUI tests to use teatest framework
- [ ] **NEW**: Add session name display tests

#### 2.3 Event System Implementation
- [ ] Implement all 6 timer events with proper timing
- [ ] Create event distribution system from timer to stage
- [ ] Add comprehensive event logging
- [ ] Test event system with existing plugin hooks

**Success Criteria:**
- TUI launches immediately on `pomodux start`
- Session names displayed prominently in timer window
- All existing timer functionality preserved
- Responsive layout adapts to terminal dimensions

### Phase 3: Plugin API Redesign (Medium Risk)
**Duration:** 2-3 weeks
**Focus:** Plugin system modernization

#### 3.1 Plugin API Design
- [ ] Design new plugin API from scratch for TUI architecture
- [ ] Create plugin window spawning rules and validation
- [ ] Implement simple information display system
- [ ] Add auto-dismissing notifications
- [ ] Create plugin SDK with Lipgloss styling helpers
- [ ] **NEW**: Migrate existing tview dialogs to Bubbletea (ADR 007 compliance)

#### 3.2 Plugin System Implementation
- [ ] Implement plugin window management
- [ ] Add plugin information zones in timer panel
- [ ] Create real-time update system via Bubbletea messages
- [ ] Implement plugin development tools
- [ ] Add comprehensive plugin logging
- [ ] **NEW**: Replace tview dependencies with Bubbletea components

#### 3.3 Test Plugin Development
- [ ] Create custom configuration for test plugins
- [ ] Develop test plugins to validate functionality
- [ ] Test plugin window spawning and interaction
- [ ] Validate plugin event system integration
- [ ] **NEW**: Create plugin development documentation

**Success Criteria:**
- Modern plugin API designed for TUI architecture
- Backward compatibility maintained for existing plugins
- Plugin window spawning rules properly enforced
- Comprehensive plugin development experience

### Phase 4: Enhanced Plugin Information System (Low Risk)
**Duration:** 1-2 weeks
**Focus:** Plugin information display

#### 4.1 Information Display Implementation
- [ ] Design and implement information zones in timer panel
- [ ] Create status line updates (max 50 characters)
- [ ] Implement auto-dismissing notifications (3-5 second timeout)
- [ ] Add progress metrics display (key-value pairs)
- [ ] Create real-time update system

#### 4.2 Plugin SDK Development
- [ ] Develop plugin SDK with Lipgloss styling helpers
- [ ] Create plugin development examples
- [ ] Add plugin development tools
- [ ] Implement hot-reload capabilities
- [ ] Create plugin marketplace foundation

**Success Criteria:**
- Simple plugin information display during active timing
- Auto-dismissing notifications work properly
- Plugin SDK provides comprehensive development tools
- Real-time updates via Bubbletea message system

### Phase 5: Advanced Plugin Windows & Finalization (Medium Risk)
**Duration:** 2-3 weeks
**Focus:** Advanced plugin capabilities and polish

#### 5.1 Advanced Plugin Windows
- [ ] Implement complex plugin window support during allowed events
- [ ] Create comprehensive plugin API for window creation
- [ ] Add Lipgloss-based plugin styling system
- [ ] Implement plugin development tools
- [ ] Test all plugin interaction patterns

#### 5.2 Finalization and Polish
- [ ] Comprehensive testing of all plugin interaction patterns
- [ ] Performance optimization and benchmarking
- [ ] Cross-platform compatibility validation
- [ ] Documentation completion
- [ ] User acceptance testing

**Success Criteria:**
- Complex plugin windows work during appropriate timer events
- Plugin API provides comprehensive window management
- All plugin interaction patterns tested and validated
- Performance meets or exceeds current implementation

## Technical Specifications

### Architecture Changes
- **Single Process Architecture**: Eliminate cross-process synchronization
- **File-Based Locking**: XDG-compliant lock files with automatic recovery
- **Global Stage Pattern**: Singleton stage manages all TUI components
- **Event-Driven System**: 6 timer events with comprehensive logging
- **Plugin API Redesign**: Modern API optimized for TUI architecture

### Dependencies
- **Existing**: Bubbletea, Lipgloss, gopher-lua (no new dependencies)
- **Platform Support**: Linux, macOS, Windows
- **File System**: XDG-compliant lock file locations
- **Removed**: tview (replaced with Bubbletea components per ADR 007)

### Performance Requirements
- Timer accuracy equivalent to current implementation
- Lock file operations complete within 1 second
- Responsive UI updates (< 100ms for user interactions)
- Memory usage within reasonable bounds (< 100MB)

## TDD Approach

### Test Strategy
1. **Lock Manager Testing**
   - Unit tests for file locking operations with mocked file system
   - Integration tests for process validation and recovery
   - End-to-end tests for concurrent timer prevention

2. **Timer Core Testing**
   - Unit tests for session name integration
   - Integration tests for event system with mocked plugin manager
   - State persistence tests with session names

3. **TUI Integration Testing**
   - teatest framework for global stage pattern
   - Window management and event distribution tests
   - Lipgloss styling integration tests

4. **Plugin API Testing**
   - Unit tests for new plugin API
   - Integration tests for window spawning rules
   - Backward compatibility tests with existing plugins

### Test Coverage Requirements
- **Overall Coverage**: Minimum 80% for all code
- **Critical Paths**: Minimum 95% for lock manager and timer core
- **Public APIs**: 100% coverage for all public interfaces
- **Integration Tests**: Comprehensive component interaction testing

## Risk Assessment & Mitigation

### Medium Risk Items

#### 1. Plugin API Redesign (Medium)
- **Risk**: Breaking existing plugins and ecosystem
- **Probability**: Medium
- **Impact**: Medium
- **Mitigation**: 
  - Create test plugins to validate functionality
  - Single user environment simplifies testing
  - Phase the changes incrementally
  - Extensive testing with custom test plugins

#### 2. File Locking Implementation (Medium)
- **Risk**: Cross-platform complexity and edge cases
- **Probability**: Medium
- **Impact**: Medium
- **Mitigation**:
  - Use proven file locking patterns
  - Start with Linux implementation
  - Add platforms incrementally
  - Comprehensive cross-platform testing

#### 3. Global Stage Pattern (Medium)
- **Risk**: Architectural complexity and potential conflicts
- **Probability**: Medium
- **Impact**: Medium
- **Mitigation**:
  - Start with minimal viable stage
  - Validate pattern with proof-of-concept
  - Iterate carefully with extensive testing

### Low Risk Items

#### 1. Performance Impact (Low)
- **Risk**: Single process handling all logic may impact performance
- **Probability**: Low
- **Impact**: Low
- **Mitigation**:
  - Establish performance baselines
  - Profile and optimize hot paths
  - Monitor performance throughout development

#### 2. User Learning Curve (Low)
- **Risk**: TUI-only interface may confuse CLI users
- **Probability**: Low
- **Impact**: Low
- **Mitigation**:
  - Clear documentation and examples
  - Single user environment simplifies testing
  - User acceptance testing

### Low Risk Items

#### 1. Cross-Platform Issues (Low)
- **Risk**: Platform-specific implementation issues
- **Probability**: Low
- **Impact**: Low
- **Mitigation**:
  - Use Bubbletea's cross-platform abstractions
  - Test on all target platforms during development

## Success Criteria

### Functional Requirements
- [ ] **Single Timer Enforcement**: Only one timer instance can exist at any time (process-safe)
- [ ] **Robust Lock Management**: Automatic recovery from crashed processes and orphaned locks
- [ ] **Generic Session Support**: Timer accepts any string as session name
- [ ] **Simplified Command Interface**: Single `pomodux start` command replaces multiple session-specific commands
- [ ] **Zero Cross-Process Synchronization**: Single process architecture with file-based state locking
- [ ] **Maintained Functionality**: All existing timer features work identically
- [ ] **Enhanced Plugin Capabilities**: Modern plugin API designed specifically for TUI architecture
- [ ] **Improved User Experience**: Immediate visual feedback with session context and clear error messages

### Non-Functional Requirements
- [ ] **Performance**: Timer accuracy equivalent to current implementation
- [ ] **Cross-Platform**: Works on Linux, macOS, and Windows
- [ ] **Backward Compatibility**: Existing plugins continue to work
- [ ] **Memory Usage**: Within reasonable bounds (< 100MB)
- [ ] **Responsive UI**: Updates within 100ms for user interactions
- [ ] **Lock Operations**: Complete within 1 second
- [ ] **ADR Compliance**: Aligns with all existing architectural decisions

### Quality Requirements
- [ ] **Test Coverage**: 80% overall, 95% critical paths
- [ ] **Code Quality**: Passes all linting and formatting checks
- [ ] **Documentation**: Complete and up-to-date
- [ ] **User Testing**: Comprehensive user acceptance testing
- [ ] **Cross-Platform Testing**: Validated on all supported platforms

## Release Artifacts

### Source Code
- `internal/timer/lock.go` - File-based lock manager
- `internal/timer/timer.go` - Updated timer core with session names
- `internal/tui/stage.go` - Global stage manager
- `internal/plugin/api.go` - New plugin API
- Updated CLI commands and integration

### Documentation
- Updated user documentation with new session naming
- Plugin development guide for new API
- Migration guide for existing plugins
- Technical documentation for lock manager

### Tests
- Comprehensive test suite for lock manager
- TUI integration tests with teatest
- Plugin API tests and examples
- Cross-platform compatibility tests

## Supporting Documentation

### Technical Specifications
- **[Technical Specifications](release-0.6.0-technical-specs.md)** - Detailed implementation approach and architecture changes
- **[TDD Plan](release-0.6.0-tdd-plan.md)** - Comprehensive test-driven development strategy

### Feature Documentation
- **[TUI Timer Feature Specification](tui_timer_feature_spec.md)** - Original feature specification with detailed requirements

## Codebase Audit Summary

### Current Implementation Status
✅ **Already Implemented:**
- Basic Bubbletea TUI with Lipgloss styling (`internal/tui/tui.go`)
- Lua-based plugin system with tview dialogs (`internal/plugin/manager.go`)
- SessionType enum-based timer core (`internal/timer/timer.go`)
- Complete CLI command set (`internal/cli/`)
- Comprehensive configuration system (`internal/config/`)
- Structured logging with logrus (`internal/logger/`)
- Unit tests, integration tests, and UAT (`tests/`)

⚠️ **Key Issues Identified:**
- **Mixed TUI Libraries**: Plugin system uses tview (violates ADR 007)
- **SessionType Architecture**: Hardcoded enum system (needs replacement)
- **Cross-Process Architecture**: Current system has cross-process synchronization
- **No File Locking**: No existing lock manager implementation
- **No Global Stage**: No existing stage management system
- **No teatest Integration**: TUI tests not using teatest framework

### Documentation Requirements
- **ADR for Global Stage**: ✅ **COMPLETED** - ADR 008 created for global stage pattern
- **Plugin Development Guide**: ✅ **COMPLETED** - Updated for Release 0.6.0 API
- **Migration Guide**: ✅ **COMPLETED** - Comprehensive migration guide created
- **Configuration File Specifications**: ✅ **COMPLETED** - Comprehensive configuration documentation created

### Impact on Release Scope
- **Phase 1**: Must create lock manager from scratch, update existing timer tests
- **Phase 2**: Must migrate existing TUI to global stage, add teatest integration
- **Phase 3**: Must migrate tview dialogs to Bubbletea, remove tview dependency
- **Testing**: Must update existing test infrastructure for new components

## Gate 1 Preparation

### Work Plan Approval Requirements
- [x] Work scope is defined and approved
- [x] Technical approach is planned and documented
- [x] TDD approach is planned for all work
- [x] Success criteria are measurable and approved
- [x] Risk assessment and mitigation strategies documented
- [x] Implementation phases and milestones defined

### Approval Status
- **Status**: Ready for Gate 1 Review
- **Next Step**: Stakeholder approval of release plan
- **Dependencies**: None (feature is in backlog)

## Release Notes Preview

### What's New in 0.6.0
- **TUI-Only Interface**: Unified TUI architecture eliminates cross-process synchronization
- **Generic Session Support**: Use any session name instead of predefined types
- **Single Timer Enforcement**: File-based locking prevents multiple timer instances
- **Enhanced Plugin System**: Modern plugin API designed for TUI architecture
- **Improved User Experience**: Immediate visual feedback and clear error messages

### Breaking Changes
- **Session Types**: Enum-based session types replaced with generic session names
- **Plugin API**: Complete redesign of plugin system (existing plugins will not work)
- **Commands**: `pomodux break` and `pomodux long-break` commands removed
- **Session History**: Existing session history will be cleared

### Migration Guide
- **Session Names**: Use `pomodux start 25m "work"` instead of predefined types
- **Plugins**: Create new plugins using the updated API (test plugins provided)
- **Configuration**: No configuration changes required
- **History**: Session history will start fresh with new format

---

**Note**: This release represents a major architectural improvement that simplifies the codebase while expanding capabilities. The phased approach minimizes risk while delivering significant user value. 