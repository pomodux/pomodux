# Base Requirements Specification

### Epic 1: Timer Management

#### US-1.1: Start Timer with Preset
**As a** technical professional  
**I want to** start a timer using a named preset  
**So that** I can quickly begin focused work sessions

**Acceptance Criteria:**
- Given a config file with `timers.work: 25m` defined
- When I run `pomodux start work "Implementing auth"`
- Then a 25-minute timer starts with label "Implementing auth"
- And the TUI displays progress bar and remaining time
- And the session is recorded when completed

**Priority:** P0 (MVP Critical)

---

#### US-1.2: Start Timer with Custom Duration
**As a** technical professional  
**I want to** start a timer with an arbitrary duration  
**So that** I can track non-standard work sessions

**Acceptance Criteria:**
- Given no specific preset required
- When I run `pomodux start 45m "Client meeting"`
- Then a 45-minute timer starts with label "Client meeting"
- And the TUI displays progress bar and remaining time
- And the session is recorded when completed

**Priority:** P0 (MVP Critical)

---

#### US-1.3: Pause and Resume Timer
**As a** technical professional  
**I want to** pause and resume a running timer  
**So that** I can handle interruptions without losing my session

**Acceptance Criteria:**
- Given a running timer at 15:23 remaining
- When I press 'p' key
- Then the timer pauses and displays "PAUSED" status
- And the remaining time freezes at 15:23
- When I press 'r' key
- Then the timer resumes from 15:23
- And paused duration is tracked in session data

**Priority:** P0 (MVP Critical)

---

#### US-1.4: Stop Timer Early
**As a** technical professional  
**I want to** stop a timer before completion  
**So that** I can end a session when interrupted or task is complete

**Acceptance Criteria:**
- Given a running timer at 10:15 remaining
- When I press 's' or 'q' key
- Then the timer stops immediately
- And the session is saved with `end_status: "stopped"`
- And the TUI exits gracefully

**Priority:** P0 (MVP Critical)

---

#### US-1.5: Default Label for Quick Start
**As a** technical professional  
**I want to** omit the label when starting a timer  
**So that** I can start timers quickly for generic tasks

**Acceptance Criteria:**
- Given I want to start a timer without typing a label
- When I run `pomodux start work`
- Then a timer starts with label "Work" (prettified preset name)
- And the session is recorded with this label
- When I run `pomodux start 45m`
- Then a timer starts with label "Generic timer session" (no preset to prettify)
- And the session is recorded with this label

**Priority:** P0 (MVP Critical)

**See also:** [UXDR: Session Label Defaults](../uxdr/session-label-defaults.md) for detailed design rationale

---

### Epic 2: Configuration Management

#### US-2.1: Define Custom Timer Presets
**As a** technical professional  
**I want to** define custom timer presets in a config file  
**So that** I can reuse common durations without typing them

**Acceptance Criteria:**
- Given a config file at `~/.config/pomodux/config.yaml`
- When I define:
  ```yaml
  timers:
    work: 25m
    break: 5m
    meeting: 60m
    review: 15m
  ```
- Then `pomodux start work` uses 25 minutes
- And `pomodux start meeting` uses 60 minutes
- And unrecognized presets show error message

**Priority:** P0 (MVP Critical)

---

#### US-2.2: Select Theme via Configuration
**As a** technical professional  
**I want to** choose a visual theme in my config file  
**So that** the TUI matches my terminal aesthetic

**Acceptance Criteria:**
- Given a config file with `theme: "catppuccin-mocha"`
- When I start `pomodux`
- Then the TUI uses Catppuccin Mocha color scheme
- And progress bar, text, and borders reflect theme colors
- And invalid theme name falls back to default theme

**Priority:** P0 (MVP Critical)

---

#### US-2.3: XDG-Compliant Configuration
**As a** power user  
**I want to** store config files in XDG-compliant directories  
**So that** my system stays organized per standards

**Acceptance Criteria:**
- Given XDG environment variables are set
- When pomodux looks for config
- Then it checks `$XDG_CONFIG_HOME/pomodux/config.yaml`
- Or falls back to `~/.config/pomodux/config.yaml`
- And session history goes to `$XDG_STATE_HOME/pomodux/history.json`
- Or falls back to `~/.local/state/pomodux/history.json`

**Priority:** P0 (MVP Critical)

---

### Epic 3: History & Statistics

#### US-3.1: View Recent Sessions
**As a** technical professional  
**I want to** view my recent timer sessions  
**So that** I can see what I've worked on today

**Acceptance Criteria:**
- Given I have completed 5 sessions today
- When I run `pomodux-stats`
- Then I see a list of sessions with:
  - Start time
  - Duration
  - Label
  - End status (completed/stopped)
- And sessions are sorted newest first

**Priority:** P0 (MVP Critical)

---

#### US-3.2: View Daily Statistics
**As a** technical professional  
**I want to** see summary statistics for today  
**So that** I can track my productivity

**Acceptance Criteria:**
- Given I have completed 3 work sessions and 2 breaks today
- When I run `pomodux-stats --today`
- Then I see:
  - Total time worked
  - Number of completed sessions
  - Number of stopped sessions
  - Completion rate percentage
  
**Priority:** P0 (MVP Critical)

---

#### US-3.3: Filter Sessions by Date Range
**As a** technical professional  
**I want to** filter session history by date  
**So that** I can review weekly or monthly productivity

