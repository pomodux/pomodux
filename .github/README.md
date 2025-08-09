# GitHub Issues Setup for Pomodux

This directory contains the GitHub Issues configuration for the Pomodux project, designed to support the 4-gate approval process and continuous development approach.

## Issue Templates

### 1. Feature Request (`feature-request.md`)
Use this template for new feature requests that align with the Pomodux roadmap.

**When to use:**
- Requesting new functionality
- Suggesting improvements
- Proposing new commands or features

**Key sections:**
- Feature alignment with project architecture
- User story format
- Acceptance criteria
- Technical requirements (TDD, test coverage)

### 2. Bug Report (`bug-report.md`)
Use this template for reporting bugs or issues in Pomodux.

**When to use:**
- Reporting unexpected behavior
- Documenting crashes or errors
- Reporting performance issues

**Key sections:**
- Environment details (OS, Go version, etc.)
- Steps to reproduce
- Error details and stack traces
- Configuration and state information

### 3. Release Task (`release-task.md`)
Use this template for tracking specific tasks within a release.

**When to use:**
- Implementing planned features
- Technical implementation tasks
- Release-specific deliverables

**Key sections:**
- Release and component information
- TDD approach with specific test cases
- Implementation plan
- Technical considerations

### 4. UAT Testing (`uat-testing.md`)
Use this template for tracking User Acceptance Testing.

**When to use:**
- Gate 3 testing sessions
- User feedback collection
- Feature validation

**Key sections:**
- Test scenarios and steps
- Test results and issues found
- User feedback
- Acceptance decisions

## Issue Workflow

### 1. Issue Creation
1. Choose the appropriate template
2. Fill in all required sections
3. Add appropriate labels
4. Add appropriate component and priority labels

### 2. Issue Triage
- All new issues start with `needs-triage` label
- Review and categorize the issue
- Assign priority and component
- Remove `needs-triage` label

### 3. Development Process
- Issues follow TDD approach
- Tests must be written first
- Code review required
- Test coverage requirements must be met

### 4. Approval Gates
- **Gate 1**: Release plan approval
- **Gate 2**: Development completion
- **Gate 3**: User acceptance testing
- **Gate 4**: Release approval

## Labels

### Issue Types
- `bug` - Bug reports
- `enhancement` - Feature requests
- `release-task` - Release-specific tasks
- `uat` - User acceptance testing

### Priority
- `high` - High priority
- `medium` - Medium priority
- `low` - Low priority

### Status
- `needs-triage` - New issues awaiting review
- `in-progress` - Currently being worked on
- `blocked` - Blocked by other issues
- `ready-for-review` - Ready for code review
- `ready-for-uat` - Ready for user testing

### Components
- `timer-engine` - Core timer functionality
- `cli` - Command-line interface
- `tui` - Terminal user interface
- `plugin-system` - Plugin and extension system
- `configuration` - Configuration management
- `testing` - Testing and quality assurance

## Project Architecture

**Core Components:**
- **Timer Engine**: State management, persistence, and session tracking
- **CLI Interface**: Cobra-based commands with completion and configuration
- **Plugin System**: Lua-based extensibility with event hooks
- **Configuration**: XDG-compliant YAML configuration system
- **TUI System**: Terminal interface with Bubbletea framework

For detailed technical information, see [Architecture Decision Records](../docs/adr/).

## Requirements Management

The main project requirements are maintained in documentation and include:

- **User Stories** with acceptance criteria
- **Technical Tasks** with TDD requirements  
- **Bug fixes** and improvements
- **Non-functional requirements**
- **Architecture decisions** documented in ADRs

### Creating Issues from Architecture
1. Reference relevant Architecture Decision Records
2. Use the appropriate template
3. Include technical context and constraints
4. Add component-specific implementation details

### Architecture Alignment
- Ensure issues align with approved ADRs
- Update documentation when decisions change
- Add new architectural decisions as needed
- Prioritize based on technical debt and user needs

## Best Practices

### Writing Good Issues
- Use clear, descriptive titles
- Include user stories for features
- Provide specific acceptance criteria
- Include technical requirements
- Reference related documentation

### Issue Management
- Keep issues focused and single-purpose
- Use labels consistently
- Update issue status regularly
- Link related issues
- Close issues when complete

### Collaboration
- Use issue comments for discussion
- Reference issues in commit messages
- Link issues to pull requests
- Update issues with progress

## Configuration

The GitHub configuration is set up to:
- Require issue templates (no blank issues)
- Provide links to documentation  
- Support the 4-gate approval process
- Align with continuous development

## Resources

- [Pomodux Documentation](../docs/)
- [Architecture Decision Records](../docs/adr/)
- [Development Setup](../docs/development-setup.md)
- [Go Development Standards](../docs/go-standards.md) 