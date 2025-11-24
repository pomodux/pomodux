---
status: approved
---

# Auto-Exit Behavior After Timer Completion

## 1. Context / Background

### 1.1 Problem Statement

When a timer completes (reaches 0:00), the application must decide what to do next:
- Exit immediately back to the command line
- Wait for user acknowledgment (keypress)
- Auto-exit after a delay (configurable timeout)

This decision affects the user's workflow and attention management.

### 1.2 User Impact

**Primary Use Case:**
Users start a timer and focus on their work. When the timer completes, they need to be notified and decide what to do next (start break, start new session, etc.).

**Workflow Considerations:**
- Users may not be watching the terminal when timer completes
- Users want minimal friction to start next timer
- Users may want to acknowledge completion before dismissing
- Users can check history anytime with `pomodux-stats`

## 2. Decision

**Selected Solution:** Option C - Show completion message, then exit immediately

### 2.1 Rationale

**Why Immediate Exit?**

1. **Clean Command Line Return**: User gets their prompt back immediately, ready for next command
2. **Minimal Friction**: No need to dismiss notification, just run next timer command
3. **History Preservation**: Session is saved and can be reviewed with `pomodux-stats`
4. **Notification Alternatives**: Users wanting acknowledgment can use notification plugins
5. **Unix Philosophy**: Do one thing (run timer) and exit cleanly

**User Experience Flow:**
```
$ pomodux start work "Code review"
[TUI shows running timer for 25 minutes]
[Timer completes - shows "✓ COMPLETED" for 1 second]
[Exits to shell]
$ pomodux-stats --today
[Shows completed session]
$ pomodux start break
```

## 3. Alternatives Considered

### 3.1 Option A: Auto-Exit After N Seconds (Configurable)

**Approach:**
```
[Timer completes]
[Shows completion message]
[Waits 3-5 seconds (configurable)]
[Exits automatically]
```

**Pros:**
- Gives user time to notice completion
- Still exits automatically
- Configurable for user preference

**Cons:**
- Arbitrary delay feels awkward
- Users must wait before getting prompt back
- No clear benefit over immediate exit
- Adds configuration complexity

**Rejected:** Delay doesn't provide meaningful value, feels arbitrary

---

### 3.2 Option B: Wait for User Keypress

**Approach:**
```
[Timer completes]
[Shows completion message]
Status: ✓ COMPLETED
Press any key to continue...
[Waits indefinitely for keypress]
```

**Pros:**
- User explicitly acknowledges completion
- Can review completion stats before dismissing
- Familiar pattern from some applications

**Cons:**
- Requires user to return to terminal and press key
- Blocks command line until acknowledged
- Extra friction for starting next timer
- Doesn't match typical Unix tool behavior
- User forgets and terminal stays blocked

**Rejected:** Too much friction, blocks user workflow

---

### 3.3 Option C: Immediate Exit (Selected)

**Approach:**
```
[Timer completes]
[Shows "✓ COMPLETED" for ~500ms]
[Exits immediately]
[User gets command prompt back]
```

**Pros:**
- Clean return to command line
- No waiting or extra keypress needed
- User can immediately start next command
- History preserved in `pomodux-stats`
- Matches Unix tool expectations
- Notification handled by plugins if desired

**Cons:**
- User might miss completion if not watching
- No visual acknowledgment requirement
- Cannot review stats before dismissing (but can use `pomodux-stats`)

**Selected:** Best balance of simplicity and workflow efficiency

## 4. Consequences

### 4.1 Positive

**Workflow Efficiency:**
- Users can chain timer commands quickly
- No blocked terminal sessions
- Clean command line experience
- Follows Unix philosophy

**Flexibility:**
- Users wanting notifications can use plugins
- History always available via `pomodux-stats`
- Scriptable and automatable
- Predictable behavior

### 4.2 Negative

**Potential Missed Completions:**
- Users not watching terminal might miss completion
- No forced acknowledgment
- Relies on external notifications (terminal bell, plugins)

**Mitigation:**
- Terminal bell option (configurable, default: off)
- Notification plugins (desktop notifications)
- Completion message visible briefly
- Session logged for later review

### 4.3 Implementation Details

**Completion Sequence:**
1. Timer reaches 0:00
2. Update display to show "✓ COMPLETED"
3. Update progress bar to 100%
4. Display "Session saved!" message
5. Wait 500-1000ms for visual feedback
6. Exit TUI with status code 0
7. Return to shell prompt

**Future Enhancements:**
- Configuration option to enable "wait for keypress" mode (if users request it)
- Plugin hook for completion actions before exit
- Configurable completion display duration

## 5. User Guidance

### 5.1 Recommended Workflow

**For Users Who Want Notifications:**
```yaml
# config.yaml
timer:
  bell_on_complete: true

plugins:
  enabled:
    - notify
```

**For Users Who Want to Review Stats:**
```bash
# Complete timer
pomodux start work "Task"
# [exits on completion]

# Review immediately
pomodux-stats --today
```

**For Rapid Timer Chains:**
```bash
# Work session
pomodux start work "Feature A" && \
# Break
pomodux start break && \
# Next work session
pomodux start work "Feature B"
```

### 5.2 Documentation Notes

**README should highlight:**
- Timer exits immediately on completion
- Use `pomodux-stats` to review session history
- Enable terminal bell or notification plugins for alerts
- Session data is always preserved

