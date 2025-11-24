# Pomodux

**Pomodux** is a terminal-native timer application that provides a clean, keyboard-driven interface for time management. Whether you're using the Pomodoro Technique or just need a simple countdown timer, Pomodux keeps you focused without leaving your terminal.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Interactive TUI**: Real-time countdown with progress bar using [Bubbletea](https://github.com/charmbracelet/bubbletea)
- **Timer Presets**: Configure common durations (work: 25m, break: 5m, etc.)
- **Session History**: Track completed sessions with statistics
- **Theming**: Built-in themes (default, nord, catppuccin-mocha)
- **Plugin System**: Extend functionality with Lua plugins
- **XDG Compliant**: Follows filesystem standards
- **Crash Recovery**: Resume interrupted sessions
- **Separate Binaries**: `pomodux` for timers, `pomodux-stats` for statistics

## Project Structure

```
pomodux/
├── cmd/
│   ├── pomodux/
│   │   └── main.go           # Entry point for timer binary
│   └── pomodux-stats/
│       └── main.go           # Entry point for stats binary
├── internal/
│   ├── config/
│   │   ├── config.go         # Config struct and loading
│   │   ├── config_test.go    # Config tests
│   │   └── defaults.go       # Default config values
│   ├── timer/
│   │   ├── timer.go          # Timer engine (wall-clock)
│   │   ├── timer_test.go     # Timer tests
│   │   └── state.go          # State persistence
│   ├── tui/
│   │   ├── model.go          # Bubbletea model
│   │   ├── update.go         # Update function
│   │   ├── view.go           # View function
│   │   └── messages.go       # Custom messages
│   ├── history/
│   │   ├── history.go        # Session persistence
│   │   ├── history_test.go   # History tests
│   │   └── query.go          # Query/filter functions
│   ├── theme/
│   │   ├── theme.go          # Theme interface
│   │   ├── themes.go         # Built-in themes
│   │   └── lipgloss.go       # Lipgloss style helpers
│   ├── plugin/
│   │   ├── manager.go        # Plugin manager
│   │   ├── loader.go         # Lua plugin loader
│   │   └── events.go         # Event definitions
│   └── logger/
│       └── logger.go         # Logrus wrapper
├── docs/
│   ├── requirements/
│   │   └── base.md           # Comprehensive requirements
│   ├── adr/                  # Architecture Decision Records
│   │   ├── programming-language-selection.md
│   │   ├── tui-framework-selection.md
│   │   ├── plugin-system-architecture.md
│   │   ├── logger-selection.md
│   │   ├── prettification-library-selection.md
│   │   ├── maintain-separate-binaries.md
│   │   └── concurrency-model.md
│   └── uxdr/                 # UX Design Records
│       ├── terminal-bell-configuration.md
│       ├── keyboard-controls-design.md
│       ├── session-label-defaults.md
│       └── auto-exit-behavior.md
├── go.mod
├── go.sum
├── Makefile
├── LICENSE
└── README.md
```

## Documentation

This project has extensive documentation prepared for implementation:

- **[Requirements](docs/requirements/base.md)**: Comprehensive functional and non-functional requirements with user stories
- **[Architecture Decision Records (ADRs)](docs/adr/)**: Detailed technical decisions and rationales
- **[UX Design Records (UXDRs)](docs/uxdr/)**: User experience design decisions

## Planned Usage

```bash
# Start a 25-minute work session
pomodux start work "Implementing authentication"

# Start a 5-minute break
pomodux start break

# Start a custom duration timer
pomodux start 45m "Client meeting"

# View today's statistics
pomodux-stats --today

# View recent sessions
pomodux-stats --limit 10
```

### Keyboard Controls

| Key       | Action              |
|-----------|---------------------|
| `p`       | Pause timer         |
| `r`       | Resume timer        |
| `s` / `q` | Stop and exit       |
| `Ctrl+C`  | Emergency exit      |


## Configuration

Configuration file location: `~/.config/pomodux/config.yaml`

```yaml
version: "1.0"

timers:
  work: 25m
  break: 5m
  longbreak: 15m
  meeting: 60m

theme: "catppuccin-mocha"

timer:
  bell_on_complete: false

logging:
  level: "info"
  file: ""
```

## Implementation Roadmap

See [docs/requirements/base.md](docs/requirements/base.md) for detailed implementation requirements.

**Planned phases:**
1. Foundation (Config, Logger, Project Structure)
2. Core Timer (Wall-clock engine, State persistence)
3. TUI (Bubbletea, Progress bar, Keyboard controls)
4. History & Statistics (Session persistence, Stats binary)
5. Plugin System (Lua runtime, Event dispatch)
6. Polish & Release (Themes, Documentation, Packaging)

## Technical Decisions

Key architectural decisions:
- **Language**: Go 1.21+ ([ADR-001](docs/adr/programming-language-selection.md))
- **TUI Framework**: Bubbletea + Bubbles + Lipgloss
- **Plugin System**: Lua via gopher-lua
- **Concurrency**: Event loop only, no manual goroutines
- **Timer Accuracy**: Wall-clock calculation using `time.Since()` for drift-free timing
- **State Persistence**: Event-driven + periodic (every 5s) hybrid approach

## License

MIT License - See [LICENSE](LICENSE) file for details.

## Links

- Repository: https://github.com/pomodux/pomodux
- Documentation: See [docs/](docs/) directory
