---
status: approved
---

# Keyboard Controls Design - Single-Key Actions

## 1. Context / Background

### 1.1 Problem Statement

Terminal applications require keyboard input for user interaction. The design of keyboard controls significantly impacts usability and user efficiency.

**Design Questions:**
- Single-key commands vs key chords (Ctrl+X, Alt+Y)?
- Case-sensitive vs case-insensitive?
- Vim-like keybindings vs traditional?
- Mnemonic keys vs position-based?
- Mouse support vs keyboard-only?

### 1.2 Requirements

**From Requirements Doc:**
- All functionality must be accessible via keyboard only
- No mouse required for any operation
- Keyboard shortcuts displayed in TUI at all times
- Single-key actions (no chords except Ctrl+C)
- Works with screen readers (basic support)

**User Context:**
- Users are CLI-focused, comfortable with keyboard
- Timer is running in foreground (not background daemon)
- Limited actions needed (pause, resume, stop)
- Quick reactions important (pause on interruption)

## 2. Decision

**Selected Solution:** Single-key, case-insensitive, mnemonic controls

### 2.1 Control Scheme

| Key       | Action              | States Available     | Mnemonic            |
|-----------|---------------------|----------------------|---------------------|
| `p`       | Pause timer         | Running              | **P**ause           |
| `r`       | Resume timer        | Paused               | **R**esume          |
| `s`       | Stop and exit       | Running, Paused      | **S**top            |
| `q`       | Stop and exit       | Running, Paused      | **Q**uit            |
| `Ctrl+C`  | Emergency exit      | All states           | (Universal)         |

**Behavior:**
- Keys work in both uppercase and lowercase (`p` = `P`)
- Invalid keys in current state are silently ignored
- Control legend in TUI shows available keys for current state
- No key combinations required (except Ctrl+C)

### 2.2 Rationale

**Why Single-Key?**

1. **Fastest Response**: One keypress, no modifier keys
2. **Low Cognitive Load**: Easy to remember and execute
3. **Accessibility**: Easier for users with mobility limitations
4. **Focus-Friendly**: Quick pause when interrupted
5. **Terminal Standard**: Matches tools like `less`, `htop`, `vim` (in some contexts)

**Why Case-Insensitive?**

1. **Forgiving**: Don't punish Caps Lock mistakes
2. **Accessibility**: Easier for some users
3. **Modern Expectation**: Most modern CLIs are case-insensitive
4. **No Ambiguity**: We don't need 26+ commands

**Why Mnemonic?**

1. **Easy to Learn**: 'p' for pause, 'r' for resume
2. **Self-Documenting**: Keys match action names
3. **Intuitive**: New users can guess controls
4. **Memorable**: Harder to forget

## 3. Alternatives Considered

### 3.1 Option A: Key Chords (Vim-Style)

**Approach:**
```
Ctrl+P  - Pause
Ctrl+R  - Resume
Ctrl+S  - Stop
Ctrl+Q  - Quit
```

**Pros:**
- Familiar to Vim users
- Less risk of accidental keypresses
- Clearer "command" vs "text input" distinction
- Standard in many terminal apps

