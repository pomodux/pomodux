---
status: approved
---

# Terminal User Interface Framework Selection

## 1. Context / Background

### 1.1 Problem Statement

Pomodux requires an interactive Terminal User Interface (TUI) to display:
- Real-time countdown timer
- Progress bar visualization
- Session information (label, preset, status)
- Keyboard controls and status indicators
- Responsive layout that adapts to terminal resize

Building a TUI from scratch involves handling:
- Terminal control sequences (ANSI escape codes)
- Keyboard input (raw mode, non-blocking I/O)
- Screen rendering and buffering
- Terminal resize events
- Cross-platform terminal compatibility
- Efficient screen updates (avoiding flicker)

This complexity warrants using a TUI framework rather than manual implementation.

### 1.2 Requirements

**Functional Requirements:**
- Real-time updates (countdown, progress bar)
- Keyboard input handling (single keypresses: p, r, s, q, Ctrl+C)
- Terminal resize detection and responsive layout
- Unicode support for progress bar characters
- Color/theme support
- Cross-platform (Linux, macOS, Windows)

**Non-Functional Requirements:**
- Minimal resource overhead (timer accuracy critical)
- Active maintenance and community support
- Well-documented with examples
- Idiomatic Go patterns
- Testable architecture
- No heavy dependencies

### 1.3 Design Constraints

- **Timer Accuracy**: UI updates must not affect timer precision
- **Single-Pane Layout**: MVP requires simple layout (no tabs/splits)
- **Keyboard-First**: No mouse support needed
- **Theme Integration**: Must support custom color schemes
- **Learning Curve**: Should be learnable within the Go learning timeline

## 2. Decision

**Selected Solution:** Charm's Bubbletea (https://github.com/charmbracelet/bubbletea)

### 2.1 Rationale

**Why Bubbletea?**

1. **Elm Architecture**: Clean, predictable state management pattern
   - Model-Update-View pattern separates concerns
   - Immutable state updates
   - Easy to reason about and test

2. **Active Development**: Part of Charm ecosystem
   - 26k+ GitHub stars
   - Active maintenance and releases
   - Large community and ecosystem
   - Well-funded company backing

3. **Rich Ecosystem**: Complementary libraries
   - `lipgloss`: Styling and layout
   - `bubbles`: Pre-built components (progress bars, spinners, etc.)
   - `harmonica`: Animation helpers
   - All designed to work together

4. **Production-Ready**: Battle-tested in many projects
   - Glow (markdown reader)
   - Soft Serve (Git server TUI)
   - VHS (terminal recorder)
   - Hundreds of community projects

5. **Excellent Documentation**:
   - Comprehensive tutorials
   - Many examples and templates
   - Active Discord community
   - Well-documented API

6. **Performance**: Efficient rendering
   - Smart screen diffing (only updates changed regions)
   - Non-blocking I/O
   - Minimal overhead

7. **Cross-Platform**: Works everywhere
   - Linux, macOS, Windows support
   - Handles terminal quirks automatically
   - Graceful degradation for limited terminals

### 2.2 Component Integration

**Bubbles Components We'll Use:**

1. **Progress Bar** (`bubbles/progress`)
   - Built-in progress bar component
   - Customizable width and style
   - Percentage display
   - Unicode character support

2. **Styling** (`lipgloss`)
   - Color management
   - Layout primitives
   - Border styles
   - Text alignment

**Example Implementation:**
```go
import (
    "github.com/charmbracelet/bubbles/progress"
    "github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    progress progress.Model
    // ... other fields
}

func (m model) Init() tea.Cmd {
    return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "p":
            m.paused = true
            return m, nil
        case "r":
            m.paused = false
            return m, tickCmd()
        case "q", "s":
            return m, tea.Quit
        }
    case tickMsg:
        if !m.paused {
            m.updateTimer()
            return m, tickCmd()
        }
    }
    return m, nil
}

func (m model) View() string {
    progressBar := m.progress.ViewAs(m.remaining / m.duration)
    
    return lipgloss.NewStyle().
        Padding(1, 2).
        Render(
            lipgloss.JoinVertical(
                lipgloss.Center,
                m.label,
                progressBar,
                m.formatTime(),
                m.renderControls(),
            ),
        )
}
```

---

### 2.3 Available Bubbles Components

The Bubbles library provides additional pre-built components we may use as the application evolves:

**Confirmed for MVP:**
- ✅ `progress`: Progress bars for timer visualization
- ✅ `lipgloss`: Styling and layout management

**Available for Future Enhancement:**
- `key`: Key binding management (useful for keyboard controls)
- `list`: Scrollable lists (useful for session history in TUI)
- `textinput`: Text input fields (if we add inline label editing)
- `table`: Tables (alternative for session display)
- `spinner`: Loading indicators (for async operations)
- `paginator`: Pagination (for long session lists)
- `viewport`: Scrollable content (for help text, logs)
- `textarea`: Multi-line text input (for notes, descriptions)
- `stopwatch`: Timer utilities (complementary to our timer logic)
- `help`: Help text display (for keyboard shortcut legend)
- `filepicker`: File selection (if we add import/export features)

**Component Reference:** https://github.com/charmbracelet/bubbles

