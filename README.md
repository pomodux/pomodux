# Pomodux

A powerful terminal-based timer and Pomodoro application built in Go, designed for productivity and time management.

## 📋 Features

- **Core Timer Engine**: Robust timer with state management
- **Persistent Timer Sessions**: Interactive keypress controls (p, r, q, s, Ctrl+C)
- **Real-time Progress Display**: Visual progress bars and time remaining
- **Pomodoro Technique**: Dedicated work/break commands
- **Enhanced CLI Interface**: Start, status, history, config, and plugin commands with TUI controls
- **State Persistence**: Timer state survives process restarts
- **Configuration System**: XDG-compliant configuration management
- **Session History**: Track and display session statistics
- **Tab Completion**: Shell completion for all commands
- **Cross-Platform**: Linux, macOS, and Windows support
- **Comprehensive Testing**: 80%+ test coverage with TDD approach

### 🔄 Planned
- **Enhanced CLI Functionality**: Improved status reporting and user experience
- **Plugin System Foundation**: Architecture for extensibility
- **Advanced Notifications**: Enhanced notification system
- **Performance Optimizations**: Improved performance and resource usage
- Plugin system and advanced features
- **TUI**: Terminal User Interface implementation

## 🛠️ Installation

### Prerequisites
- Go 1.21+ (minimum), Go 1.24.4+ recommended
- Git

### Build from Source
```bash
git clone https://github.com/yourusername/pomodux.git
cd pomodux
make build
```

### Install Binary
```bash
# Download the latest release binary for your platform
# Add to your PATH
```

## 🎯 Quick Start

### Basic Timer Usage
```bash
# Start a 25-minute timer
pomodux start 25m

# Check timer status
pomodux status

# Use TUI controls for timer management (p=pause, r=resume, q=stop)
```

### Supported Duration Formats
- `25m` - 25 minutes
- `1h30m` - 1 hour 30 minutes
- `1500s` - 1500 seconds
- `1.5h` - 1.5 hours

## 📁 Project Structure

```
pomodux/
├── cmd/pomodux/          # Main application entry point
├── internal/
│   ├── cli/             # CLI commands and interface
│   ├── config/          # Configuration management
│   └── timer/           # Core timer engine
├── docs/                # Documentation and ADRs
├── .github/             # GitHub templates and workflows
└── Makefile            # Build and development tasks
```

## 🧪 Development

### Prerequisites
- Go 1.21+ (minimum), Go 1.24.4+ recommended
- golangci-lint
- Make

### Development Setup
See [CLAUDE.md](CLAUDE.md) for complete development commands and setup instructions.

### CI/CD Pipeline

Pomodux uses a comprehensive CI/CD pipeline with automated testing, linting, and releases:

- **Continuous Integration**: Runs on every push and pull request
- **Automated Releases**: Triggered by git tags (e.g., `v1.0.0`)
- **Multi-Platform Builds**: Linux, macOS, and Windows (amd64/arm64)
- **Quality Gates**: Automated testing, linting, and security scanning

**Quick Start**:
```bash
# Create a new release
./scripts/release.sh 1.2.3
```

For detailed information, see:
- [CI/CD Pipeline Documentation](docs/ci-cd-pipeline.md) - Complete guide with quick reference

For detailed development commands, testing procedures, and build instructions, see [CLAUDE.md](CLAUDE.md).

## 📚 Documentation

- **[Requirements](docs/requirements.md)** - Project requirements and goals
- **[Technical Specifications](docs/technical_specifications.md)** - Technical architecture and design
- **[Development Setup](docs/development-setup.md)** - Development environment and tools
- **[Go Standards](docs/go-standards.md)** - Go coding standards and conventions
- **[Release Management](docs/release-management.md)** - Release process and approval gates
- **[ADR](docs/adr/)** - Architecture Decision Records

## 🤝 Contributing

### Development Process
Pomodux follows a structured 4-gate approval process:
1. **Gate 1**: Work Plan Approval
2. **Gate 2**: Development Completion
3. **Gate 3**: Testing/Validation
4. **Gate 4**: Final Approval (Releases)

### Issue Management
- Use GitHub issue templates for bug reports and feature requests
- Follow TDD (Test-Driven Development) approach
- Reference requirements and technical specifications for planning
- Link issues to appropriate release milestones

### Code Standards
- Follow Go best practices and standards
- Maintain 80%+ test coverage
- Use golangci-lint for code quality
- Document all public APIs

## 🔧 Configuration

Pomodux uses XDG-compliant configuration:
- **Linux/macOS**: `~/.config/pomodux/config.yaml`
- **Windows**: `%APPDATA%\pomodux\config.yaml`

### Default Configuration
```yaml
timer:
  default_duration: 25m
  auto_start: false
  notifications: true

cli:
  output_format: text
  verbose: false
```

## 🐛 Known Issues

Currently tracking issues in GitHub Issues. See [GitHub Issues](https://github.com/pomodux/pomodux/issues) for active bugs and feature requests.

## 📈 Future Development

- Enhanced CLI functionality and user experience  
- Plugin system development and community extensions
- Terminal User Interface (TUI) implementation
- Advanced notification and integration features
- Custom workflows and automation
- Extended configuration options

For detailed technical plans, see [Architecture Decision Records](docs/adr/).

## 📄 License

[License information to be added]

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- Following [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
- Inspired by the Pomodoro Technique

---

**Note**: Pomodux is actively developed with continuous integration. For the latest updates, check the [issue tracker](https://github.com/yourusername/pomodux/issues) and [Architecture Decision Records](docs/adr/).