**Acceptance Criteria:**
- Given sessions exist for the past 2 weeks
- When I run `pomodux-stats --week`
- Then I see only sessions from the last 7 days
- When I run `pomodux-stats --all`
- Then I see all recorded sessions

**Priority:** P1 (Post-MVP)

---

### Epic 4: Theming

#### US-4.1: Built-in Theme Selection
**As a** technical professional  
**I want to** choose from built-in themes  
**So that** I can quickly match my terminal aesthetic

**Acceptance Criteria:**
- Given pomodux ships with 2 themes: `default`, `catppuccin-mocha`
- When I set `theme: "catppuccin-mocha"` in config
- Then the TUI uses Catppuccin Mocha color scheme
- And all UI elements (progress bar, text, borders) are themed
- And invalid theme name falls back to default theme

**Priority:** P0 (MVP Critical)

---

#### US-4.2: Custom Theme Files (Post-MVP)
**As a** power user  
**I want to** create custom theme files  
**So that** I can define my own color schemes

**Acceptance Criteria:**
- Given I want to create a custom theme
- When I create `~/.config/pomodux/themes/nord.yml` with theme definition
- And I set `theme: "nord"` in config
- Then the TUI uses my custom Nord color scheme
- And all UI elements apply the custom theme

**Priority:** P1 (Post-MVP)

---

### Epic 5: Plugin-Ready Architecture

#### US-5.1: Event Emission (Internal)
**As a** future plugin developer  
**I want to** the application to emit lifecycle events internally  
**So that** plugins can hook into these events later

**Acceptance Criteria:**
- Given the application is running
- When timer starts, the system emits `TimerStarted` event with session data
- When timer completes, the system emits `TimerCompleted` event
- When timer is paused, the system emits `TimerPaused` event
- When timer is resumed, the system emits `TimerResumed` event
- When timer is stopped early, the system emits `TimerStopped` event
- When app starts, the system emits `ApplicationStarted` event
- When app exits normally, the system emits `ApplicationStopping` event
- When app receives SIGINT/SIGTERM/SIGHUP, the system emits `ApplicationInterrupted` event
- When config loads, the system emits `ConfigurationLoaded` event
- And all events are logged via logrus but not dispatched to plugins in MVP

**Priority:** P0 (MVP Critical - Architecture Only)

---

#### US-5.2: Plugin System Implementation
**As a** technical professional  
**I want to** load Lua plugins that react to timer events  
**So that** I can extend pomodux with custom integrations

**Acceptance Criteria:**
- Given a plugin file at `~/.config/pomodux/plugins/notify/init.lua`
- When pomodux starts
- Then the plugin is loaded and registered
- When a timer completes
- Then the plugin's `on_timer_completed` function is called
- And the plugin can send notifications or perform actions
- When app is interrupted
- Then the plugin's `on_app_interrupted` function is called
- And the plugin can perform cleanup

**Priority:** P1 (Post-MVP - See ADR 004)

---

## Functional Requirements

### Timer Management (FR-TIMER)

#### FR-TIMER-001: Timer Start
**Description:** User must be able to start a timer with specified duration or preset.

**Requirements:**
- Accept duration in format: `25m`, `1h`, `90s`, `1h30m`
- Accept preset name from config file
- Require label (with default: "Generic timer session")
- Validate duration is positive and reasonable (max 24 hours)
- Create session record on start
- Display TUI with progress visualization

**Dependencies:** FR-CONFIG-001 (Config Loading)

---

#### FR-TIMER-002: Timer Pause/Resume
**Description:** User must be able to pause and resume running timer.

**Requirements:**
- Pause timer on 'p' keypress
- Resume timer on 'r' keypress
- Track total paused duration
- Track number of pauses
- Display "PAUSED" status in TUI
- Preserve timer state if application exits while paused

**Dependencies:** FR-TIMER-001

---

#### FR-TIMER-003: Timer Stop
**Description:** User must be able to stop timer before completion.

**Requirements:**
- Stop timer on 's' or 'q' keypress
- Save session with `end_status: "stopped"`
- Record actual duration (not full preset duration)
- Exit TUI gracefully
- Clean up state files

**Dependencies:** FR-TIMER-001, FR-HISTORY-001

---

#### FR-TIMER-004: Timer Completion
**Description:** Timer must detect completion and notify user.

**Requirements:**
- Detect when timer reaches 0:00
- Save session with `end_status: "completed"`
- Display completion message briefly
- Emit `TimerCompleted` event (internal)
- Optional: Ring terminal bell (configurable, see [UXDR: Terminal Bell](../uxdr/terminal-bell-configuration.md))
- Exit immediately and return to command line (see [UXDR: Auto-Exit](../uxdr/auto-exit-behavior.md))

**Dependencies:** FR-TIMER-001, FR-HISTORY-001

---

### Configuration Management (FR-CONFIG)

#### FR-CONFIG-001: Config File Loading
**Description:** Application must load configuration from XDG-compliant location.

**Requirements:**
- Check `$XDG_CONFIG_HOME/pomodux/config.yaml`
- Fallback to `~/.config/pomodux/config.yaml`
- Create default config if none exists
- Validate config schema on load
- Log errors for invalid config
- Continue with defaults if config is invalid

**Dependencies:** None

---

#### FR-CONFIG-002: Timer Presets
**Description:** Users must be able to define reusable timer durations.

**Requirements:**
- Support YAML map: `timers: { work: 25m, break: 5m }`
- Validate duration format on load
- Allow arbitrary preset names
- Provide sensible defaults (work: 25m, break: 5m, longbreak: 15m)
- Error message if user references undefined preset

