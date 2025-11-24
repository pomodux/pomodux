---
status: proposed
---

# CLI Framework Selection for Pomodux

## 1. Context / Background

### 1.1 Problem Statement

Pomodux requires a CLI framework to handle:
- Command parsing for `pomodux start <duration|preset> [label]`
- Help text generation (`pomodux --help`)
- Version display (`pomodux --version`)
- Shell completion (bash, zsh, fish)
- Flag parsing and validation
- Error handling for invalid commands

While Go's standard library `flag` package provides basic functionality, a dedicated CLI framework offers better UX, completion support, and maintainability.

### 1.2 Requirements

**Functional Requirements:**
- Parse commands: `start`, `--help`, `--version`
- Parse arguments: duration strings (`25m`, `1h30m`), preset names, optional labels
- Generate help text automatically
- Support shell completion (bash, zsh, fish)
- Validate input and provide clear error messages
- Cross-platform compatibility (Linux, macOS, Windows)

**Non-Functional Requirements:**
- Minimal learning curve (project goal: learn Go)
- Active maintenance and community support
- Well-documented with examples
- Lightweight (minimal binary size impact)
- Idiomatic Go patterns

**Project Constraints:**
- Simple command structure (no complex subcommands)
- Two separate binaries (`pomodux` and `pomodux-stats`)
- Learning-focused project (prefer simpler solutions)

## 2. Decision

**Selected Solution:** `urfave/cli` (v2)

### 2.1 Rationale

**Why urfave/cli?**

1. **Simplicity**: Clean, straightforward API that's easy to learn
   - Perfect for a learning-focused project
   - Less boilerplate than `cobra`
   - Intuitive command definition

2. **Sufficient Features**: Meets all requirements without over-engineering
   - Command parsing: ✅
   - Help generation: ✅
   - Flag parsing: ✅
   - Shell completion: ✅ (via `urfave/cli/v2/autocomplete`)
   - Error handling: ✅

3. **Active Maintenance**: Well-maintained with regular releases
   - 11k+ GitHub stars
   - Active development and community
   - Used by major projects (Docker CLI originally used v1)

4. **Lightweight**: Minimal dependencies and binary size impact
   - Fewer transitive dependencies than `cobra`
   - Smaller binary footprint
   - Faster compilation

5. **Good Documentation**: Clear examples and tutorials
   - Comprehensive README
   - Many examples in repository
   - Active community support

6. **Project Fit**: Matches Pomodux's simple command structure
   - No complex subcommand hierarchies needed
   - Two separate binaries (each can use urfave/cli independently)
   - Simple argument parsing requirements

### 2.2 Usage Example

**pomodux binary:**
```go
package main

import (
    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "pomodux",
        Usage: "Terminal-based Pomodoro timer",
        Commands: []*cli.Command{
            {
                Name:  "start",
                Usage: "Start a timer session",
                UsageText: "pomodux start <duration|preset> [label]",
                Action: startTimer,
            },
        },
        Flags: []cli.Flag{
            &cli.BoolFlag{
                Name:    "version",
                Aliases: []string{"v"},
                Usage:   "Show version information",
            },
        },
    }
    
    app.Run(os.Args)
}
```

**pomodux-stats binary:**
```go
app := &cli.App{
    Name:  "pomodux-stats",
    Usage: "View timer statistics",
    Flags: []cli.Flag{
        &cli.IntFlag{
            Name:    "limit",
            Aliases: []string{"l"},
            Usage:   "Show last N sessions",
            Value:   20,
        },
        &cli.BoolFlag{
            Name:    "today",
            Aliases: []string{"t"},
            Usage:   "Show today's statistics",
        },
        &cli.BoolFlag{
            Name:  "all",
            Usage: "Show all sessions",
        },
    },
    Action: showStats,
}
```

## 3. Alternatives Considered

### 3.1 Option A: cobra

**Approach:** Use `spf13/cobra` for CLI framework

**Pros:**
- Industry standard for complex CLIs (kubectl, docker, hugo)
- Rich feature set (subcommands, command groups, plugins)
- Excellent shell completion support
- Powerful help generation
- Large ecosystem and community
- Well-documented with many examples

**Cons:**
- More complex API (overkill for simple commands)
- Steeper learning curve
- More boilerplate code required
- Larger dependency footprint
- More features than needed for Pomodux

**Example Complexity:**
```go
// cobra requires more setup
var startCmd = &cobra.Command{
    Use:   "start [duration|preset] [label]",
    Short: "Start a timer",
    Long:  "Start a timer session...",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(startCmd)
}
```

**Rejected:** Over-engineered for Pomodux's simple command structure. Better suited for complex CLIs with many subcommands.

---

### 3.2 Option B: Standard Library `flag`

**Approach:** Use Go's built-in `flag` package

**Pros:**
- No external dependencies
- Simple and lightweight
- Part of standard library (always available)
- Minimal learning curve

**Cons:**
- No automatic help generation
- No shell completion support
- Manual command routing required
- More verbose error handling
- Less user-friendly (poor help text)
- No subcommand support (would need manual parsing)

**Example:**
```go
// Manual command parsing
if len(os.Args) < 2 {
    fmt.Fprintf(os.Stderr, "Usage: pomodux start <duration|preset> [label]\n")
    os.Exit(1)
}

command := os.Args[1]
switch command {
case "start":
    // Manual argument parsing
case "--help":
    // Manual help text
case "--version":
    // Manual version display
}
```

