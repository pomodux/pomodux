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
- **Main Window:** Persistent components (border, header, progress, time, status)
- **Overlays:** Transient components (control legend) that fade after initial display
- **Separation:** Overlays rendered separately, positioned below main window

**Minimum Dimensions:**
- 80 columns minimum
- 24 rows minimum
- Graceful degradation below minimum

**Theme Application:**
- All colors from theme
- Border style from theme
- Consistent styling throughout

**Source:** [Requirements Section 8.1](../requirements/base.md#screen-layouts)

### Overlay Architecture

**Control Legend Overlay (Option 1: Below Window, Centered):**
- **Position:** Below main window border, centered horizontally
- **Initial Display:** Appears when timer starts
- **Fade Timing:** Visible 3-5 seconds, then fades over ~1 second
- **Fade Method:** Color-based dimming (4-5 steps, ~200-250ms per step)
- **Behavior:** Independent of main window, fades without affecting window content
- **Layout Impact:** None - window stays fixed, legend fades independently
- **Rationale:** Provides initial guidance, then fades to reduce clutter

**Benefits:**
- **No Layout Impact:** Window dimensions stay fixed during fade
- **Smooth Fade:** Color-based dimming provides smooth visual transition
- **Cleaner Main Window:** No persistent legend cluttering the interface
- **Initial Guidance:** Shows controls when users need them most
- **Minimal UI:** After fade, interface is clean and focused
- **Independent Behavior:** Overlay fades without affecting other components

---

## Screen States

### State 1: Running

**Purpose:** Display active timer counting down

**Wireframe:**

**Initial State (Control Legend Visible):**
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header
│                                                              │
│  ████████████████████████░░░░░░░░░░░░░  60%   15:23          │ ← Progress Bar + Time
│                                                              │
│  Status: RUNNING                                             │ ← Status Indicator
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window)
                                                              │
                                                              │
         [p] pause  [q] quit                                  ← Overlay: Control Legend (fades after 3-5s)
```

**After Fade (Control Legend Hidden):**
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header
│                                                              │
│  ████████████████████████░░░░░░░░░░░░░  60%   15:23          │ ← Progress Bar + Time
│                                                              │
│  Status: RUNNING                                             │ ← Status Indicator
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window)
                                                              │
                                                              │
                                                              ← Overlay: Control Legend (faded, not displayed)
```

**Component Composition:**

**Main Window (Persistent):**
1. **Window/Border** - Rounded border, theme border color
2. **Session Header** - "Work Session: Implementing authentication"
3. **Progress Bar** - 60% filled, shows percentage, theme colors
4. **Time Display** - "15:23" (MM:SS format)
5. **Status Indicator** - "RUNNING" (green/success color)

**Overlay (Transient):**
6. **Control Legend** - Minimal display: `[p] pause [q] quit` - appears below window, fades after 3-5 seconds

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

**Paused State (Legend Always Visible):**
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header
│                                                              │
│  ████████████████████████░░░░░░░░░░░░░  60%   15:23          │ ← Progress Bar (frozen) + Time (frozen)
│                                                              │
│  Status: ⏸ PAUSED                                            │ ← Status Indicator
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window)
                                                              │
                                                              │
         [r] resume  [q] quit                                 ← Overlay: Control Legend (always visible when paused)
```

**Component Composition:**

**Main Window (Persistent):**
1. **Window/Border** - Same as Running
2. **Session Header** - Same as Running
3. **Progress Bar** - Frozen at 60% (no updates)
4. **Time Display** - Frozen at "15:23" (no countdown)
5. **Status Indicator** - "⏸ PAUSED" (yellow/warning color)

**Overlay (Transient):**
6. **Control Legend** - Always visible when paused, shows `[r] resume [q] quit`

**Changes from Running:**
- Status: "⏸ PAUSED" (yellow) instead of "RUNNING" (green)
- Progress bar: No animation (frozen)
- Time display: No countdown (frozen)
- Control legend: Reappears immediately, shows `[r] resume` instead of `[p] pause` (always visible when paused)

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
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header
│                                                              │
│  ███████████████████████████████████████████████  100%  0:00 │ ← Progress Bar (100%) + Time (0:00)
│                                                              │
│  Status: ✓ COMPLETED                                         │ ← Status Indicator
│                                                              │
│  Session saved!                                              │ ← Completion Message
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window)
                                                              │
                                                              │
                                                              ← Overlay: Control Legend (not shown, already faded)
```

**Component Composition:**

**Main Window (Persistent):**
1. **Window/Border** - Same as Running
2. **Session Header** - Same as Running
3. **Progress Bar** - 100% filled, static
4. **Time Display** - "0:00"
5. **Status Indicator** - "✓ COMPLETED" (green/success color)
6. **Completion Message** - "Session saved!" (green/success color)

**Overlay (Transient):**
7. **Control Legend** - Not shown (exits immediately, legend already faded)

**Changes from Running:**
- Progress bar: 100% filled
- Time: "0:00"
- Status: "✓ COMPLETED" (green)
- Message: "Session saved!" appears
- Controls: Not shown (exits after ~500ms)