**Cons:**
- Slower (requires Ctrl modifier)
- Higher cognitive load
- Less accessible (two-key combinations)
- Overkill for simple timer (we're not editing text)
- Ctrl+S often mapped to terminal flow control (XOFF)
- Ctrl+Q often mapped to terminal flow control (XON)

**Rejected:** Unnecessary complexity for timer application

---

### 3.2 Option B: Numeric Menu

**Approach:**
```
[1] Pause
[2] Resume
[3] Stop
[4] Quit
```

**Pros:**
- Unambiguous
- Easy to display
- Works well with long command lists
- Familiar from some TUI apps

**Cons:**
- Not mnemonic (have to remember numbers)
- Slower than letter keys
- Numbers less ergonomic than home row
- Doesn't scale well
- Less intuitive

**Rejected:** Less usable than mnemonic letter keys

---

### 3.3 Option C: Arrow Keys + Enter

**Approach:**
```
[↑/↓] Navigate menu
[Enter] Select action
```

**Pros:**
- Very explicit
- Hard to trigger accidentally
- Familiar from menu-driven apps

**Cons:**
- Much slower (multiple keypresses)
- Overkill for 3-4 actions
- Poor UX for timer (need quick pause)
- Requires menu navigation state
- Takes focus away from work

**Rejected:** Too slow for timer interactions

---

### 3.4 Option D: Single-Key, Case-Sensitive

**Approach:**
```
p - Pause
r - Resume
S - Stop (uppercase)
Q - Quit (uppercase)
```

**Pros:**
- More actions possible (26 + 26 = 52 keys)
- Could reserve uppercase for "dangerous" actions

**Cons:**
- Caps Lock confusion
- Harder to remember which is uppercase
- Accessibility issues
- We don't need 52 actions
- Error-prone

**Rejected:** Unnecessary complexity, worse UX

---

### 3.5 Option E: Mouse Support

**Approach:**
```
Click [Pause] button
Click [Resume] button
Click [Stop] button
```

**Pros:**
- Familiar to GUI users
- Visual feedback on hover
- Can show button states

**Cons:**
- Breaks keyboard-first philosophy
- Slower than keyboard
- Doesn't work over SSH easily
- Terminal mouse support inconsistent
- Against project goals (CLI-native)
- Mouse not always available in terminals

**Rejected:** Violates keyboard-first principle

## 4. Implementation

### 4.1 Control Legend Display

**Running State:**
```
[p] pause  [s] stop  [q] quit  [Ctrl+C] emergency exit
```

**Paused State:**
```
[r] resume  [s] stop  [q] quit  [Ctrl+C] emergency exit
```

**Display Rules:**
- Only show valid keys for current state
- Legend always visible at bottom of TUI
- Brief descriptions (verb form)
- Highlight key in brackets
- Ctrl+C always shown

### 4.2 Invalid Key Handling

**Behavior:**
- Invalid keys for current state: silently ignore
- No error messages or visual feedback
- Avoids clutter and distraction
- User can reference legend if unsure

**Rationale:**
- Users might press keys while thinking
- No harm in ignoring invalid input
- Keeps TUI clean and focused
- Legend provides guidance

## 5. Consequences

### 5.1 Positive

**Usability:**
- Fast, one-key actions
- Easy to learn and remember
- Forgiving (case-insensitive)
- Minimal cognitive load
- Accessible to all users

**Efficiency:**
- Pause instantly on interruption
- Resume quickly after break
- Stop without fumbling for keys
- No modifier key combinations

**Clarity:**
- Legend always visible
- Mnemonic keys (intuitive)
- Only relevant keys shown
- No ambiguity

### 5.2 Negative

**Accidental Keypresses:**
- Easy to trigger accidentally if terminal has focus
- Could pause timer while typing in another window
- Mitigated by: timer is foreground app, user watching it

**Limited Actions:**
- Single-key limits to ~26 actions (more than enough for timer)
- Can't add many more features without key conflicts
- Mitigated by: timer is simple, few actions needed

### 5.3 Edge Cases

**Terminal Flow Control:**
- Ctrl+S and Ctrl+Q often mapped to XOFF/XON
- Solution: Use plain 's' and 'q', not Ctrl versions
- Already handled in current design

**International Keyboards:**
- Keys like 'p', 'r', 's', 'q' universal on QWERTY
- May be in different positions on non-QWERTY layouts
- Mitigated by: mnemonic still works (letters are same)
- Future: Could support custom key bindings

## 6. Accessibility Considerations

### 6.1 Screen Reader Support

**Current Design:**
- Text-based TUI (screen reader compatible)
- Control legend is text (readable)
- State changes are text (readable)
- No mouse required

**Future Improvements:**
- ARIA-like labels for state changes
- Audio cues (via plugins)
- Configurable verbosity

### 6.2 Motor Impairments

**Benefits:**
- Single-key actions (no chords)
- Case-insensitive (no precision needed)
- No timing requirements
- No mouse needed (hard for some users)

### 6.3 Visual Impairments

**Benefits:**
- High contrast themes available
- Large, clear text
- Progress bar (multiple visual cues)
- Text-based (works with screen magnifiers)

## 7. User Guidance

### 7.1 Learning the Controls

**First-Time Users:**
1. Start timer: `pomodux start 25m "My task"`
2. See legend at bottom: `[p] pause  [s] stop...`
3. Try pressing 'p' → timer pauses
4. See legend updates: `[r] resume  [s] stop...`
5. Press 'r' → timer resumes

**Mental Model:**
- Look at legend for available actions
- Press the letter that makes sense
- Timer responds immediately
- Legend updates to show new options

### 7.2 Help Documentation

**README:**
```markdown
## Keyboard Controls

During a timer session:

| Key     | Action              |
|---------|---------------------|
| p       | Pause timer         |
| r       | Resume timer        |
| s, q    | Stop and exit       |
| Ctrl+C  | Emergency exit      |

All keys are case-insensitive. Available controls are shown at the bottom of the timer display.
```

## 8. Future Considerations

### 8.1 Custom Key Bindings

**Potential:**
```yaml
keybindings:
  pause: "p"
  resume: "r"
  stop: "s"
  quit: "q"
```

**Use Cases:**
- Users with different keyboard layouts
- Users preferring different mnemonics
- Users with muscle memory from other apps

**Complexity:**
- Config validation needed
- Conflict detection (no duplicate bindings)
- Documentation complexity
- Testing burden

**Evaluation:**
- Wait for user requests
- Assess if default bindings cause issues
- Consider if benefit outweighs complexity

### 8.2 Additional Actions

**Potential Future Actions:**
- `a` - Add time to timer
- `h` - Show help overlay
- `t` - Show timer statistics
- `n` - Add note to session

**Considerations:**
- Keep it simple (focus on core timer functionality)
- Use plugins for advanced features
- Don't add features without clear use cases
- Preserve single-key principle