**Dependencies:** FR-CONFIG-001

---

#### FR-CONFIG-003: Theme Configuration
**Description:** Users must be able to select visual theme.

**Requirements:**
- Support `theme: "name"` in config
- Built-in themes: `default`, `nord`, `catppuccin-mocha`
- Fallback to `default` if theme not found
- Theme defines: primary color, secondary color, background, text, progress bar style

**Dependencies:** FR-CONFIG-001

---

#### FR-CONFIG-004: Logging Configuration
**Description:** Users must be able to configure debug logging.

**Requirements:**
- Support `logging.level: "debug|info|warn|error"`
- Support `logging.file: "/path/to/log"` (empty = stderr only)
- Default: `level: info`, `file: ""`
- Integrate with logrus (ADR 005)

**Dependencies:** FR-CONFIG-001

---

### History & Statistics (FR-HISTORY)

#### FR-HISTORY-001: Session Persistence
**Description:** Completed sessions must be saved to persistent storage.

**Requirements:**
- Save to `~/.local/state/pomodux/history.json`
- JSON array of session objects (see section 6.2)
- Append new sessions to array
- Create file if doesn't exist
- Handle file corruption gracefully (backup + recreate)

**Dependencies:** None

---

#### FR-HISTORY-002: Session Listing
**Description:** Users must be able to view recent sessions.

**Requirements:**
- `pomodux-stats` displays tabular list of sessions
- Columns: Start Time, Duration, Label, Status
- Default: Show last 20 sessions
- Sort: Newest first
- Support `--limit N` flag
- Format timestamps in local timezone

**Dependencies:** FR-HISTORY-001

---

#### FR-HISTORY-003: Daily Statistics
**Description:** Users must be able to view daily summary statistics.

**Requirements:**
- `pomodux-stats --today` calculates:
  - Total time (sum of all session durations)
  - Completed sessions count
  - Stopped sessions count
  - Completion rate percentage
- Filter sessions by current day (00:00 to 23:59 local time)
- Display in human-readable format

**Dependencies:** FR-HISTORY-001

---

#### FR-HISTORY-004: Date Range Filtering (Post-MVP)
**Description:** Users must be able to filter sessions by date range.

**Requirements:**
- `pomodux-stats --week` shows last 7 days
- `pomodux-stats --month` shows last 30 days
- `pomodux-stats --all` shows all sessions
- Support custom range: `--from YYYY-MM-DD --to YYYY-MM-DD`

**Dependencies:** FR-HISTORY-001  
**Priority:** P1 (Post-MVP)

---

### Plugin Architecture (FR-PLUGIN)

#### FR-PLUGIN-001: Event System (Internal - MVP)
**Description:** Application must emit internal events for future plugin system.

**Requirements:**

**Timer Events:**
- `TimerStarted`: Timer begins countdown
- `TimerPaused`: Timer paused by user
- `TimerResumed`: Timer resumed from pause
- `TimerStopped`: Timer stopped before completion (user pressed 's' or 'q')
- `TimerCompleted`: Timer reached 0:00

**Application Lifecycle Events:**
- `ApplicationStarted`: Application initialization complete, config loaded, ready to start timer
- `ApplicationStopping`: Normal shutdown initiated (timer completed or user quit)
- `ApplicationInterrupted`: Emergency shutdown (Ctrl+C, SIGTERM, SIGHUP)
- `ConfigurationLoaded`: Configuration file loaded and validated (fires on every app start)

**Event Data Structure:**
- Each event includes: event type, timestamp, relevant context data
- Timer events include: session ID, label, preset, duration, remaining time
- Application events include: version, config path, signal type (for interrupts)
- All events use consistent JSON-serializable structure

**Event Handling (MVP):**
- Events logged via logrus with structured fields
- Events do not block application execution
- Event emission is synchronous (logged immediately)
- Events fire even if no plugins loaded (architecture preparation)
- All signals (SIGINT, SIGTERM, SIGHUP) emit single `ApplicationInterrupted` event with signal info in data

**Future (Post-MVP - See ADR 004):**
- Asynchronous dispatch to loaded plugins
- Plugin hooks registered per event type
- Plugin execution in isolated Lua environment
- Plugin errors don't crash application

**Signal Handling:**
- Graceful handling of SIGINT (Ctrl+C), SIGTERM (kill), SIGHUP (terminal close)
- `ApplicationInterrupted` event emitted before cleanup
- Timer state saved before exit
- SIGKILL cannot be caught (timer state may not be saved)

**Dependencies:** FR-TIMER-001, FR-CONFIG-001  
**Priority:** P0 (Architecture Only)

---

#### FR-PLUGIN-002: Lua Plugin Loading (Post-MVP)
**Description:** Application must load and execute Lua plugins.

**Requirements:**
- Load plugins from `~/.config/pomodux/plugins/`
- Use `gopher-lua` runtime (ADR 004)
- Each plugin runs in isolated Lua state
- Plugins register hooks: `on_timer_started()`, `on_timer_completed()`, etc.
- Plugin errors logged but don't crash application
- Plugins can be enabled/disabled via config

**Dependencies:** FR-PLUGIN-001  
**Priority:** P1 (Post-MVP - See ADR 004)

---

## Non-Functional Requirements

### Performance (NFR-PERF)

#### NFR-PERF-001: Startup Time
**Requirement:** Application must start within 2 seconds on modern hardware.

