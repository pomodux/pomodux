---
status: approved
---

# Separate Binaries for Timer and Statistics

## 1. Context / Background

### 1.1 Problem Statement
Pomodux requires two distinct user-facing interfaces:
1. **Timer Interface**: Interactive TUI for running timers with real-time updates
2. **Statistics Interface**: Read-only view of historical session data

The architectural decision is whether to implement these as:
- A single binary with subcommands (`pomodux start`, `pomodux stats`)
- Separate binaries (`pomodux` for timer, `pomodux-stats` for statistics)
- A hybrid approach (single codebase, multiple entry points)

### 1.2 Requirements
- **User Experience**: Commands should feel natural and intuitive
- **Separation of Concerns**: Timer and statistics have different UX patterns
- **Development Velocity**: Architecture should not complicate development
- **Maintenance**: Code should be maintainable as features grow
- **Unix Philosophy**: Tools should do one thing well
- **Distribution**: Package managers should handle both tools easily

## 2. Decision

**Selected Solution:** Separate binaries (`pomodux` and `pomodux-stats`)

### 2.1 Rationale

**1. Clear Mental Model**
- `pomodux` = "Start a timer" (action-oriented, interactive)
- `pomodux-stats` = "View my data" (query-oriented, read-only)
- Users don't need to remember subcommands or flags

**2. Unix Philosophy Alignment**
- Each binary does one thing well
- Follows patterns like `git`/`git-*`, `docker`/`docker-compose`
- Enables composition: `pomodux-stats --today | grep work`

**3. Different UX Patterns**
- **pomodux**: Interactive TUI, real-time updates, keyboard controls
- **pomodux-stats**: Simple output (table/list), scriptable, pipe-friendly
- Mixing these in one binary creates UX confusion

**4. Simplified Development**
- Timer logic isolated from statistics logic
- Easier to test each component independently
- Clearer code organization and responsibilities

**5. Future Extensibility**
- Easy to add more tools: `pomodux-export`, `pomodux-config`, etc.
- Plugin system can target specific binaries
- Users can install only what they need

**6. Distribution Benefits**
- Package managers handle multiple binaries naturally
- Users can symlink/alias as preferred
- `pomodux` can be in PATH, stats tool in extended PATH

## 3. Alternatives Considered

### 3.1 Single Binary with Subcommands
```bash
pomodux start [duration|preset] [label]
pomodux stats [--today|--week]
pomodux history
```

**Pros:**
- Single installation
- Unified help system
- Common pattern (git, docker, kubectl)

**Cons:**
- Timer and stats mixed in one codebase
- More complex CLI parsing
- `pomodux stats` feels verbose for frequent use
- Harder to compose with Unix pipes

### 3.2 Hybrid Approach
Single codebase, symlinks create multiple entry points:
```bash
pomodux -> main binary
pomodux-stats -> symlink to main binary (detects name)
```

**Pros:**
- Single binary distribution
- Multiple command names
- Shared code naturally

**Cons:**
- More complex entry point logic
- Confusing for debugging ("which binary am I?")
- Package manager complexity

### 3.3 Monolithic Single Binary
```bash
pomodux [duration|preset] [label]  # Implicit start
pomodux --stats                     # Flag for stats mode
```

**Pros:**
- Simplest distribution
- Minimal typing for timer start

**Cons:**
- Ambiguous interface (`pomodux 25` vs `pomodux --today`)
- Difficult to extend
- Poor Unix philosophy fit

## 4. Consequences

### 4.1 Positive
- **Clear Responsibilities**: Each binary has well-defined purpose
- **Better UX**: Users have intuitive mental model
- **Easier Testing**: Independent test suites for timer and stats
- **Maintainability**: Clear code boundaries between components
- **Unix Composability**: Stats tool works naturally with pipes and scripts
- **Future-Proof**: Easy to add new binaries without breaking existing ones

### 4.2 Negative
- **Distribution Complexity**: Package managers must install both binaries
- **Potential Duplication**: Must ensure shared code is truly shared
- **User Discovery**: Users might not know about `pomodux-stats` initially
- **Version Sync**: Both binaries must stay version-compatible

### 4.3 Risks and Mitigations

**Risk: Users don't discover `pomodux-stats`**
- **Mitigation**: `pomodux` displays hint on first completion: "View stats: pomodux-stats"
- **Mitigation**: Documentation prominently features both tools

**Risk: Version incompatibility between binaries**
- **Mitigation**: Both binaries built from same codebase/version
- **Mitigation**: Shared data models versioned explicitly
- **Mitigation**: Package managers install both together

**Risk: Code duplication in shared components**
- **Mitigation**: Strong `internal/` package organization
- **Mitigation**: Shared models in `internal/models/`
- **Mitigation**: Integration tests ensure compatibility

**Risk: Users expect `pomodux stats` to work**
- **Mitigation**: `pomodux stats` shows friendly error: "Use pomodux-stats instead"
- **Mitigation**: Shell completion suggests `pomodux-stats`
- **Mitigation**: Documentation explains rationale

