# Interaction Specifications

**Version:** 1.0  
**Date:** 2026-01-27

## Overview

This document specifies all keyboard controls, focus management, input handling, and transition behaviors for the Pomodux TUI.

**Source:** [Requirements Section 8.2 Keyboard Controls](../requirements/base.md#keyboard-controls), [UXDR: Keyboard Controls Design](../uxdr/keyboard-controls-design.md), [Components](components.md)

---

## Keyboard Controls

### Action Classification

**Primary Actions:**
- Core workflow actions used during normal timer operation
- Prominently displayed in control legend
- Single-key, mnemonic controls
- Part of expected user journey
- Examples: Pause (`p`), Resume (`r`), Stop (`s`/`q`)

**Secondary Actions:**
- Emergency or exceptional actions
- De-emphasized in control legend
- Available but not part of normal workflow
- Used for exceptional circumstances
- Example: Emergency exit (`Ctrl+C`)

**Rationale for Classification:**
- **Primary actions** are the main controls users interact with during normal timer use
- **Secondary actions** are emergency/exceptional controls that should be available but not prominent
- Visual hierarchy guides users to primary actions while keeping secondary actions accessible
- Reduces cognitive load by emphasizing what users need most often
- Follows UX principle: Make primary actions prominent, secondary actions available but not distracting

### Control Reference Table

#### Primary Actions

| Key | Action | Available States | Mnemonic |
|-----|--------|------------------|----------|
| `p` | Pause timer | Running | **P**ause |
| `r` | Resume timer | Paused | **R**esume |
| `q` | Stop timer (shows confirmation) | Running, Paused | **Q**uit |
| `s` | Stop timer (alias, shows confirmation) | Running, Paused | **S**top |
| `y` | Confirm stop | Confirmation | **Y**es |
| `n` | Cancel stop | Confirmation | **N**o |
| `Esc` | Cancel stop (same as `n`) | Confirmation | Escape |

**Note:** `s` is also available as an alias for stop (same action as `q`), but only `q` is displayed in the legend for minimal UI.

#### Secondary Actions

| Key | Action | Available States | Type |
|-----|--------|------------------|------|
| `Ctrl+C` | Emergency exit | All states | Emergency (not displayed - universal terminal convention) |
| `s` | Stop timer (alias) | Running, Paused | Alias (not displayed - `q` shown instead) |

**Source:** [Requirements Section 8.2](../requirements/base.md#keyboard-controls), [UXDR: Keyboard Controls Design](../uxdr/keyboard-controls-design.md)

---

### Control Behavior

**Single-Key Actions:**
- All controls are single-key presses (no chords except Ctrl+C)
- Fast, immediate response
- Low cognitive load

**Case-Insensitive:**
- Keys work in both uppercase and lowercase (`p` = `P`)
- Forgiving of Caps Lock mistakes
- Accessible to all users

**State-Based Controls:**
- Controls available based on current timer state
- Invalid keys in current state are silently ignored
- Control legend updates to show available controls

**Invalid Key Handling:**
- Invalid keys for current state: Silently ignored
- No error messages or visual feedback
- Avoids clutter and distraction
- User can reference legend if unsure

**Source:** [UXDR: Keyboard Controls Design](../uxdr/keyboard-controls-design.md#22-rationale)

---

### Control Legend (Overlay)

**Display Rules:**
- **Initial Display:** Appears below main window when timer starts
- **Fade Behavior:** Visible for 3-5 seconds, then fades out over ~1 second
- **Position:** Below main window border, centered
- **Component Type:** Overlay (not part of main window content)
- Shows only valid primary keys for current state
- Brief descriptions (verb form)
- Keys highlighted in brackets

**Visual Hierarchy:**
- Primary actions: Normal/muted text color
- Minimal display: Only essential controls shown
- Secondary actions: Not displayed (Ctrl+C is universal terminal convention, `s` is alias)

**Running State Legend:**
```
[p] pause  [q] quit
```

**Paused State Legend:**
```
[r] resume  [q] quit
```

**Fade Behavior:**
- **Running State:** 
  - Visible for 3-5 seconds after timer starts
  - Fades out over ~1 second using color-based dimming (4-5 steps)
  - Stays hidden during active countdown
- **Paused State:**
  - Legend reappears immediately when paused
  - Always visible while paused (does not fade)
  - Shows updated controls: `[r] resume [q] quit`
- **Resume to Running:**
  - Legend remains visible for 3-5 seconds after resume
  - Then fades out again using same fade mechanism
- **Fade Duration:** ~1 second fade out (4-5 color steps)
- **Fade Method:** Color brightness reduction (TUI doesn't support true opacity)
- **Layout Impact:** None - window stays fixed, legend fades independently
- **Rationale:** Paused state is interactive, so showing controls provides helpful context

**Rationale for Fading:**
- **Initial Guidance:** Shows controls when user first sees timer
- **Reduces Clutter:** Fades away to keep UI minimal and focused
- **User Memory:** After seeing once, users remember controls
- **Clean Interface:** Timer display becomes cleaner after initial period
- **Overlay Architecture:** Separate from main window, allows independent fade behavior

**Rationale for Minimal Display:**
- **Only show `q` for stop:** `s` and `q` are aliases (same action), showing both is redundant
- **Don't show Ctrl+C:** Universal terminal convention, users already know it works
- **Minimal UI principle:** Show only what's necessary, reduce visual clutter
- **Both `s` and `q` still work:** Just not displayed (users can use either)
- **Ctrl+C still works:** Universal terminal convention, doesn't need to be shown

**Completed State:**
- Legend not shown (exits immediately)

**Legend Updates:**
- **Timer Start:** Legend appears when timer starts (overlay below window)
- **Running State:** Legend fades after 3-5 seconds, stays hidden
- **Pause Action:** Legend reappears immediately when paused (becomes visible)
- **Paused State:** Legend always visible (does not fade while paused)
- **Resume Action:** Legend remains visible, then fades after 3-5 seconds
- **Content Updates:** Legend text changes immediately on state transition (`p` vs `r`)
- **Stop Action:** `q` remains constant (available in both states)

**Rationale:**
- **Initial Guidance:** Shows controls when timer starts, provides orientation
- **Fades Away:** Reduces visual clutter after initial period
- **Minimal UI:** Show only essential, state-specific controls
- **Reduce redundancy:** Don't show aliases (`s` works but not displayed)
- **Universal conventions:** Don't show Ctrl+C (users already know it)
- **Focus on workflow:** Display only what changes with state (`p`/`r`) plus constant stop (`q`)
- **Less is more:** Fewer displayed controls = cleaner, less cluttered interface
- **Overlay Architecture:** Separate component allows fade without affecting main window

**Source:** [UXDR: Keyboard Controls Design](../uxdr/keyboard-controls-design.md#41-control-legend-display), [Component: Control Legend](components.md#component-6-control-legend)

---

## Focus Management

### Focus Model

**Single Focus:**
- TUI has single focus (the timer display)
- No focus rings or focus indicators needed
- All keyboard input goes to timer controls
- No tab navigation or focus movement

**Focus Behavior:**
- Focus is always on timer (no focus changes)
- Keyboard input always processed as timer controls
- No focus management needed (single focus)

**Rationale:**
- Simple TUI with single purpose
- No multiple interactive elements
- Keyboard controls are global (not element-specific)

**Source:** [Requirements NFR-USE-002](../requirements/base.md#nfr-use-002) - Keyboard Accessibility

---

## Input Handling

### Keypress Event Handling

**Event Flow:**
1. User presses key
2. Bubble Tea receives `tea.KeyMsg`
3. TUI `Update()` function processes key
4. Timer state updated (if valid key)
5. TUI redraws with new state

**Key Processing:**
- Keys processed synchronously in `Update()` function
- State changes are immediate
- UI updates on next render cycle

**Invalid Key Handling:**
- Invalid keys: Silently ignored (no error)
- Invalid keys for state: Silently ignored
- No visual feedback for invalid keys
- Legend provides guidance

**Source:** [ADR: TUI Framework Selection](../adr/tui-framework-selection.md) - Bubble Tea Event Loop

---

### Key Event Propagation

**Event Handling:**
- All keys processed in TUI `Update()` function
- No event bubbling or propagation
- Direct key-to-action mapping
- No intermediate handlers

**Key Mapping:**
- `p` → Pause action (if running)
- `r` → Resume action (if paused)
- `q` → Stop action (if running/paused) - displayed in legend
- `s` → Stop action (if running/paused) - alias, works but not displayed
- `Ctrl+C` → Interrupt action (always) - works but not displayed (universal convention)

**Source:** [Component: Control Legend](components.md#component-6-control-legend)

---

## Transitions

### State Transitions

**Transition Triggers:**
- **Running → Paused:** User presses `p`
- **Paused → Running:** User presses `r`
- **Running → Completed:** Timer reaches 0:00
- **Running/Paused → Confirmation:** User presses `q` or `s` (timer paused, shows confirmation dialog)
- **Confirmation → Running:** User presses `n`, `N`, or `Esc` (cancel, timer resumes)
- **Confirmation → Stopped:** User presses `y` or `Y` (confirm stop)
- **Any → Interrupted:** User presses `Ctrl+C` (bypasses confirmation)

**Transition Timing:**
- Immediate (no delay)
- State changes synchronously
- UI updates on next render cycle (~250ms)

**Transition Feedback:**
- Visual feedback: Status indicator changes
- Progress bar: Freezes (pause) or continues (resume)
- Control legend: Updates immediately
- No animation delays

**Source:** [State Machines](state-machines.md#application-state-machine), [Application Flows](application-flows.md#screen-flows)

---

### Visual Transitions

**Progress Bar Updates:**
- Updates every 250ms (via tick message)
- Smooth progress visualization
- No animation delays
- Immediate freeze on pause

**Status Indicator Updates:**
- Updates immediately on state change
- Text and color change synchronously
- No fade or animation
- Instant feedback

**Control Legend Updates:**
- **Running State:** Fades after 3-5 seconds, stays hidden
- **Pause Action:** Reappears immediately when paused (becomes visible)
- **Paused State:** Always visible, does not fade
- **Resume Action:** Remains visible, then fades after 3-5 seconds
- Text changes synchronously on state transition
- Instant feedback on state changes

**Completion Message:**
- Appears immediately on completion
- Displays for ~500ms
- Then exits automatically
- No user interaction needed

**Source:** [Components](components.md), [State Machines](state-machines.md)

---

## Interaction Patterns

### Immediate Feedback Pattern

**Pattern:**
- All user actions provide immediate visual feedback
- State changes are synchronous
- UI updates within one render cycle (~250ms)
- No perceived delay

**Examples:**
- Press `p`: Status immediately changes to "PAUSED"
- Press `r`: Status immediately changes to "RUNNING"
- Timer completes: Status immediately changes to "COMPLETED"

**Rationale:**
- Users expect immediate response
- Delays feel unresponsive
- Immediate feedback builds confidence

**Source:** [UXDR: Keyboard Controls Design](../uxdr/keyboard-controls-design.md#22-rationale)

---

### State-Based Controls Pattern

**Pattern:**
- Available controls depend on current state
- Control legend shows only valid controls
- Invalid keys are silently ignored
- No error messages for invalid keys

**Examples:**
- Running state: Shows `[p] pause`
- Paused state: Shows `[r] resume` (not `[p] pause`)
- Invalid `p` in paused state: Silently ignored

**Rationale:**
- Reduces cognitive load
- Prevents user errors
- Clear guidance via legend

**Source:** [UXDR: Keyboard Controls Design](../uxdr/keyboard-controls-design.md#42-invalid-key-handling)

---

### Error Handling Pattern

**Pattern:**
- Invalid keys: Silently ignored
- No error messages
- Legend provides guidance
- User can reference legend if unsure

**Rationale:**
- Avoids clutter
- Keeps TUI clean and focused
- Legend provides sufficient guidance

**Source:** [UXDR: Keyboard Controls Design](../uxdr/keyboard-controls-design.md#42-invalid-key-handling)

---

## Keyboard Control Details

### Pause Control (`p`)

**Action:** Pause running timer

**Available In:** Running state only

**Behavior:**
1. User presses `p`
2. Timer state changes to "paused"
3. `pausedAt` set to current time
4. `pausedCount` incremented
5. State saved (event-driven)
6. TUI updates:
   - Status: "⏸ PAUSED" (yellow)
   - Progress bar: Frozen
   - Time display: Frozen
   - Control legend: Reappears immediately, shows `[r] resume [q] quit` (always visible when paused)

**Source:** [US-1.3](../requirements/base.md#us-13), [FR-TIMER-002](../requirements/base.md#fr-timer-002)

---

### Resume Control (`r`)

**Action:** Resume paused timer

**Available In:** Paused state only

**Behavior:**
1. User presses `r`
2. `totalPaused` updated with pause duration
3. `pausedAt` cleared
4. Timer state changes to "running"
5. State saved (event-driven)
6. TUI updates:
   - Status: "RUNNING" (green)
   - Progress bar: Continues
   - Time display: Resumes countdown
   - Control legend: Updates to `[p] pause [q] quit`, remains visible for 3-5 seconds, then fades

**Source:** [US-1.3](../requirements/base.md#us-13), [FR-TIMER-002](../requirements/base.md#fr-timer-002)

---

### Stop Control (`q` or `s`)

**Action:** Stop timer and exit (with confirmation)

**Available In:** Running, Paused states

**Display:** Only `q` shown in legend (minimal UI), `s` works as alias but not displayed

**Behavior:**
1. User presses `q` or `s` (both work identically)
2. **Timer is automatically paused** (if running, timer pauses; if already paused, stays paused)
3. **Confirmation dialog appears:** "Stop timer and exit? [y]es / [n]o" (overlays entire main window)
4. **User confirms or cancels:**
   - **Confirm (`y` or `Y`):**
     - Timer state changes to "stopped"
     - Session saved with `end_status: "stopped"`
     - State file cleaned up
     - TUI exits immediately
     - User returns to command line
   - **Cancel (`n`, `N`, or `Esc`):**
     - Confirmation dialog dismissed
     - Timer resumes (unpauses and continues from where it was)
     - Timer returns to "running" state
     - User continues with timer

**Confirmation Dialog:**
- **Question:** "Stop timer and exit? [y]es / [n]o"
- **Position:** Centered overlay over entire main window (takes over whole timer)
- **Modal:** Blocks other input until confirmed or cancelled
- **Timer Behavior:** Timer automatically paused when dialog appears
- **Keys:** `y`/`Y` to confirm, `n`/`N`/`Esc` to cancel (resumes timer)
- **Other Keys:** Ignored (wait for y/n)

**Rationale for Confirmation:**
- **Prevents Accidents:** Avoids accidental stops when user meant to pause
- **User Control:** Gives user chance to reconsider
- **Clear Intent:** Confirms user wants to stop (destructive action)
- **Common Pattern:** Standard UX pattern for destructive actions

**Rationale for Display Choice:**
- `q` displayed: Common "quit" convention in terminal apps
- `s` not displayed: Alias, showing both is redundant
- Both keys work: Users can use either, but legend shows only `q` for minimal UI

**Source:** [US-1.4](../requirements/base.md#us-14), [FR-TIMER-003](../requirements/base.md#fr-timer-003)

---

### Emergency Exit (`Ctrl+C`) - Not Displayed

**Action:** Emergency exit (interrupt)

**Available In:** All states

**Display:** Not shown in control legend

**Behavior:**
1. User presses `Ctrl+C`
2. Application receives SIGINT
3. State saved (for recovery)
4. `ApplicationInterrupted` event emitted
5. TUI exits immediately
6. User returns to command line

**Note:** 
- State is saved for potential recovery on next start
- This is an emergency action, not part of normal workflow
- Primary stop action (`q` or `s`) should be used for normal exits

**Rationale for Not Displaying:**
- **Universal terminal convention:** Users already know Ctrl+C works in all terminal apps
- **Not Pomodux-specific:** Standard terminal behavior, doesn't need explanation
- **Minimal UI principle:** Don't display universal conventions
- **Reduces clutter:** Legend focuses on Pomodux-specific controls only
- **Still works:** Ctrl+C functions normally, just not displayed

**Source:** [FR-PLUGIN-001](../requirements/base.md#fr-plugin-001), [NFR-REL-001](../requirements/base.md#nfr-rel-001)

---

## References

### Requirements
- [Section 8.2 Keyboard Controls](../requirements/base.md#keyboard-controls)
- [NFR-USE-002](../requirements/base.md#nfr-use-002) - Keyboard Accessibility
- [US-1.3](../requirements/base.md#us-13) - Pause and Resume Timer
- [US-1.4](../requirements/base.md#us-14) - Stop Timer Early
- [FR-TIMER-002](../requirements/base.md#fr-timer-002) - Timer Pause/Resume
- [FR-TIMER-003](../requirements/base.md#fr-timer-003) - Timer Stop

### UXDRs
- [Keyboard Controls Design](../uxdr/keyboard-controls-design.md)

### ADRs
- [TUI Framework Selection](../adr/tui-framework-selection.md) - Bubble Tea Event Loop

### Components
- [Component: Control Legend](components.md#component-6-control-legend)
- [Component: Status Indicator](components.md#component-5-status-indicator)

### Flows
- [Application Flows](application-flows.md)
- [Flow 3: Pause and Resume](application-flows.md#flow-3-pause-and-resume-timer)
- [Flow 4: Stop Timer](application-flows.md#flow-4-stop-timer-early)

### State Machines
- [State Machines](state-machines.md)

---

**Last Updated:** 2026-01-27