**Measurement:**
- Measured from process start to TUI render
- Hardware baseline: 2020+ laptop (4-core CPU, 8GB RAM, SSD)
- Includes: config loading, state restoration, TUI initialization

**Acceptance:** `time pomodux start work` completes in <2s

---

#### NFR-PERF-002: Memory Usage
**Requirement:** Application must use less than 50MB of RAM during normal operation.

**Measurement:**
- Measured with single running timer
- Includes: TUI, timer engine, logging
- Excludes: Go runtime overhead (counted separately)

**Acceptance:** `ps aux | grep pomodux` shows <50MB RSS

---

#### NFR-PERF-003: Timer Accuracy
**Requirement:** Timer calculation must be exact based on wall-clock time.

**Implementation Approach:**
- Use `time.Since(startTime)` for calculation, not tick-based countdown
- Calculate remaining time on every render from absolute start time
- Display always shows actual remaining time, no accumulated drift
- Pause duration tracked separately and subtracted from elapsed time

**Measurement:**
- Timer completion occurs at exact configured duration (within system clock precision)
- Display updates at 250ms intervals for smooth countdown visualization
- Accuracy limited only by Go's `time.Time` precision (nanoseconds) and OS clock accuracy
- Test on multiple platforms (Linux, macOS, Windows)

**Acceptance:** Timer completes within 100ms of configured duration (system clock precision)

---

### 4.2 Reliability (NFR-REL)

#### NFR-REL-001: Crash Recovery
**Requirement:** Application must gracefully recover from unexpected termination.

**Requirements:**
- Timer state persisted using hybrid strategy:
  - Event-driven: On start, pause, resume, stop, interrupt
  - Periodic backup: Every 5 seconds while running (not while paused)
- On restart, detect interrupted session
- Offer to resume from last saved state
- Mark unrecovered sessions as `end_status: "interrupted"`
- Handle SIGINT (Ctrl+C), SIGTERM (kill), and SIGHUP (terminal close) gracefully
- Emit `ApplicationInterrupted` event before cleanup
- SIGKILL (kill -9) cannot be caught - maximum 5 seconds data loss

**Recovery Accuracy:**
- Paused timers: Exact state recovery (saved on pause event)
- Running timers: Recovery within 5 seconds of crash time
- Interrupted timers: Emergency save captures final state

**Testing:** Kill process with various signals, verify recovery accuracy on restart

---

#### NFR-REL-002: Data Integrity
**Requirement:** Session history must not be corrupted by crashes.

**Requirements:**
- Use atomic file writes (write to temp file, then rename)
- Validate JSON on load
- Backup corrupt files to `.backup` suffix
- Continue operation with empty history if unrecoverable

**Testing:** Corrupt history.json, verify graceful degradation

---

#### NFR-REL-003: Config Validation
**Requirement:** Config and theme load failure must prevent the timer from loading; no TUI is shown.

**Requirements:**
- Validate config schema on load
- Log errors for invalid fields
- Use default values for invalid/missing fields where applicable
- On config load failure or theme resolution failure (e.g. unreadable config, unknown theme name): do not start TUI; timer does not load; return error to CLI and exit
- No in-TUI config or theme error banner
- Emit `ConfigurationLoaded` event even with partial config when load succeeds

**Testing:** Provide malformed YAML or unknown theme name; verify app does not start TUI and returns error

---

### Usability (NFR-USE)

#### NFR-USE-001: Learning Curve
**Requirement:** New users must be able to start their first timer within 5 minutes.

**Requirements:**
- Clear help text: `pomodux --help`
- Intuitive commands: `pomodux start 25m "My task"`
- Control legend visible in TUI at all times
- Error messages suggest correct usage

**Testing:** User testing with CLI newcomers

---

#### NFR-USE-002: Keyboard Accessibility
**Requirement:** All functionality must be accessible via keyboard only.

**Requirements:**
- No mouse required for any operation
- Keyboard shortcuts displayed in TUI
- Single-key actions (no chords except Ctrl+C)
- Works with screen readers (basic support)

**Testing:** Operate entirely via keyboard, verify all functions accessible

---

#### NFR-USE-003: Error Messages
**Requirement:** Error messages must be clear and actionable.

**Requirements:**
- Specify what went wrong
- Suggest how to fix it
- Include relevant context (file paths, values, etc.)
- No technical jargon unless necessary

---

### Maintainability (NFR-MAINT)

#### NFR-MAINT-001: Code Quality
**Requirement:** Code must follow Go best practices and idioms.

**Requirements:**
- Pass `golangci-lint` with default settings
- Pass `go vet` with no warnings
- Formatted with `gofmt`
- Clear package boundaries (`internal/timer`, `internal/config`, etc.)

**Validation:** CI pipeline enforces all checks

---

#### NFR-MAINT-002: Documentation
**Requirement:** Code must be well-documented for contributors.

**Requirements:**
- Package-level godoc for all packages
- Function-level godoc for exported functions
- Inline comments for complex logic
- Architecture decision records (ADRs) for major decisions
- README with setup instructions

**Validation:** `go doc` generates complete documentation

---

#### NFR-MAINT-003: Plugin API Stability (Post-MVP)
**Requirement:** Plugin API must remain stable across minor versions.

**Requirements:**
- Semantic versioning for API changes
- Deprecation warnings before removal (1 minor version)
- Plugin API versioning in config
- Clear migration guides for breaking changes

**Priority:** P1 (Post-MVP)

---

###  Portability (NFR-PORT)

