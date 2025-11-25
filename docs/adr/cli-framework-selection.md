---
status: accepted
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
- Active maintenance and community support
- Well-documented with examples
- Industry-standard patterns and practices
- Strong ecosystem and third-party integrations
- Idiomatic Go patterns

**Project Constraints:**
- Simple command structure (no complex subcommands)
- Two separate binaries (`pomodux` and `pomodux-stats`)
- Learning-focused project (exposure to industry standards)

## 2. Decision

**Selected Solution:** `spf13/cobra`

### 2.1 Rationale

**Why Cobra?**

1. **Industry Standard**: Dominant choice in the Go ecosystem
   - 41.7K+ GitHub stars (as of 2025)
   - Imported by 182,338+ projects
   - Used by Kubernetes, GitHub CLI, Hugo, Docker, etcd
   - De facto standard for Go CLIs

2. **Comprehensive Features**: Production-ready feature set
   - Command parsing: ✅
   - Help generation: ✅ (highly polished)
   - Flag parsing: ✅ (including persistent flags)
   - Shell completion: ✅ (bash, zsh, fish, PowerShell)
   - Error handling: ✅
   - Configuration integration: ✅ (Viper integration)

3. **Active Maintenance**: Continuously maintained and improved
   - Last updated August 2025
   - Regular releases and security updates
   - Large contributor base
   - Well-established governance

4. **Strong Ecosystem**: Rich tooling and integrations
   - `cobra-cli` generator for scaffolding
   - Seamless Viper integration for config management
   - Extensive third-party plugins and extensions
   - Abundant community resources and examples

5. **Professional UX**: Best-in-class user experience
   - Automatic help and usage generation
   - Sophisticated completion support
   - Consistent error messages
   - Support for command aliases and deprecation

6. **Learning Value**: Exposure to industry-standard tools
   - Learning Cobra means understanding most Go CLI tools
   - Transferable knowledge to major projects
   - Industry-recognized patterns and practices
   - Better for portfolio and future contributions

7. **Future-Proofing**: Well-suited for growth
   - Handles simple commands elegantly
   - Can scale to complex command hierarchies if needed
   - No migration needed if requirements evolve
   - Proven at scale in production environments

## 2. Alternatives Considered

### 2.1 Option A: urfave/cli (v2)

**Approach:** Use `urfave/cli` for CLI framework

**Pros:**
- Simpler API than Cobra
- Less boilerplate code
- Lightweight with fewer dependencies
- Good documentation
- Active maintenance (23.7K stars, v3.4.1 released August 2025)
- Used by Docker Machine, Drone, Gogs

**Cons:**
- Significantly less popular (23.7K vs 41.7K stars)
- Smaller ecosystem (270 contributors vs 182K+ importing projects)
- Less sophisticated completion support
- Fewer third-party integrations
- Some projects migrating away to Cobra
- Less exposure to industry-standard patterns

**Migration Patterns:**
Several projects have migrated from urfave/cli to Cobra, citing:
- Better persistent flag handling
- Better configuration management integration (Viper)
- More robust for complex command structures

**Rejected:** While simpler, Cobra's industry dominance, larger ecosystem, and better long-term support outweigh the modest increase in complexity. For a learning-focused project, exposure to industry-standard tools provides more value.

---

### 2.2 Option B: Standard Library `flag`

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

**Rejected:** Too low-level, requires too much manual work. Poor UX for help and completion. Not suitable for user-facing CLI tools.

---

### 2.3 Option C: kingpin

**Approach:** Use `alecthomas/kingpin` for CLI parsing

**Pros:**
- Clean API
- Good help generation
- Active maintenance

**Cons:**
- Much smaller community than Cobra/urfave
- Less documentation and examples
- Less familiar to most Go developers
- Limited ecosystem

**Rejected:** Insufficient community support and resources for learning.

---

### 2.4 Option D: kong

**Approach:** Use `alecthomas/kong` for struct-based CLI parsing

**Pros:**
- Declarative struct-based approach
- Minimal boilerplate
- Modern design

**Cons:**
- Very different paradigm from most Go CLIs
- Much smaller adoption
- Less documentation
- Newer project with less production validation

**Rejected:** Too unconventional for a learning-focused project. Limited industry adoption.

## 3. Consequences

### 3.1 Positive

**Industry Alignment:**
- Learning industry-standard tools and patterns
- Code structure familiar to Go developers
- Easier to attract contributors familiar with Cobra
- Portfolio value (recognizable by employers)

**Ecosystem Benefits:**
- Extensive community resources and examples
- Rich third-party integrations available
- `cobra-cli` generator for scaffolding
- Seamless Viper integration for configuration

**User Experience:**
- Polished help generation and formatting
- Superior shell completion support
- Consistent with major Go tools (kubectl, gh, hugo)
- Professional-grade CLI UX

**Future-Proofing:**
- Handles simple commands without overhead
- Room to grow if complexity increases
- No migration needed as project evolves
- Proven at scale in production

**Maintainability:**
- Clear command structure and separation
- Well-documented patterns
- Large knowledge base for troubleshooting
- Active community support

### 3.2 Negative

**Increased Complexity:**
- More boilerplate than urfave/cli
- Steeper initial learning curve
- More concepts to understand (root commands, subcommands, persistent flags)
- Slightly larger binary size

**Learning Overhead:**
- More time needed to understand framework
- More moving parts than simpler alternatives
- May be overkill for simple command structure
- Additional concepts beyond core Go

### 3.3 Trade-offs Accepted

**Complexity vs. Industry Standard:**
- Accept: Modest increase in complexity
- Gain: Exposure to industry-standard patterns and massive ecosystem

**Binary Size vs. Features:**
- Accept: Slightly larger binary (still reasonable for CLI tool)
- Gain: Superior UX, completion support, and professional polish

**Learning Curve vs. Future Value:**
- Accept: More upfront learning time
- Gain: Transferable knowledge applicable to major Go projects

### 3.4 Risks and Mitigations

**Risk: Framework overhead for simple commands**
- **Likelihood:** Medium (Pomodux has simple command structure)
- **Impact:** Low (extra boilerplate is manageable)
- **Mitigation:** Cobra handles simple commands well; overhead is minimal for our use case

**Risk: Steeper learning curve slows development**
- **Likelihood:** Low (Cobra is well-documented)
- **Impact:** Low (one-time learning cost)
- **Mitigation:** Excellent documentation and examples; `cobra-cli` generator simplifies setup

**Risk: Framework becomes unmaintained**
- **Likelihood:** Very low (182K+ projects depend on it)
- **Impact:** High (would require migration)
- **Mitigation:** Cobra is critical infrastructure for Go ecosystem; too big to fail