**Design Principle:** This ADR approves the Bubbles ecosystem as a whole. Specific component selection for post-MVP features is an implementation detail, not requiring separate architectural decisions. Components should be evaluated on a case-by-case basis for appropriateness and fit.

---

### 2.4 Theme Integration

Bubbletea + Lipgloss make theming straightforward:

```go
// internal/theme/theme.go
type Theme struct {
    Primary    lipgloss.Color
    Secondary  lipgloss.Color
    Background lipgloss.Color
    // ...
}

func (t Theme) ProgressStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(t.Primary).
        Background(t.Background)
}

// Apply theme in TUI
func (m model) View() string {
    theme := m.config.Theme
    
    title := theme.TitleStyle().Render(m.label)
    progress := theme.ProgressStyle().Render(m.progressBar())
    
    return lipgloss.JoinVertical(lipgloss.Center, title, progress)
}
```

## 3. Alternatives Considered

### 3.1 tcell (https://github.com/gdamore/tcell)

**Approach:** Low-level terminal cell buffer library

**Pros:**
- More control over rendering
- Mature and stable (v2.7+)
- Used by many TUI apps
- Good performance

**Cons:**
- Lower-level API (more code to write)
- Manual state management
- More complex event handling
- Steeper learning curve
- No built-in components (progress bars, etc.)

**Example:**
```go
// More verbose, manual screen management
screen, _ := tcell.NewScreen()
screen.Init()
defer screen.Fini()

for {
    ev := screen.PollEvent()
    // Manual event handling
    // Manual screen drawing
}
```

**Rejected:** Too low-level for our needs, slower development

---

### 3.2 tview (https://github.com/rivo/tview)

**Approach:** High-level TUI framework with widgets

**Pros:**
- Rich widget library (tables, forms, lists)
- Event-driven architecture
- Good documentation
- Active maintenance

**Cons:**
- Heavier framework (more features than needed)
- Widget-based approach less flexible for custom layouts
- Not as idiomatic Go patterns
- Smaller ecosystem than Charm

**Example:**
```go
app := tview.NewApplication()
textView := tview.NewTextView()
// Widget-centric approach
```

**Rejected:** Overkill for simple timer display, less flexible

---

### 3.3 termui (https://github.com/gizak/termui)

**Approach:** Dashboard-style TUI with widgets

**Pros:**
- Good for dashboards
- Chart and graph support
- Simple API

**Cons:**
- Less active maintenance (last release 2021)
- Dashboard-focused (not ideal for interactive apps)
- Smaller community
- Limited keyboard interaction patterns

**Rejected:** Stale maintenance, not ideal for our use case

---

### 3.4 Custom Implementation (ANSI Escape Codes)

**Approach:** Direct terminal control using ANSI sequences

**Pros:**
- No dependencies
- Full control
- Minimal overhead

**Cons:**
- High implementation complexity
- Cross-platform challenges
- Manual keyboard handling (raw mode, etc.)
- Manual screen buffering
- High maintenance burden
- Reinventing the wheel

**Example:**
```go
// Manual ANSI sequences
fmt.Print("\033[2J")        // Clear screen
fmt.Print("\033[H")         // Move cursor to home
fmt.Print("\033[31mRed\033[0m") // Color text
// ... hundreds of lines of terminal control
```

**Rejected:** Not worth the effort, too much complexity

---

### 3.5 Lip Gloss Only (No Bubbletea)

**Approach:** Use lipgloss for styling, manual event loop

**Pros:**
- Lighter than full Bubbletea
- Good styling capabilities

**Cons:**
- Manual event handling
- No structured state management
- Missing Bubbletea's update loop pattern
- Still need keyboard input handling
- More complex testing

**Rejected:** Loses Bubbletea's architectural benefits

## 4. Consequences

### 4.1 Positive

- **Rapid Development**: Pre-built components and patterns accelerate development
- **Testable**: Clean state management makes unit testing straightforward
- **Maintainable**: Well-structured code with clear separation of concerns
- **Extensible**: Easy to add features (animations, more components, etc.)
- **Community**: Large ecosystem and active support
- **Documentation**: Excellent tutorials and examples
- **Future-Proof**: Active development ensures long-term viability

### 4.2 Negative

- **Dependency**: Adds external dependency (Charm ecosystem)
- **Learning Curve**: Team must learn Elm Architecture pattern
- **Abstraction**: Less control than low-level libraries
- **Binary Size**: Adds ~2-3MB to binary (acceptable for desktop app)

### 4.3 Risks and Mitigations

**Risk: Framework becomes unmaintained**
- **Likelihood**: Very low (well-funded company, active development)
- **Impact**: Medium (would need to migrate)
- **Mitigation**: Charm has track record of maintenance, large user base

**Risk: Performance issues with frequent updates**
- **Likelihood**: Low (Bubbletea is well-optimized)
- **Impact**: High (timer accuracy critical)
- **Mitigation**: Profile early, optimize update frequency, use smart diffing

**Risk: Terminal compatibility issues**
- **Likelihood**: Low (Bubbletea handles most terminal quirks)
- **Impact**: Medium (some users might have display issues)
- **Mitigation**: Test on multiple terminals, graceful degradation

**Risk: Learning curve delays development**
- **Likelihood**: Low (good documentation, simple pattern)
- **Impact**: Low (pattern is intuitive)
- **Mitigation**: Follow official tutorials, reference examples

