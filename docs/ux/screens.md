# Screen Specifications

**Version:** 1.0  
**Date:** 2026-01-27

## Overview

This document provides complete screen wireframes for all states of the Pomodux TUI. Each screen shows component composition, entry/exit conditions, and responsive behavior.

**Source:** [Requirements Section 8.1 Screen Layouts](../requirements/base.md#screen-layouts), [Components](components.md)

---

## Screen Layout Principles

### Layout Approach

**Centered Layout:**
- Window centered horizontally and vertically
- Content centered within window
- Responsive to terminal size

**Component Architecture:**
- **Main Window:** All components inline (border, header, progress, time, status, action selection)
- **No Overlays:** All components are part of main window content
- **Action Selection:** Always visible, changes content based on timer state

**Minimum Dimensions:**
- 80 columns minimum
- 24 rows minimum
- Graceful degradation below minimum

**Theme Application:**
- All colors from theme
- Border style from theme
- Consistent styling throughout

**Source:** [Requirements Section 8.1](../requirements/base.md#screen-layouts)

### Action Selection Architecture

**Inline Action Selection:**
- **Position:** Inside main window, below status indicator
- **Always Visible:** No fading, always displayed
- **State-Based Content:** Content changes based on timer state
- **Confirmation Inline:** Confirmation prompt appears within action selection component
- **Rationale:** Simpler than overlay approach, always visible, no fade complexity

---

## Screen States

### State 1: Running

**Purpose:** Display active timer counting down

**Wireframe:**
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header
│                                                              │
│  ████████████████████████░░░░░░░░░░░░░  60%   15:23          │ ← Progress Bar + Time
│                                                              │
│  Status: RUNNING                                             │ ← Status Indicator
│                                                              │
│  [p]ause  [s]top                                             │ ← Action Selection (always visible)
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window)
```

**Component Composition:**

**Main Window (Persistent):**
1. **Window/Border** - Rounded border, theme border color
2. **Session Header** - "Work Session: Implementing authentication"
3. **Progress Bar** - 60% filled, shows percentage, theme colors
4. **Time Display** - "15:23" (MM:SS format)
5. **Status Indicator** - "RUNNING" (green/success color)
6. **Action Selection** - `[p]ause  [s]top` (always visible, no fade)

**Entry Conditions:**
- Timer started successfully
- Timer state is "running"
- TUI initialized

**Exit Conditions:**
- User presses `p` → Transition to Paused
- User presses `q` (or `s` alias) → Transition to Confirmation
- Timer reaches 0:00 → Transition to Completed
- User presses `Ctrl+C` → Exit (interrupted)

**Responsive Behavior:**
- Window width adapts to terminal (centered, max content width)
- Progress bar width adapts (max 80 columns)
- Layout recalculates on terminal resize
- If terminal too small: Show warning, continue with degraded layout

**Requirements References:**
- [Section 8.1.1 Running Timer State](../requirements/base.md#running-timer-state)
- [FR-TIMER-001](../requirements/base.md#fr-timer-001) - Timer Start

**Flow Reference:** [Flow 1: Start Timer](application-flows.md#flow-1-start-timer-with-preset)

---

### State 2: Paused

**Purpose:** Display paused timer

**Wireframe:**
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header
│                                                              │
│  ████████████████████████░░░░░░░░░░░░░  60%   15:23          │ ← Progress Bar (frozen) + Time (frozen)
│                                                              │
│  Status: ⏸ PAUSED                                            │ ← Status Indicator
│                                                              │
│  [r]esume  [s]top                                            │ ← Action Selection (always visible)
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window)
```

**Component Composition:**

**Main Window (Persistent):**
1. **Window/Border** - Same as Running
2. **Session Header** - Same as Running
3. **Progress Bar** - Frozen at 60% (no updates)
4. **Time Display** - Frozen at "15:23" (no countdown)
5. **Status Indicator** - "⏸ PAUSED" (yellow/warning color)
6. **Action Selection** - `[r]esume  [s]top` (always visible)

**Changes from Running:**
- Status: "⏸ PAUSED" (yellow) instead of "RUNNING" (green)
- Progress bar: No animation (frozen)
- Time display: No countdown (frozen)
- Action Selection: Content changes to `[r]esume  [s]top`

**Entry Conditions:**
- Timer is running
- User presses `p` key
- Timer state changes to "paused"

**Exit Conditions:**
- User presses `r` → Transition to Running
- User presses `q` (or `s` alias) → Transition to Confirmation
- User presses `Ctrl+C` → Exit (interrupted)

**Responsive Behavior:**
- Same as Running state

**Requirements References:**
- [Section 8.1.2 Paused Timer State](../requirements/base.md#812-paused-timer-state)
- [FR-TIMER-002](../requirements/base.md#fr-timer-002) - Timer Pause/Resume

**Flow Reference:** [Flow 3: Pause and Resume](application-flows.md#flow-3-pause-and-resume-timer)

---

### State 3: Completed

**Purpose:** Display completion confirmation (brief)

**Wireframe:**

**Initial (Closing in 3):**
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header
│                                                              │
│  ███████████████████████████████████████████████  100%  0:00 │ ← Progress Bar (100%) + Time (0:00)
│                                                              │
│  Status: ✓ COMPLETED                                         │ ← Status Indicator
│                                                              │
│  Session saved! Closing in 3.                                │ ← Completion Message (countdown: 3)
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window)
```

**After 1 second (Closing in 2):**
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header
│                                                              │
│  ███████████████████████████████████████████████  100%  0:00 │ ← Progress Bar (100%) + Time (0:00)
│                                                              │
│  Status: ✓ COMPLETED                                         │ ← Status Indicator
│                                                              │
│  Session saved! Closing in 2.                                │ ← Completion Message (countdown: 2)
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window)
```

**After 2 seconds (Closing in 1):**
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header
│                                                              │
│  ███████████████████████████████████████████████  100%  0:00 │ ← Progress Bar (100%) + Time (0:00)
│                                                              │
│  Status: ✓ COMPLETED                                         │ ← Status Indicator
│                                                              │
│  Session saved! Closing in 1.                                │ ← Completion Message (countdown: 1)
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window)
```

**Component Composition:**

**Main Window (Persistent):**
1. **Window/Border** - Same as Running
2. **Session Header** - Same as Running
3. **Progress Bar** - 100% filled, static
4. **Time Display** - "0:00"
5. **Status Indicator** - "✓ COMPLETED" (green/success color)
6. **Completion Message** - "Session saved! Closing in 3." (countdown: 3 → 2 → 1)

**Changes from Running:**
- Progress bar: 100% filled
- Time: "0:00"
- Status: "✓ COMPLETED" (green)
- Message: "Session saved! Closing in 3." appears, counts down to 1
- Action Selection: Not shown (exits after countdown)

**Entry Conditions:**
- Timer is running
- Timer reaches 0:00
- Timer state changes to "completed"

**Exit Conditions:**
- After countdown completes (~3 seconds) → Exit automatically
- No user interaction needed

**Display Duration:**
- ~3 seconds total
- Countdown updates every 1 second: "Closing in 3." → "Closing in 2." → "Closing in 1."
- Then automatic exit

**Requirements References:**
- [Section 8.1.3 Completion State](../requirements/base.md#813-completion-state)
- [FR-TIMER-004](../requirements/base.md#fr-timer-004) - Timer Completion
- [UXDR: Auto-Exit Behavior](../uxdr/auto-exit-behavior.md)

**Flow Reference:** [Flow 5: Timer Completion](application-flows.md#flow-5-timer-completion)

---

### State 4: Confirmation

**Purpose:** Confirm user intent before stopping timer (inline confirmation)

**Wireframe:**
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header
│                                                              │
│  ████████████████████████░░░░░░░░░░░░░  60%   15:23          │ ← Progress Bar (frozen) + Time (frozen)
│                                                              │
│  Status: ⏸ PAUSED                                            │ ← Status Indicator (paused automatically)
│                                                              │
│  Stop timer and exit? [y]es / [n]o                           │ ← Action Selection (confirmation prompt)
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window)
```

**Component Composition:**

**Main Window (Persistent):**
1. **Window/Border** - Same as Running/Paused
2. **Session Header** - Same as Running/Paused
3. **Progress Bar** - Frozen at current progress (no updates)
4. **Time Display** - Frozen at current time (no countdown)
5. **Status Indicator** - "⏸ PAUSED" (timer paused automatically)
6. **Action Selection** - "Stop timer and exit? [y]es / [n]o" (confirmation prompt inline)

**Entry Conditions:**
- Timer is running or paused
- User presses `q` or `s` key

**Exit Conditions:**
- User presses `y` or `Y` → Confirm stop, transition to Stopped, then exit
- User presses `n`, `N`, or `Esc` → Cancel, return to previous state (timer resumes if was running)
- User presses `Ctrl+C` → Emergency exit (bypasses confirmation)

**Keyboard Controls:**
- `y` / `Y`: Confirm stop and exit
- `n` / `N`: Cancel, return to timer
- `Esc`: Cancel (same as `n`)
- Other keys: Ignored (wait for y/n)

**Behavior:**
- Timer is automatically paused when confirmation appears (if running, pauses; if paused, stays paused)
- Confirmation prompt appears inline within Action Selection component
- No overlay - confirmation happens within main window
- Timer does not continue running (paused during confirmation)
- On cancel (`n`/`Esc`): Timer resumes (unpauses and continues from where it was) if it was running
- Action Selection content changes back to state-appropriate actions after cancel
- User must explicitly confirm or cancel

**Requirements References:**
- [US-1.4](../requirements/base.md#us-14) - Stop Timer Early
- [FR-TIMER-003](../requirements/base.md#fr-timer-003) - Timer Stop

**Flow Reference:** [Flow 4: Stop Timer Early](application-flows.md#flow-4-stop-timer-early)

---

## Error States

### Error State 1: Terminal Too Small

**Purpose:** Warn user if terminal below minimum size

**Wireframe:**
```
┌─────────────────────────────────┐
│ Terminal too small!             │
│ Minimum: 80x24                  │
│ Current: 60x20                  │
└─────────────────────────────────┘
```

**Component Composition:**
1. **Terminal Size Warning** - Warning message, centered

**Display Conditions:**
- Terminal width < 80 columns OR
- Terminal height < 24 rows
- Warning displayed, continue with degraded layout OR exit

**Resolution:**
- User resizes terminal to adequate size
- Warning disappears when size adequate

**Requirements References:**
- [Section 8.3 Terminal Resize Handling](../requirements/base.md#terminal-resize-handling)

**Flow Reference:** [Error Flow 5: Terminal Too Small](application-flows.md#error-flow-5-terminal-too-small)

---

### Config/Theme Errors (No TUI)

**Purpose:** Clarify that config or theme load errors prevent the TUI from starting.

**Behavior:**
- Config load failure (e.g. invalid YAML, unreadable file) or theme resolution failure (e.g. unknown theme name) causes startup to fail.
- The timer does not load; the TUI is not shown.
- The application returns an error to the CLI and exits. No in-TUI banner or error state.

**Requirements References:**
- [NFR-REL-003](../requirements/base.md#nfr-rel-003) - Config Validation
- [FR-CONFIG-001](../requirements/base.md#fr-config-001) - Config File Loading

**Flow Reference:** [Error Flow 4: Config/Theme Load Failure](application-flows.md#error-flow-4-configtheme-load-failure)

---

## Responsive Behavior

### Terminal Resize Handling

**Behavior:**
- Detect resize via `tea.WindowSizeMsg`
- Recalculate layout immediately
- Maintain timer accuracy (no interruption)
- Adapt component widths/heights

**Minimum Size:**
- 80 columns x 24 rows minimum
- Below minimum: Show warning, continue with degraded layout

**Component Adaptation:**
- Window: Centers within terminal
- Progress bar: Width adapts (max 80 columns)
- Text: Wraps if needed (unlikely for short labels)

**Source:** [Requirements Section 8.3](../requirements/base.md#terminal-resize-handling), [ADR: TUI Framework](../adr/tui-framework-selection.md)

---

## Theme Application

### Theme Integration

**All Components Use Theme:**
- Window border: Theme border color and style
- Text colors: Theme foreground, primary, muted colors
- Progress bar: Theme progress colors
- Status indicator: Theme semantic colors (success/warning/error)

**Theme Selection:**
- From config: `theme: "catppuccin-mocha"`
- Applied to all UI elements
- Consistent throughout application

**Source:** [Section 8.4 Theme Application](../requirements/base.md#theme-application), [Components: Theme Integration](components.md#theme-integration)

---

## Entry/Exit Conditions Summary

### Entry Conditions

**Running State:**
- Timer started successfully
- Timer state is "running"

**Paused State:**
- Timer is running
- User presses `p`

**Completed State:**
- Timer is running
- Timer reaches 0:00

**Error States:**
- Terminal too small: Terminal size < 80x24

**Note:** Config or theme load errors prevent the TUI from starting (timer does not load); there is no TUI screen for that case.

### Exit Conditions

**Running State:**
- User presses `p` → Paused
- User presses `s`/`q` → Exit (stopped)
- Timer completes → Completed
- User presses `Ctrl+C` → Exit (interrupted)

**Paused State:**
- User presses `r` → Running
- User presses `s`/`q` → Exit (stopped)
- User presses `Ctrl+C` → Exit (interrupted)

**Completed State:**
- After ~500ms → Exit automatically

---

## References

### Requirements
- [Section 8.1 Screen Layouts](../requirements/base.md#screen-layouts)
- [Section 8.3 Terminal Resize Handling](../requirements/base.md#terminal-resize-handling)
- [Section 8.4 Theme Application](../requirements/base.md#theme-application)

### Components
- [Component Specifications](components.md)
- [Component: Window/Border](components.md#component-1-windowborder)
- [Component: Progress Bar](components.md#component-3-progress-bar)
- [Component: Status Indicator](components.md#component-5-status-indicator)

### Flows
- [Application Flows](application-flows.md)
- [Flow 1: Start Timer](application-flows.md#flow-1-start-timer-with-preset)
- [Flow 3: Pause and Resume](application-flows.md#flow-3-pause-and-resume-timer)
- [Flow 5: Timer Completion](application-flows.md#flow-5-timer-completion)

### Research

---

**Last Updated:** 2026-01-28
