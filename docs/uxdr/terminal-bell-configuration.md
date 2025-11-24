---
status: approved
---

# Terminal Bell Configuration for Timer Completion

## 1. Context / Background

### 1.1 Problem Statement

When a timer completes, users need some form of notification to draw their attention back to the timer. The terminal bell (BEL character, ASCII 7) is a built-in terminal feature that can provide audio feedback.

**Key Questions:**
- Should the bell always ring on completion?
- Should it never ring?
- Should it be configurable?
- What should the default behavior be?

### 1.2 User Considerations

**Terminal Bell Characteristics:**
- Built-in to all terminals (no dependencies)
- Sound/behavior varies by terminal emulator
- Some terminals flash screen instead of sound
- Many users disable terminal bell system-wide
- Considered annoying by some users, helpful by others
- Works in all environments (local, SSH, tmux, etc.)

**User Preferences:**
- Some users rely on terminal bell for notifications
- Some users find terminal bells intrusive
- Some terminals have bell disabled at OS level
- Some users prefer visual-only feedback
- Some users use notification plugins instead

## 2. Decision

**Selected Solution:** Option C - Configurable (default: off)

### 2.1 Configuration

**Default Config:**
```yaml
timer:
  bell_on_complete: false
```

**To Enable:**
```yaml
timer:
  bell_on_complete: true
```

### 2.2 Rationale

**Why Configurable?**

1. **User Preference Varies Widely**: No single default satisfies all users
2. **Environment Dependent**: Bell behavior varies by terminal/system
3. **Plugin Alternative**: Users wanting notifications can use desktop notification plugins
4. **Visual Feedback Exists**: TUI shows completion message
5. **Modern Conventions**: Most modern CLI tools don't ring bell by default

**Why Default to Off?**

1. **Least Intrusive**: Avoids surprising users with sound
2. **Modern Expectations**: Terminal bells considered dated by many
3. **Plugin Preference**: Encourages use of modern notification systems
4. **Opt-In Philosophy**: Users who want it can easily enable it
5. **Professional Environments**: Avoids embarrassment in quiet offices

## 3. Alternatives Considered

### 3.1 Option A: Always Ring Bell

**Approach:**
```go
// Always ring bell on completion
fmt.Print("\a")  // BEL character
```

**Pros:**
- Guaranteed notification (if bell enabled in terminal)
- Simple implementation
- Works everywhere
- Traditional Unix behavior

**Cons:**
- Annoying to users who don't want it
- No way to disable without editing code
- Unprofessional in some environments
- Many terminals have bell disabled anyway
- No control over volume/sound

**Rejected:** Too inflexible, likely to annoy users

---

### 3.2 Option B: Never Ring Bell

**Approach:**
```go
// Never ring bell
// Completion notification only via:
// 1. Visual TUI update
// 2. Plugins (if configured)
```

**Pros:**
- Never intrusive
- Modern approach
- Encourages plugin usage
- No sound surprises

**Cons:**
- No built-in audio notification
- Users must configure plugins for sound
- Removes a simple, universal feature
- May miss completions without plugins

**Rejected:** Removes useful option for users who want it

---

### 3.3 Option C: Configurable (Default: Off) - Selected

**Approach:**
```yaml
# config.yaml
timer:
  bell_on_complete: false  # Default
```

**Implementation:**
```go
if config.Timer.BellOnComplete {
    fmt.Print("\a")
}
```

**Pros:**
- User choice
- Easy to enable for those who want it
- Default doesn't surprise users
- Simple configuration
- Works with all terminals

**Cons:**
- Requires configuration to enable
- Users might not discover the feature

**Selected:** Best balance of flexibility and modern defaults

## 4. Implementation

### 4.1 Configuration Schema

```yaml
timer:
  bell_on_complete: boolean  # default: false
```

### 4.2 Code Implementation

```go
// internal/config/config.go
type TimerConfig struct {
    BellOnComplete bool `yaml:"bell_on_complete"`
}

// internal/tui/timer.go
func (m model) onTimerComplete() tea.Cmd {
    // Save session
    m.saveSession()

    // Ring bell if configured
    if m.config.Timer.BellOnComplete {
        fmt.Print("\a")
    }

    // Show completion message
    return tea.Quit
}
```

### 4.3 Terminal Bell Behavior

**What Happens:**
- `\a` (BEL) character sent to terminal
- Terminal interprets based on its settings:
  - Audio beep (if bell enabled)
  - Visual flash (some terminals)
  - Nothing (if bell disabled)

**Cross-Platform:**
- Works on Linux, macOS, Windows
- Works in all terminal emulators
- Works over SSH
- Works in tmux/screen

## 5. Consequences

### 5.1 Positive

**User Empowerment:**
- Users choose notification method
- Simple toggle in config
- No code changes needed
- Can enable/disable anytime

**Modern Defaults:**
- Doesn't surprise users with sound
- Encourages notification plugins
- Professional out-of-box experience
- Respects user environment

**Flexibility:**
- Works with all terminals
- No dependencies
- Simple implementation
- Easy to document

### 5.2 Negative

**Discoverability:**
- Users might not know feature exists
- Requires reading documentation
- Not obvious in default config

**Limitations:**
- Can't control volume
- Can't control sound type
- Depends on terminal configuration
- Might not work if bell disabled

### 5.3 Mitigations

**Documentation:**
- Mention in README Quick Start
- Include in default config.yaml with comments
- Document in `pomodux --help` output
- Provide example in config file

**Default Config Comments:**
```yaml
# Timer behavior
timer:
  # Ring terminal bell when timer completes
  # Note: Behavior depends on terminal settings
  # Consider using notification plugins for more control
  bell_on_complete: false
```

## 6. User Guidance

### 6.1 When to Enable Terminal Bell

**Good Use Cases:**
- Working in a single terminal window
- Terminal always visible
- Familiar with terminal bell sound
- Simple notification preference
- Working in minimal environment

**Example:**
```yaml
timer:
  bell_on_complete: true
```

### 6.2 When to Use Plugins Instead

**Better Use Cases:**
- Need desktop notifications
- Want custom sounds
- Working across multiple desktops
- Need persistent notifications
- Want notification actions (snooze, etc.)

**Example:**
```yaml
timer:
  bell_on_complete: false

plugins:
  enabled:
    - notify  # Desktop notifications
```

### 6.3 Using Both

```yaml
timer:
  bell_on_complete: true  # Terminal bell

plugins:
  enabled:
    - notify  # Plus desktop notification
```

## 7. Future Enhancements

### 7.1 Advanced Bell Configuration

**Could Add:**
```yaml
timer:
  bell_on_complete: true
  bell_repeat: 3          # Ring 3 times
  bell_interval: 1000     # 1 second apart
```

**Evaluation Criteria:**
- User requests for feature
- Clear use case
- Terminal support for repeated bells
- Implementation complexity

### 7.2 Custom Bell Sounds

**Potential:**
```yaml
timer:
  bell_on_complete: true
  bell_command: "afplay /path/to/sound.mp3"
```

**Concerns:**
- Platform-specific
- Better handled by plugins
- Adds complexity
- Breaks from simple terminal bell

**Recommendation:** Use notification plugins for custom sounds