#### NFR-PORT-001: Cross-Platform Support
**Requirement:** Application must work on Linux, macOS, and Windows.

**Requirements:**
- Linux: Primary target, test on Ubuntu 22.04+ and Arch Linux
- macOS: Test on macOS 12+ (Intel and Apple Silicon)
- Windows: Test on Windows 10+ (best-effort support)
- No OS-specific dependencies in core features
- Use Go standard library for file operations

**Testing:** CI runs tests on all three platforms

---

#### NFR-PORT-003: XDG Compliance
**Requirement:** Application must follow XDG Base Directory specification.

**Requirements:**
- Config: `$XDG_CONFIG_HOME/pomodux/` (default: `~/.config/pomodux/`)
- State: `$XDG_STATE_HOME/pomodux/` (default: `~/.local/state/pomodux/`)
- Cache: `$XDG_CACHE_HOME/pomodux/` (default: `~/.cache/pomodux/`) [future use]
- No files in home directory root

**Validation:** Check file locations on fresh install

---

### Security (NFR-SEC)

#### NFR-SEC-001: Config File Permissions
**Requirement:** Config files must not expose sensitive data.

**Requirements:**
- Config files created with 0600 permissions (user read/write only)
- History files created with 0600 permissions
- Warn if config file has overly permissive permissions
- No secrets stored in config (future plugin auth uses keyring)

**Testing:** Verify file permissions after creation

---

#### NFR-SEC-003: Input Validation
**Requirement:** All user inputs must be validated.

**Requirements:**
- Timer duration: positive integer, reasonable maximum (24 hours)
- Label: maximum length (200 characters), sanitize special characters
- Config values: type-checked, range-checked
- File paths: prevent path traversal attacks

**Testing:** Fuzzing with invalid inputs

---

## Data Specifications

### Configuration File Schema

**Location:** `~/.config/pomodux/config.yaml`

```yaml
# Pomodux Configuration File
# Version: 1.0.0

# Application version (do not edit manually)
version: "1.0"

# Timer presets
timers:
  work: 25m        # Standard work session
  break: 5m        # Short break
  longbreak: 15m   # Long break after 4 work sessions
  meeting: 60m     # One-hour meeting
  review: 15m      # Code review session

# Visual theme
theme: "catppuccin-mocha"  # Options: default, nord, catppuccin-mocha

# Timer behavior
timer:
  bell_on_complete: false   # Ring terminal bell on completion
  # See: docs/uxdr/terminal-bell-configuration.md

# Logging (for debugging)
logging:
  level: "info"   # Options: debug, info, warn, error
  file: ""        # Empty = stderr only, or path to log file

# Plugin configuration (Post-MVP)
plugins:
  enabled: []     # List of enabled plugin names
  directory: ""   # Custom plugin directory (empty = default)
```

**Validation Rules:**
- `timers`: Map of string to duration (e.g., "25m", "1h30m")
- `theme`: Must be known theme name or "custom"
- `timer.bell_on_complete`: Boolean
- `logging.level`: One of "debug", "info", "warn", "error"
- `logging.file`: Valid file path or empty string

**Default Values:**
```yaml
version: "1.0"
timers:
  work: 25m
  break: 5m
  longbreak: 15m
timer:
  bell_on_complete: false
theme: "default"
logging:
  level: "info"
  file: ""
plugins:
  enabled: []
  directory: ""
```

---

### 5.2 Session History Schema

**Location:** `~/.local/state/pomodux/history.json`

```json
{
  "version": "1.0",
  "sessions": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "started_at": "2025-01-15T14:00:00Z",
      "ended_at": "2025-01-15T14:25:00Z",
      "duration": "25m",
      "preset": "work",
      "label": "Implementing auth module",
      "end_status": "completed",
      "paused_count": 2,
      "paused_duration": "3m"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "started_at": "2025-01-15T14:30:00Z",
      "ended_at": "2025-01-15T14:35:00Z",
      "duration": "5m",
      "preset": "break",
      "label": "Generic timer session",
      "end_status": "completed",
      "paused_count": 0,
      "paused_duration": "0s"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "started_at": "2025-01-15T14:40:00Z",
      "ended_at": "2025-01-15T14:55:00Z",
      "duration": "25m",
      "preset": "work",
      "label": "Code review PR #123",
      "end_status": "stopped",
      "paused_count": 1,
      "paused_duration": "2m"
    }
  ]
}
```

**Field Specifications:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | UUID string | Yes | Unique session identifier |
| `started_at` | ISO 8601 timestamp | Yes | Session start time (UTC) |
| `ended_at` | ISO 8601 timestamp | Yes | Session end time (UTC) |
| `duration` | Duration string | Yes | Configured duration (e.g., "25m") |
| `preset` | String | No | Preset name used (null if custom duration) |
| `label` | String | Yes | Session label/description |
| `end_status` | Enum string | Yes | Session outcome: "completed", "stopped", "cancelled", "interrupted" |
| `paused_count` | Integer | Yes | Number of times paused |
| `paused_duration` | Duration string | Yes | Total time spent paused |

**`end_status` Values:**
- `"completed"`: Timer ran to 0:00 successfully
- `"stopped"`: User stopped timer early (pressed 's' or 'q')
- `"cancelled"`: User cancelled with Ctrl+C during active session
- `"interrupted"`: Application crashed or was killed (SIGKILL)

**File Operations:**
- Append-only for new sessions
- Atomic writes (temp file + rename)
- Periodic compaction (remove old sessions, configurable)
- Backup on corruption: `history.json.backup`

