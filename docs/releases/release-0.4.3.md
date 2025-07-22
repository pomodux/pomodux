# Release 0.4.3 - Config Flag Fix and Path Expansion

**Release Date**: 2025-07-21  
**Status**: Planning  
**Type**: Bug Fix Release

## Overview

Release 0.4.3 addresses critical configuration management issues and improves path handling throughout the application. This release fixes the `--config` flag bug and resolves path expansion problems that were causing unwanted directory creation.

## Key Features

### 🔧 **Config Flag Fix**
- **Issue**: The `--config` flag was not properly respected by the plugin system
- **Solution**: Implemented config injection pattern to ensure consistent configuration usage
- **Impact**: Plugin system now correctly uses the specified configuration file

### 🛠️ **Path Expansion Improvements**
- **Issue**: Tilde (`~`) in configuration paths was not being expanded
- **Solution**: Added comprehensive path expansion for configuration paths
- **Impact**: Eliminates unwanted "~" folder creation in repository

### 🧪 **Enhanced Testing**
- **New Tests**: Added path expansion validation tests
- **Coverage**: Improved test coverage for configuration management
- **Validation**: All CI/CD tests pass with 100% success rate

## Technical Changes

### Architecture Improvements

#### 1. Config Injection Pattern
- **File**: `internal/timer/manager.go`
- **Change**: Added `SetGlobalConfig()` function and global config storage
- **Benefit**: Ensures consistent configuration usage across all components

#### 2. Path Expansion System
- **File**: `internal/config/config.go`
- **Change**: Added `expandPath()` function for tilde and environment variable expansion
- **Benefit**: Proper path resolution and elimination of literal path issues

#### 3. Plugin System Enhancement
- **File**: `internal/plugin/manager.go`
- **Change**: Updated `registerGetConfigFn` to use actual configuration values
- **Benefit**: Plugins can now access correct configuration settings

### Code Quality

#### Test Coverage
- **New Tests**: `TestExpandPath()` for path expansion validation
- **Coverage**: All existing tests continue to pass
- **Validation**: Comprehensive test suite for configuration scenarios

#### Code Standards
- **Linting**: All code passes golangci-lint checks
- **Security**: No security vulnerabilities identified
- **Documentation**: Updated inline documentation for new functions

## Bug Fixes

### Fixed Issues

1. **Config Flag Bug** (#config-flag-issue)
   - **Problem**: `--config` flag loaded correct file but plugin system used default config
   - **Root Cause**: Timer manager was calling `config.Load()` internally
   - **Solution**: Implemented config injection pattern
   - **Verification**: Plugin system now respects `--config` flag

2. **Path Expansion Bug** (#path-expansion-issue)
   - **Problem**: Literal "~" folder created in repository
   - **Root Cause**: Tilde not expanded in configuration paths
   - **Solution**: Added comprehensive path expansion
   - **Verification**: No more unwanted directory creation

3. **Plugin Configuration Access** (#plugin-config-issue)
   - **Problem**: Plugins couldn't access actual configuration values
   - **Root Cause**: `get_config` function was not implemented
   - **Solution**: Implemented proper config access for plugins
   - **Verification**: Plugins can now access `plugins.kimai` and other settings

## Testing

### Test Results
- **Unit Tests**: ✅ All passing
- **Integration Tests**: ✅ All passing
- **UAT Tests**: ✅ 27/27 tests passing (100% success rate)
- **CI/CD Pipeline**: ✅ All jobs successful

### Test Coverage
- **Overall Coverage**: Maintained at high levels
- **New Code Coverage**: 100% for path expansion functions
- **Regression Testing**: All existing functionality verified

## Documentation Updates

### Updated Documentation
- **Release Notes**: This document
- **Code Comments**: Enhanced inline documentation
- **Test Documentation**: Added test descriptions for new functionality

### New Documentation
- **Path Expansion Guide**: Documentation for configuration path handling
- **Config Injection Pattern**: Architecture documentation for the new pattern

## Release Criteria

### ✅ **Approval Gate 0: Architecture Review**
- [x] ADR audit completed
- [x] Architecture proposal documented
- [x] Integration approach validated
- [x] Risk assessment completed

### ✅ **Approval Gate 1: Release Plan Approval**
- [x] Feature list defined
- [x] Technical design documented
- [x] TDD approach followed
- [x] Success criteria established

### ✅ **Approval Gate 2: Development Completion**
- [x] All features implemented
- [x] Test coverage requirements met
- [x] Documentation updated
- [x] Code quality standards met

### ✅ **Approval Gate 3: User Acceptance**
- [x] UAT tests passing
- [x] User testing completed
- [x] Bug fixes implemented
- [x] Performance verified

### ⏳ **Approval Gate 4: Release Approval**
- [ ] Final documentation review
- [ ] Release artifacts prepared
- [ ] Stakeholder approval obtained
- [ ] Release deployment ready

## Installation and Usage

### Prerequisites
- Go 1.24.4 or later
- Standard Unix/Linux environment

### Installation
```bash
# Build from source
make build

# Or download pre-built binary
# (will be available after release)
```

### Configuration
The application now properly handles configuration paths with tilde expansion:

```yaml
plugins:
  directory: ~/.config/pomodux/plugins  # Will expand to /home/user/.config/pomodux/plugins
  enabled:
    kimai: true
```

### Usage Examples

#### Using Custom Config File
```bash
# Load specific configuration file
bin/pomodux --config ~/.config/pomodux/config-test.yaml start 25m
```

#### Plugin System
```bash
# Plugin system now respects the config flag
bin/pomodux --config ~/.config/pomodux/config-test.yaml start 2s
```

## Migration Guide

### From Previous Versions
- **No Breaking Changes**: This release is fully backward compatible
- **Automatic Migration**: Existing configurations continue to work
- **Path Expansion**: Existing paths with tilde will be automatically expanded

### Configuration Updates
- **Optional**: Update configuration files to use tilde expansion
- **Recommended**: Use absolute paths or tilde expansion for clarity
- **Validation**: All existing configuration validation rules still apply

## Known Issues

### None Identified
- All known issues from previous releases have been resolved
- No new issues identified during testing

## Future Considerations

### Potential Enhancements
- **Environment Variable Support**: Enhanced environment variable expansion
- **Configuration Templates**: Pre-built configuration templates
- **Validation Improvements**: Enhanced configuration validation feedback

### Technical Debt
- **None**: This release reduces technical debt by fixing architectural issues
- **Code Quality**: Improved code organization and maintainability

## Release Notes Summary

### 🎉 **What's New**
- Fixed `--config` flag to work correctly with plugin system
- Added proper path expansion for configuration files
- Enhanced plugin configuration access
- Improved test coverage and validation

### 🔧 **Bug Fixes**
- Fixed config flag not being respected by plugin system
- Fixed literal "~" folder creation in repository
- Fixed plugin configuration access issues

### 📚 **Documentation**
- Updated release documentation
- Enhanced code comments
- Added architecture documentation

### 🧪 **Testing**
- Added comprehensive path expansion tests
- All CI/CD tests passing
- Improved test coverage

---

**Next Release**: 0.4.4 (planned for future enhancements)  
**Maintenance**: This release will be supported until 0.5.0 