**Rejected:** Too low-level, requires too much manual work. Poor UX for help and completion.

---

### 3.3 Option C: urfave/cli v1

**Approach:** Use `urfave/cli` v1 (legacy version)

**Pros:**
- Mature and stable
- Well-tested in production
- Similar API to v2

**Cons:**
- Deprecated (maintenance mode only)
- No new features
- Migration path to v2 exists
- Not recommended for new projects

**Rejected:** v1 is deprecated. v2 is the recommended version.

---

### 3.4 Option D: kingpin

**Approach:** Use `alecthomas/kingpin` for CLI parsing

**Pros:**
- Clean API
- Good help generation
- Active maintenance

**Cons:**
- Smaller community than cobra/urfave
- Less documentation and examples
- Less familiar to most Go developers

**Rejected:** Less popular, fewer resources for learning.

## 4. Consequences

### 4.1 Positive

**Development Velocity:**
- Simple API enables rapid development
- Less boilerplate than `cobra`
- Quick to learn and implement

**Maintainability:**
- Clean, readable code structure
- Easy for contributors to understand
- Well-documented framework

**User Experience:**
- Automatic help generation
- Shell completion support
- Clear error messages
- Professional CLI appearance

**Learning Focus:**
- Simpler framework aligns with project's learning goals
- Less cognitive overhead
- Focus on Go fundamentals, not framework complexity

**Binary Size:**
- Smaller dependency footprint than `cobra`
- Faster compilation
- Minimal impact on binary size

### 4.2 Negative

**Feature Limitations:**
- Less powerful than `cobra` for complex subcommands
- Fewer advanced features (not needed for Pomodux)
- If command structure grows complex, might need migration

**Ecosystem:**
- Smaller ecosystem than `cobra`
- Fewer third-party integrations
- Less community support (though still substantial)

**Migration Risk:**
- If requirements change significantly, might need to migrate to `cobra`
- Mitigation: Command structure is simple and unlikely to grow complex

### 4.3 Risks and Mitigations

**Risk: Command structure becomes too complex for urfave/cli**
- **Likelihood:** Low (simple timer application)
- **Impact:** Medium (would require migration)
- **Mitigation:** Keep commands simple, document decision to avoid scope creep

**Risk: Shell completion not as polished as cobra**
- **Likelihood:** Low (urfave/cli has good completion support)
- **Impact:** Low (completion is nice-to-have, not critical)
- **Mitigation:** Test completion on target shells, document usage

**Risk: Framework becomes unmaintained**
- **Likelihood:** Very low (active development, large user base)
- **Impact:** Medium (would need to migrate)
- **Mitigation:** Monitor maintenance status, migration path exists

## 5. Implementation Guidelines

### 5.1 Command Structure

**pomodux binary:**
```go
app := &cli.App{
    Name:    "pomodux",
    Usage:   "Terminal-based Pomodoro timer",
    Version: version, // Injected at build time
    
    Commands: []*cli.Command{
        {
            Name:      "start",
            Usage:     "Start a timer session",
            UsageText: "pomodux start <duration|preset> [label]",
            Action:    startTimer,
        },
    },
    
    Flags: []cli.Flag{
        &cli.BoolFlag{
            Name:    "version",
            Aliases: []string{"v"},
            Usage:   "Show version information",
        },
        &cli.BoolFlag{
            Name:    "help",
            Aliases: []string{"h"},
            Usage:   "Show help",
        },
    },
}
```

**pomodux-stats binary:**
```go
app := &cli.App{
    Name:    "pomodux-stats",
    Usage:   "View timer statistics",
    Version: version,
    
    Flags: []cli.Flag{
        &cli.IntFlag{
            Name:    "limit",
            Aliases: []string{"l"},
            Usage:   "Show last N sessions",
            Value:   20,
        },
        &cli.BoolFlag{
            Name:    "today",
            Aliases: []string{"t"},
            Usage:   "Show today's statistics",
        },
        &cli.BoolFlag{
            Name:  "all",
            Usage: "Show all sessions",
        },
    },
    
    Action: showStats,
}
```

### 5.2 Shell Completion

**Installation:**
```bash
# Generate completion script
pomodux --generate-bash-completion > /etc/bash_completion.d/pomodux
pomodux --generate-zsh-completion > ~/.zsh/completions/_pomodux
```

**Implementation:**
- Use `urfave/cli/v2/autocomplete` package
- Generate completion scripts at build time or runtime
- Document installation in README

### 5.3 Error Handling

**Invalid Commands:**
```go
app := &cli.App{
    // ... config ...
    ExitErrHandler: func(c *cli.Context, err error) {
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            fmt.Fprintf(os.Stderr, "Run '%s --help' for usage\n", c.App.Name)
            os.Exit(1)
        }
    },
}
```

**Custom Validation:**
```go
func startTimer(c *cli.Context) error {
    args := c.Args()
    if args.Len() == 0 {
        return cli.Exit("Error: duration or preset required", 1)
    }
    
    durationOrPreset := args.Get(0)
    label := args.Get(1) // Optional
    
    // Validate and start timer
    return nil
}
```