**Entry Conditions:**
- Timer is running
- Timer reaches 0:00
- Timer state changes to "completed"

**Exit Conditions:**
- After ~500ms display → Exit automatically
- No user interaction needed

**Display Duration:**
- ~500ms (2 tick cycles at 250ms)
- Then automatic exit

**Requirements References:**
- [Section 8.1.3 Completion State](../requirements/base.md#813-completion-state)
- [FR-TIMER-004](../requirements/base.md#fr-timer-004) - Timer Completion
- [UXDR: Auto-Exit Behavior](../uxdr/auto-exit-behavior.md)

**Flow Reference:** [Flow 5: Timer Completion](application-flows.md#flow-5-timer-completion)

---

### State 4: Confirmation

**Purpose:** Confirm user intent before stopping timer

**Wireframe:**

**Confirmation Overlay (takes over entire timer):**
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header (dimmed/background)
│                                                              │
│  ████████████████████████░░░░░░░░░░░░░  60%   15:23          │ ← Progress Bar + Time (frozen, dimmed)
│                                                              │
│  Status: ⏸ PAUSED                                            │ ← Status Indicator (paused automatically)
│                                                              │
│                                                              │
│                    ┌─────────────────────────────┐           │
│                    │ Stop timer and exit?        │           │ ← Confirmation Dialog
│                    │ [y]es / [n]o               │           │   (Overlays entire window)
│                    └─────────────────────────────┘           │
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window, dimmed)
```

**Component Composition:**

**Main Window (Dimmed/Background):**
- Timer is automatically paused when confirmation appears
- Status shows "⏸ PAUSED" (timer paused automatically)
- Progress bar and time display frozen
- All components dimmed/backgrounded (visual indication of modal overlay)
- Timer does not continue running (paused)

**Overlay (Modal Dialog - Takes Over Entire Timer):**
- **Confirmation Dialog** - Centered overlay with question
- **Question:** "Stop timer and exit? [y]es / [n]o"
- **Position:** Centered horizontally and vertically, overlays entire main window
- **Modal:** Blocks other input until confirmed or cancelled
- **Visual:** Main window dimmed/backgrounded, dialog prominent

**Entry Conditions:**
- Timer is running or paused
- User presses `q` or `s` key

**Exit Conditions:**
- User presses `y` or `Y` → Confirm stop, transition to Stopped, then exit
- User presses `n`, `N`, or `Esc` → Cancel, timer resumes (unpauses and continues running)
- User presses `Ctrl+C` → Emergency exit (bypasses confirmation)

**Keyboard Controls:**
- `y` / `Y`: Confirm stop and exit
- `n` / `N`: Cancel, return to timer
- `Esc`: Cancel (same as `n`)
- Other keys: Ignored (wait for y/n)

**Behavior:**
- Timer is automatically paused when confirmation appears (if running, pauses; if paused, stays paused)
- Confirmation dialog is modal (blocks other input, takes over entire timer)
- Main window remains visible but dimmed/backgrounded (visual indication of modal)
- Timer does not continue running (paused during confirmation)
- On cancel (`n`/`Esc`): Timer resumes (unpauses and continues from where it was)
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

### Error State 2: Config Errors

**Purpose:** Display configuration validation errors

**Wireframe:**
```
┌─ Pomodux Timer ─────────────────────────────────────────────┐
│                                                              │
│  ⚠ Config errors detected, using defaults                   │ ← Config Error Banner
│                                                              │
│  Work Session: Implementing authentication                   │ ← Session Header
│                                                              │
│  ████████████████████████░░░░░░░░░░░░░  60%   15:23          │ ← Progress Bar + Time
│                                                              │
│  Status: RUNNING                                             │ ← Status Indicator
│                                                              │
└──────────────────────────────────────────────────────────────┘ ← Window/Border (Main Window)
                                                              │
                                                              │
         [p] pause  [q] quit                                  ← Overlay: Control Legend (if still visible, fades after 3-5s)
```

**Component Composition:**

**Main Window (Persistent):**
1. **Window/Border** - Same as Running
2. **Config Error Banner** - Warning message at top (always visible)
3. **Session Header** - Same as Running
4. **Progress Bar** - Same as Running
5. **Time Display** - Same as Running
6. **Status Indicator** - Same as Running

**Overlay (Transient):**
7. **Control Legend** - If still visible, shows `[p] pause [q] quit` (typically already faded)

**Display Conditions:**
- Config file has validation errors
- Invalid YAML or invalid values
- Application uses defaults for invalid fields

**Behavior:**
- Warning banner always visible (cannot dismiss)
- Timer continues normally with defaults
- No blocking of functionality

**Requirements References:**
- [NFR-REL-003](../requirements/base.md#nfr-rel-003) - Config Validation
- [FR-CONFIG-001](../requirements/base.md#fr-config-001) - Config File Loading

**Flow Reference:** [Error Flow 4: Config File Errors](application-flows.md#error-flow-4-config-file-errors)

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
- Config errors: Config validation fails

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

**Last Updated:** 2026-01-27
