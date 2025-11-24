---
status: approved
---

# Concurrency Model - No Manual Goroutines

## 1. Context / Background

### 1.1 Problem Statement

Pomodux requires several operations that could be implemented concurrently:
- Timer countdown updates (every 250ms for display)
- Periodic state persistence (every 5 seconds)
- Event dispatch to plugins (future)
- Keyboard input handling
- Terminal resize event handling
- File I/O operations (config, history, state)

The decision is whether to use manual goroutine management for these concurrent operations or rely on the Bubbletea framework's event loop model.

### 1.2 Requirements Affected

**Timer Accuracy:**
- Must complete at exact configured duration (within 100ms system clock precision)
- Display updates at 250ms intervals without affecting underlying calculation

**State Persistence:**
- Event-driven saves on state changes (start, pause, resume, stop, interrupt)
- Periodic backup every 5 seconds while running
- Maximum 5-second data loss on crash (SIGKILL)

**Code Maintainability:**
- Simple codebase suitable for learning Go
- Minimal cognitive overhead for understanding control flow
- Easy debugging and testing

### 1.3 Concurrency Challenges

Manual goroutine management introduces complexity:
- **Race conditions**: Shared state between timer, UI, and persistence
- **Synchronization**: Mutexes, channels, and WaitGroups needed
- **Debugging**: Race detector, deadlock detection, non-deterministic bugs
- **Testing**: Harder to write deterministic tests for concurrent code
- **Learning curve**: Goroutines and channels add conceptual overhead

## 2. Decision

**Selected Solution:** No manual goroutines - use Bubbletea's Cmd-based event loop exclusively

### 2.1 Rationale

**1. Bubbletea Already Provides Concurrency**
- Keyboard input handled in background goroutines (internal to framework)
- Terminal resize events managed by framework
- `tea.Tick()` provides periodic events without manual timers
- `tea.Cmd` allows async operations to be queued and executed

**2. Simplicity Over Performance**
- Single-threaded event loop is easier to reason about
- No race conditions possible (all state mutations in Update())
- Deterministic execution order
- Easier to debug: linear control flow

**3. Performance Is Adequate**
- Timer application is not computationally intensive
- All operations (state save, rendering, input) are fast (<50ms)
- Event loop can easily handle 4 FPS display updates
- No performance bottlenecks identified

**4. Learning-Focused Design**
- Project goal: "Serve as a way for me to learn the Go programming language"
- Avoiding goroutines allows focus on:
  - Go fundamentals (types, interfaces, error handling)
  - File I/O and configuration management
  - TUI development with Bubbletea
  - Testing and code organization
- Goroutines can be learned in future projects

**5. Matches Bubbletea's Philosophy**
- Elm Architecture pattern is inherently single-threaded
- All state updates go through Update() function
- Framework handles concurrency internally
- Fighting the framework leads to complexity

### 2.2 Implementation Approach

**Timer Accuracy (Wall-Clock Calculation):**
```go
type model struct {
    startTime    time.Time
    duration     time.Duration
    pausedAt     time.Time
    totalPaused  time.Duration
    isPaused     bool
}

func (m model) remaining() time.Duration {
    elapsed := time.Since(m.startTime) - m.totalPaused
    return m.duration - elapsed
}

// Update() called every 250ms via tea.Tick()
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case tickMsg:
        if m.remaining() <= 0 {
            return m.completeTimer()
        }
        return m, tea.Tick(250*time.Millisecond, ...)
}
```

**State Persistence (Event-Driven + Periodic):**
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case timerStartedMsg:
        m.saveState()
        return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
            return saveStateMsg{}
        })

    case saveStateMsg:
        m.saveState()
        if m.isRunning() {
            return m, tea.Tick(5*time.Second, ...)
        }
        return m, nil
}
```

**Plugin Event Dispatch (Synchronous):**
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case timerCompletedMsg:
        // Synchronous plugin execution
        for _, plugin := range m.plugins {
            plugin.OnTimerCompleted(event) // Blocks until complete
        }
        return m, nil
}
```

## 3. Alternatives Considered

### 3.1 Manual Goroutines with Channels

**Approach:**
```go
type Timer struct {
    remaining   time.Duration
    updateChan  chan time.Duration
    controlChan chan ControlMsg
}

func (t *Timer) Start() {
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        for {
            select {
            case <-ticker.C:
                t.remaining--
                t.updateChan <- t.remaining
            case msg := <-t.controlChan:
                // Handle pause/resume/stop
            }
        }
    }()
}
```

**Pros:**
- More control over timing precision
- Can run operations truly in parallel
- Closer to "idiomatic Go" for concurrent operations

**Cons:**
- Race conditions on shared state (timer model, TUI state)
- Need mutexes or careful channel design
- Harder to test deterministically
- More complex debugging (race detector, deadlocks)
- Violates Bubbletea's single-threaded Update() pattern
- High cognitive overhead for this simple application

**Rejected:** Complexity outweighs benefits for timer application

---

### 3.2 Actor Model (errgroup, worker pools)

**Approach:**
```go
import "golang.org/x/sync/errgroup"

func (m model) Start() error {
    g, ctx := errgroup.WithContext(context.Background())

    g.Go(func() error { return m.runTimer(ctx) })
    g.Go(func() error { return m.savePeriodically(ctx) })
    g.Go(func() error { return m.runTUI(ctx) })

    return g.Wait()
}
```