---

### 5.3 Timer State Schema

**Location:** `~/.local/state/pomodux/timer_state.json`

**Purpose:** Persist timer state for crash recovery

```json
{
  "version": "1.0",
  "session_id": "550e8400-e29b-41d4-a716-446655440003",
  "started_at": "2025-01-15T15:00:00Z",
  "duration": "25m",
  "preset": "work",
  "label": "Implementing timer persistence",
  "remaining": "15m23s",
  "is_paused": false,
  "paused_count": 1,
  "paused_duration": "2m",
  "last_updated": "2025-01-15T15:09:37Z"
}
```

**Field Specifications:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `version` | String | Yes | State file format version |
| `session_id` | UUID string | Yes | Session ID (matches history) |
| `started_at` | ISO 8601 timestamp | Yes | Session start time (UTC) |
| `duration` | Duration string | Yes | Configured duration |
| `preset` | String | No | Preset name (null if custom) |
| `label` | String | Yes | Session label |
| `remaining` | Duration string | Yes | Time remaining at last update |
| `is_paused` | Boolean | Yes | Current pause state |
| `paused_count` | Integer | Yes | Number of pauses so far |
| `paused_duration` | Duration string | Yes | Total paused time |
| `last_updated` | ISO 8601 timestamp | Yes | Last state save time |

**File Lifecycle:**
- Created when timer starts
- Updated every 5 seconds
- Deleted on normal completion or stop
- Persists on crash for recovery

---

### 5.4 Theme Definition Schema

**Location (Built-in):** `internal/theme/themes.go`  
**Location (User):** `~/.config/pomodux/themes/` (Post-MVP)

```yaml
# Example: Catppuccin Mocha Theme
name: "catppuccin-mocha"
version: "1.0"

colors:
  # Base colors
  background: "#1e1e2e"
  foreground: "#cdd6f4"
  
  # Accent colors
  primary: "#89b4fa"      # Blue
  secondary: "#cba6f7"    # Mauve
  success: "#a6e3a1"      # Green
  warning: "#f9e2af"      # Yellow
  error: "#f38ba8"        # Red
  
  # UI elements
  border: "#45475a"
  progress_filled: "#89b4fa"
  progress_empty: "#313244"
  text_muted: "#6c7086"

# Progress bar style
progress:
  filled_char: "█"
  empty_char: "░"
  show_percentage: true

# Border style
border:
  style: "rounded"  # Options: rounded, square, double, none
```

**Built-in Themes:**

1. **default**: Simple monochrome theme
2. **nord**: Nord color palette
3. **catppuccin-mocha**: Catppuccin Mocha palette

---

## CLI Reference

### `pomodux` - Timer Command

#### Start Timer with Preset
```bash
pomodux start <preset> [label]
```

**Arguments:**
- `<preset>`: Timer preset name from config (e.g., "work", "break")
- `[label]`: Optional session label (default: "Generic timer session")

**Examples:**
```bash
# Start 25m work session with default label
pomodux start work

# Start 25m work session with custom label
pomodux start work "Implementing authentication"

# Start 5m break session
pomodux start break "Coffee break"
```

**Output:** Interactive TUI with timer

---

#### Start Timer with Custom Duration
```bash
pomodux start <duration> [label]
```

**Arguments:**
- `<duration>`: Duration in format: `25m`, `1h`, `90s`, `1h30m`
- `[label]`: Optional session label (default: "Generic timer session")

**Examples:**
```bash
# Start 45-minute session
pomodux start 45m "Client meeting"

# Start 2-hour deep work session
pomodux start 2h "Deep work on refactoring"

# Start 90-second test timer
pomodux start 90s "Quick break"
```

**Output:** Interactive TUI with timer

---

#### 7.1.4 Help
```bash
pomodux --help
pomodux -h
```

**Output:**
```
Pomodux - Terminal-based Pomodoro Timer

USAGE:
  pomodux start <duration|preset> [label]
  pomodux --version
  pomodux --help

COMMANDS:
  start     Start a new timer session

OPTIONS:
  -h, --help       Show this help message
  -v, --version    Show version information

EXAMPLES:
  pomodux start work "Implementing auth"
  pomodux start 45m "Client meeting"
  pomodux start break

KEYBOARD CONTROLS (during timer):
  p       Pause timer
  r       Resume timer
  s, q    Stop timer and exit
  Ctrl+C  Emergency exit

FILES:
  Config:   ~/.config/pomodux/config.yaml
  History:  ~/.local/state/pomodux/history.json

LEARN MORE:
  Website:  https://github.com/yourusername/pomodux
  Docs:     https://github.com/yourusername/pomodux/wiki
```

---

#### 7.1.5 Version
```bash
pomodux --version
pomodux -v
```

**Output:**
```
pomodux version 1.0.0
Built with Go 1.21.5 on 2025-01-15
```

---

### `pomodux-stats` - Statistics Command

#### List Recent Sessions
```bash
pomodux-stats [options]
```

**Options:**
- `--limit N`: Show last N sessions (default: 20)
- `--all`: Show all sessions (no limit)
- `--today`: Show only today's sessions
- `--week`: Show last 7 days (Post-MVP)
- `--month`: Show last 30 days (Post-MVP)

**Examples:**
```bash
# Show last 20 sessions (default)
pomodux-stats

# Show last 50 sessions
pomodux-stats --limit 50

# Show all sessions
pomodux-stats --all

# Show only today's sessions
pomodux-stats --today
```

**Output:**
```
Recent Sessions
───────────────────────────────────────────────────────────────
Start Time          Duration  Status      Label
───────────────────────────────────────────────────────────────
2025-01-15 14:00    25m       completed   Implementing auth module
2025-01-15 14:30    5m        completed   Generic timer session
2025-01-15 14:40    25m       stopped     Code review PR #123
2025-01-15 15:10    25m       completed   Writing documentation
───────────────────────────────────────────────────────────────
Total: 4 sessions
```

---

#### Daily Statistics
```bash
pomodux-stats --today
```

**Output:**
```
Today's Statistics (2025-01-15)
───────────────────────────────────────────────────────────────
Total Time:           1h 20m
Completed Sessions:   3
Stopped Sessions:     1
Total Sessions:       4
Completion Rate:      75%
───────────────────────────────────────────────────────────────

Work Sessions:        2 (50m total)
Break Sessions:       1 (5m total)
Custom Sessions:      1 (25m total)
```

---

#### Export Sessions (Post-MVP)
```bash
pomodux-stats export [--format json|csv] [--output file]
```

**Priority:** P1 (Post-MVP)

---

#### Help
```bash
pomodux-stats --help
pomodux-stats -h
```

**Output:**
```
Pomodux Stats - View Timer Statistics

USAGE:
  pomodux-stats [options]

OPTIONS:
  --limit N        Show last N sessions (default: 20)
  --all            Show all sessions
  --today          Show today's sessions and statistics
  --week           Show last 7 days (coming soon)
  --month          Show last 30 days (coming soon)
  -h, --help       Show this help message
  -v, --version    Show version information

EXAMPLES:
  pomodux-stats
  pomodux-stats --today
  pomodux-stats --limit 50

FILES:
  History:  ~/.local/state/pomodux/history.json

LEARN MORE:
  Website:  https://github.com/yourusername/pomodux
```

---

## TUI Specification

### Screen Layouts

"Window" should be centered horizontally and vertically within the terminal window.

#### Running Timer State

```
┌─ Pomodoro Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │
│                                                              │
│  ████████████████████████░░░░░░░░░░░░░  60%   15:23          │
│                                                              │
│  Status: RUNNING                                             │
│                                                              │
│  [p] pause  [s] stop  [q] quit  [Ctrl+C] emergency exit      │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

**Elements:**
- Title: "Pomodoro Timer"
- Session Type + Label: "Work Session: Implementing authentication"
- Progress Bar: Unicode blocks, 60% filled, color-coded
- Time Remaining: "15:23" (MM:SS format)
- Status: "RUNNING" (green in themed output)
- Control Legend: Keyboard shortcuts

---

#### 8.1.2 Paused Timer State

```
┌─ Pomodoro Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │
│                                                              │
│  ████████████████████████░░░░░░░░░░░░░  60%   15:23          │
│                                                              │
│  Status: ⏸ PAUSED                                            │
│                                                              │
│  [r] resume  [s] stop  [q] quit  [Ctrl+C] emergency exit     │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

**Changes from Running:**
- Status: "⏸ PAUSED" (yellow/warning color)
- Progress bar: No animation
- Controls: Shows `[r] resume` instead of `[p] pause`

---

#### 8.1.3 Completion State