**Pros:**
- Clean error propagation
- Structured concurrency
- Graceful shutdown on errors

**Cons:**
- Still requires synchronization between actors
- Doesn't integrate with Bubbletea's model
- Overkill for single-user desktop application
- Adds complexity without clear benefits

**Rejected:** Doesn't fit Bubbletea architecture

---

### 3.3 Hybrid: Bubbletea + Background Goroutines

**Approach:**
```go
// Timer runs in background, sends updates via Cmd
func timerTickerCmd(remaining time.Duration) tea.Cmd {
    return func() tea.Msg {
        go func() {
            time.Sleep(1 * time.Second)
            // Send update back to Bubbletea
        }()
        return tickMsg{}
    }
}
```

**Pros:**
- Can have "true" background work
- Still integrates with Bubbletea

**Cons:**
- Introduces goroutine complexity anyway
- Bubbletea already handles this internally
- Easy to create race conditions accidentally
- Violates single-responsibility (mixing models)

**Rejected:** Doesn't provide sufficient benefit over pure Bubbletea

---

### 3.4 Async/Await Pattern (Future Go Feature)

**Note:** Go does not have async/await as of Go 1.22

**Hypothetical Approach:**
```go
async func saveState() {
    await file.Write(state)
}
```

**Rejected:** Not available in Go, purely theoretical

## 4. Consequences

### 4.1 Positive

**Simplicity:**
- Single-threaded mental model
- All state mutations in one place (Update function)
- Easy to trace execution flow
- No race conditions possible

**Testability:**
- Deterministic test execution
- No need for synchronization primitives in tests
- Easy to mock time-based operations
- Straightforward unit tests for Update() function

**Debugging:**
- Linear stack traces
- No race detector needed
- No deadlock concerns
- Predictable control flow

**Learning:**
- Focus on Go fundamentals, not concurrency
- Understand Elm Architecture pattern
- Practice functional programming concepts
- Learn Bubbletea framework deeply

**Maintainability:**
- Clear ownership of state
- Easy for future contributors to understand
- No subtle concurrency bugs
- Simple mental model

### 4.2 Negative

**Blocking Operations:**
- File I/O happens in event loop (can pause UI briefly)
- Plugin hooks execute synchronously (slow plugins freeze UI)
- No true parallelism for independent operations

**Perceived Performance:**
- State saves may cause brief UI freezes (~10-50ms)
- Multiple concurrent timers not possible (already not a requirement)
- Heavy plugins could impact responsiveness

**Limited Concurrency Practice:**
- Don't learn goroutines/channels in this project
- Miss opportunity to practice concurrent Go patterns

**Plugin System Constraints:**
- Plugins must complete quickly (<100ms)
- Long-running operations not supported
- Network calls in plugins will freeze UI
- Limits plugin complexity

### 4.3 Risks and Mitigations

**Risk: UI freezing from blocking operations**
- **Likelihood:** Medium (depends on filesystem performance)
- **Impact:** Low (brief freezes, <50ms acceptable)
- **Mitigation:**
  - Use fast atomic file writes
  - Test on slow storage (HDD, NFS)
  - Reduce state save frequency (5s instead of 1s)
  - Document plugin performance requirements

**Risk: Plugin system too limited**
- **Likelihood:** Medium (synchronous-only limits use cases)
- **Impact:** Medium (some plugins impossible to implement)
- **Mitigation:**
  - Document plugin constraints clearly
  - Suggest external tools for heavy operations
  - Consider allowing plugins to spawn their own goroutines (isolated)
  - Defer complex plugin use cases to post-MVP

**Risk: Missing learning opportunity for goroutines**
- **Likelihood:** High (definitely won't learn goroutines here)
- **Impact:** Low (can learn in future projects)
- **Mitigation:**
  - Document this as intentional decision
  - Plan future project to explore concurrency
  - Learn goroutines through Go tutorials separately

**Risk: Future features require goroutines**
- **Likelihood:** Low (most timer features don't need concurrency)
- **Impact:** Medium (might need to refactor)
- **Mitigation:**
  - Keep architecture flexible
  - Isolate state management in clear boundaries
  - Can add goroutines later if truly needed
  - YAGNI principle: don't add complexity until required

## 5. Implementation Guidelines

### 5.1 Rules for Developers

**DO:**
- ✅ Use `tea.Tick()` for periodic operations
- ✅ Use `tea.Cmd` for async-like behavior
- ✅ Calculate timer values from absolute time (`time.Since`)
- ✅ Keep all state mutations in `Update()` function
- ✅ Use Bubbletea's message passing for events

**DON'T:**
- ❌ Spawn goroutines with `go func()`
- ❌ Use channels for inter-component communication
- ❌ Use mutexes or sync primitives
- ❌ Store mutable state outside of the model
- ❌ Perform blocking operations in View() (read-only)

### 5.2 Exception Clause

Goroutines are permitted in the following limited cases:

1. **Inside Bubbletea's Cmd functions** (framework-managed)
2. **Inside plugin code** (isolated, plugin's responsibility)
3. **External libraries** (e.g., Bubbletea framework internals)

All exceptions must be documented and justified.