```
┌─ Pomodoro Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │
│                                                              │
│  ███████████████████████████████████████████████  100%  0:00 │
│                                                              │
│  Status: ✓ COMPLETED                                         │
│                                                              │
│  Session saved!                                              │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

**Changes from Running:**
- Progress bar: 100% filled
- Time: "0:00"
- Status: "✓ COMPLETED" (green/success color)
- Message: "Session saved!"
- Exits immediately after displaying completion (see [UXDR: Auto-Exit](../uxdr/auto-exit-behavior.md))

---

### Keyboard Controls

**See also:** [UXDR: Keyboard Controls Design](../uxdr/keyboard-controls-design.md) for detailed rationale and alternatives.

| Key | Action | Available States |
|-----|--------|------------------|
| `p` | Pause timer | Running |
| `r` | Resume timer | Paused |
| `s` | Stop timer and exit | Running, Paused |
| `q` | Stop timer and exit (alias) | Running, Paused |
| `Ctrl+C` | Emergency exit | All states |

**Behavior Notes:**
- All keys are single-press (no chords except Ctrl+C)
- Invalid keys in current state are ignored (no error)
- Keys are case-insensitive
- Control legend updates based on current state

---

### Terminal Resize Handling

**Requirements:**
- Detect terminal resize events (`SIGWINCH` on Unix)
- Recalculate layout immediately
- Maintain timer accuracy during resize
- Minimum terminal size: 80 columns x 24 rows
- Graceful degradation below minimum (show warning)

**Behavior:**
```
# If terminal too small:
┌─────────────────────────────────┐
│ Terminal too small!             │
│ Minimum: 80x24                  │
│ Current: 60x20                  │
└─────────────────────────────────┘
```

---

### Theme Application

**Theme affects:**
- Progress bar colors (filled/empty)
- Border style and color
- Text colors (primary, muted, success, warning, error)
- Background color (if terminal supports)

**Example color mappings:**

| Element | Default Theme | Nord Theme | Catppuccin Mocha |
|---------|---------------|------------|------------------|
| Progress (filled) | Cyan | Blue (#88c0d0) | Blue (#89b4fa) |
| Progress (empty) | Dark Gray | Polar Night (#3b4252) | Surface (#313244) |
| Status (running) | Green | Green (#a3be8c) | Green (#a6e3a1) |
| Status (paused) | Yellow | Yellow (#ebcb8b) | Yellow (#f9e2af) |
| Text (primary) | White | Snow Storm (#d8dee9) | Text (#cdd6f4) |

---

## Out of Scope

### Explicitly Deferred to Post-MVP

The following features are **not** included in MVP (v1.0.0) and are planned for future releases:

#### Plugin System Implementation (v1.1)
- **Deferred:** Plugin loading, Lua runtime, event dispatch
- **Included in MVP:** Event emission (internal architecture only)

#### Auto-Cycling Pomodoro Mode (v1.2)
- **Description:** Automatic 25→5→25→5→25→15 cycle
- **Rationale:** Adds complexity, not required for core use case
- **Alternative:** User manually starts work/break sessions

#### Advanced Statistics (v1.2)
- **Deferred:**
  - Weekly/monthly trends
  - Productivity heatmaps
  - Custom date range filtering (`--from`, `--to`)
  - Export to CSV/JSON
  - Charts/graphs
- **Included in MVP:** Basic daily statistics only

#### Custom Theme Loading (v1.3)
- **Deferred:** User-provided themes from `~/.config/pomodux/themes/`
- **Included in MVP:** 3 built-in themes only
- **Rationale:** Theme system architecture needs validation first

#### Multi-Pane TUI (v1.3)
- **Description:** Split screen with timer + history sidebar
- **Rationale:** Single-pane TUI is sufficient for MVP
- **Alternative:** Use separate `pomodux-stats` command

#### Custom Sound Notifications (v1.1 via plugin)
- **Description:** Custom audio files, desktop notification sounds, advanced audio alerts
- **Rationale:** Better suited for plugin implementation
- **Included in MVP:** Basic terminal bell (see [UXDR: Terminal Bell](../uxdr/terminal-bell-configuration.md))
- **Deferred:** Custom sounds, audio file playback, persistent audio notifications

#### Session Editing (v1.2)
- **Description:** Edit/delete past sessions
- **Rationale:** Read-only history is sufficient for MVP
- **Concern:** Data integrity, no undo mechanism

---

## Open Questions

### Technical Decisions

#### Q1: Timer State Persistence Strategy
**Question:** How should timer state be saved to disk for crash recovery?

**Options:**
- A. Time-based only (every N seconds, regardless of state)
- B. Event-driven only (save on state changes: start, pause, resume)
- C. Hybrid: Event-driven + periodic backup (every 5 seconds when running)

**Current Decision:** Option C (Hybrid)
**Rationale:**
- **Event-driven saves**: Immediate persistence on state changes (start, pause, resume, stop, interrupt)
- **Periodic backup**: Save every 5 seconds while timer is actively running (not paused)
- **Pause optimization**: When paused, only the pause event triggers save; no periodic saves while paused
- **Best of both worlds**: Efficient I/O, immediate response to user actions, maximum 5-second data loss on crash
- **No goroutines constraint**: Uses Bubbletea's `tea.Tick()` command system for periodic saves
- **Performance**: Reduces I/O operations from 60/session to 12/session over 25 minutes, minimizing event loop blocking

**Implementation:**
```
On TimerStarted:     → Save state, return tea.Tick(5*time.Second) command
On TimerPaused:      → Save state, stop tick commands
On TimerResumed:     → Save state, return tea.Tick(5*time.Second) command
On TimerStopped:     → Clean up state file
On TimerCompleted:   → Clean up state file
On Interrupted:      → Save state (critical)
Every 5s (running):  → Save state (backup) via tea.Tick command
```

---

#### Q2: History File Size Management
**Question:** Should we limit history file size or implement rotation?

**Options:**
- A. No limit (user manages manually)
- B. Keep last N days (e.g., 90 days)
- C. Keep last N sessions (e.g., 1000 sessions)
- D. Configurable limit in config file

**Proposed Decision:** Option D (configurable, default 90 days)  
**Rationale:** Flexibility, prevents unbounded growth

---

#### Q3: Theme Color Depth
**Question:** Should themes support 256-color or truecolor (24-bit)?

**Options:**
- A. 16-color (basic ANSI, maximum compatibility)
- B. 256-color (good balance)
- C. Truecolor (24-bit, best appearance)
- D. Adaptive (detect terminal capability)

**Proposed Decision:** Option D (adaptive)  
**Rationale:** Best appearance on capable terminals, fallback for compatibility

---

### Feature Scope Questions

#### Q4: Multiple Concurrent Timers
**Question:** Should users be able to run multiple timers simultaneously?

**Options:**
- A. No (single timer only, simpler implementation)
- B. Yes (multiple timers, each in separate TUI instance)
- C. Yes (single TUI with tabs/panes)

**Proposed Decision:** Option A (only)
**Rationale:** Single timer aligns with Pomodoro technique, reduces complexity

---

#### Q5: Session Tags/Categories (Post-MVP)
**Question:** Should sessions support tags beyond preset name?

**Example:**
```bash
pomodux start work "Auth module" --tags project:api,priority:high
```

**Proposed Decision:** Defer to later releases 
**Rationale:** Adds complexity, can be implemented via label convention for MVP

---

#### Q6: Time Tracking Integrations (Post-MVP)
**Question:** Should core app support time tracking APIs (Toggl, Harvest, etc.)?

**Proposed Decision:** Via plugins only (v1.1+)  
**Rationale:** Keep core focused, extensibility via plugin system

