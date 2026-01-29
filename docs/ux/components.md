# Component Specifications

**Version:** 1.0  
**Date:** 2026-01-27

## Overview

This document provides detailed specifications for all UI components in the Pomodux TUI. Each component is specified with its purpose, properties, states, keyboard interactions, theme integration, and implementation notes.

**Source:** [Requirements TUI Specification](../requirements/base.md#tui-specification)

---

## Component Inventory

**Total Components:** 7

1. **Window/Border** - Visual containment
2. **Session Header** - Session type and label display
3. **Progress Bar** - Visual progress representation
4. **Time Display** - Remaining time (MM:SS)
5. **Status Indicator** - Current timer state
6. **Action Selection** - Available actions based on timer state (inline)
7. **Completion Message** - Completion confirmation (brief)
8. **Terminal Size Warning** - Terminal too small warning

---

## Component Specifications

### Component 1: Window/Border

**Purpose:** Visual containment of timer display, provides window frame

**Properties:**
- **Title:** "Pomodux Timer" (static text)
- **Border Style:** Theme-dependent (rounded/square/double/none)
- **Border Color:** Theme border color
- **Width:** Responsive to terminal size (centered)
- **Height:** Responsive to content (centered vertically)
- **Padding:** 1-2 cells internal padding

**States:**
- Always visible (no state changes)
- Border style changes with theme

**Keyboard Interactions:**
- None (container only)

**Theme Integration:**
- Uses `theme.BorderStyle()` for border rendering
- Border color from `theme.Colors.Border`
- Border style from `theme.Border.Style`

**Implementation Notes:**
- Use Lipgloss `Border()` style
- Center horizontally and vertically using Lipgloss layout
- Border style: `lipgloss.RoundedBorder()` (default), configurable per theme

**Requirements References:**
- [Section 8.1 Screen Layouts](../requirements/base.md#screen-layouts)
- [Section 5.4 Theme Definition Schema](../requirements/base.md#54-theme-definition-schema)

**Wireframe Usage:** All screen states

---

### Component 2: Session Header

**Purpose:** Display session type (preset) and label

**Properties:**
- **Format:** "{Preset} Session: {Label}" or "Session: {Label}" (if no preset)
- **Preset Prettification:** 
  - "work" → "Work"
  - "longbreak" → "Long Break"
  - "custom_preset" → "Custom Preset"
- **Text Style:** Primary color, regular weight
- **Position:** Below window title, above progress bar
- **Alignment:** Left-aligned within window

**States:**
- Static (doesn't change during timer)
- Content changes based on preset/label

**Keyboard Interactions:**
- None (display only)

**Theme Integration:**
- Uses `theme.TitleStyle()` or `theme.Colors.Primary` for text color
- Text color from `theme.Colors.Primary`

**Implementation Notes:**
- Prettify preset name: Split on underscores/hyphens, capitalize each word
- Format: `fmt.Sprintf("%s Session: %s", prettifiedPreset, label)`
- If no preset: `fmt.Sprintf("Session: %s", label)`

**Requirements References:**
- [Section 8.1.1 Running Timer State](../requirements/base.md#running-timer-state)
- [US-1.5](../requirements/base.md#us-15) - Default Label for Quick Start
- [UXDR: Session Label Defaults](../uxdr/session-label-defaults.md)

**Wireframe Usage:** Running, Paused, Completed states

---

### Component 3: Progress Bar

**Purpose:** Visual representation of elapsed time

**Properties:**
- **Component:** Bubbles `progress.Model` (from `github.com/charmbracelet/bubbles/progress`)
- **Characters:** 
  - Filled: Theme `progress.filled_char` (default: "█")
  - Empty: Theme `progress.empty_char` (default: "░")
- **Width:** Responsive (max 80 columns, adapts to terminal width)
- **Percentage:** Displayed if `theme.progress.show_percentage: true` (default: true)
- **Colors:** 
  - Filled: Theme `colors.progress_filled`
  - Empty: Theme `colors.progress_empty`
- **Gradient:** Theme-based gradient (via Bubbles component)
- **Position:** Below session header, above time display

**States:**
- **Running:** Updates every 250ms, shows progress (0% → 100%)
- **Paused:** Frozen at current progress (no updates)
- **Completed:** 100% filled, static

**State Transitions:**
- Running → Paused: Freeze at current progress
- Paused → Running: Resume updates from frozen progress
- Running → Completed: Fill to 100%, stop updates

**Keyboard Interactions:**
- None (display only, affected by timer state)

**Theme Integration:**
- Uses `theme.ProgressFilledStyle()` for filled portion
- Uses `theme.ProgressEmptyStyle()` for empty portion
- Colors from `theme.Colors.ProgressFilled` and `theme.Colors.ProgressEmpty`
- Characters from `theme.Progress.FilledChar` and `theme.Progress.EmptyChar`
- Percentage display from `theme.Progress.ShowPercentage`

**Implementation Notes:**
- Use Bubbles `progress.Model` component
- Initialize: `progress.New(progress.WithDefaultGradient())`
- Set width: `progress.Width = terminalWidth - padding`
- Update: `progress.SetPercent(percent)` where `percent = elapsed / duration`
- Render: `progress.ViewAs(percent)`
- Apply theme colors via Lipgloss styles

**Requirements References:**
- [Section 8.1.1 Running Timer State](../requirements/base.md#running-timer-state)
- [ADR: TUI Framework Selection](../adr/tui-framework-selection.md) - Bubbles Components

**Wireframe Usage:** Running, Paused, Completed states

---

### Component 4: Time Display

**Purpose:** Show remaining time in MM:SS format

**Properties:**
- **Format:** "MM:SS" (e.g., "15:23", "0:45", "0:05")
- **Update Frequency:** Every 250ms (via tick message)
- **Text Style:** Primary color, regular weight, monospace font preferred
- **Position:** Below progress bar, above status indicator
- **Alignment:** Left-aligned within window

**States:**
- **Running:** Counts down (15:23 → 15:22 → ... → 0:00)
- **Paused:** Frozen at current time (e.g., "15:23")
- **Completed:** Shows "0:00"

**State Transitions:**
- Running → Paused: Freeze at current time
- Paused → Running: Resume countdown from frozen time
- Running → Completed: Show "0:00"

**Keyboard Interactions:**
- None (display only, affected by timer state)

**Theme Integration:**
- Uses `theme.Colors.Primary` for text color
- Text color from theme primary color

**Implementation Notes:**
- Format: `fmt.Sprintf("%02d:%02d", minutes, seconds)`
- Calculate from `timer.Remaining()` duration
- Update on every tick message (250ms)
- Use monospace font for consistent width

**Requirements References:**
- [Section 8.1.1 Running Timer State](../requirements/base.md#running-timer-state)
- [FR-TIMER-001](../requirements/base.md#fr-timer-001) - Timer Start

**Wireframe Usage:** Running, Paused, Completed states

---

### Component 5: Status Indicator

**Purpose:** Display current timer state

**Properties:**
- **Text:** 
  - Running: "RUNNING"
  - Paused: "⏸ PAUSED"
  - Completed: "✓ COMPLETED"
- **Color:** Theme-dependent semantic colors
  - Running: Success (green)
  - Paused: Warning (yellow)
  - Completed: Success (green)
- **Position:** Below time display, above control legend
- **Alignment:** Left-aligned within window

**States:**
- **Running:** "RUNNING" (green/success color)
- **Paused:** "⏸ PAUSED" (yellow/warning color)
- **Completed:** "✓ COMPLETED" (green/success color)

**State Transitions:**
- Running → Paused: Text changes to "⏸ PAUSED", color changes to warning
- Paused → Running: Text changes to "RUNNING", color changes to success
- Running → Completed: Text changes to "✓ COMPLETED", color stays success

**Keyboard Interactions:**
- None (display only, reflects timer state)

**Theme Integration:**
- Uses `theme.StatusStyle(status)` for color
- Colors from `theme.Colors.Success` (running/completed) and `theme.Colors.Warning` (paused)
- Status-based color mapping

**Implementation Notes:**
- Use `theme.StatusStyle(status)` method
- Status values: "running", "paused", "completed"
- Unicode symbols: ⏸ (pause), ✓ (checkmark)

**Requirements References:**
- [Section 8.1.1 Running Timer State](../requirements/base.md#running-timer-state)
- [Section 8.1.2 Paused Timer State](../requirements/base.md#812-paused-timer-state)
- [Section 8.1.3 Completion State](../requirements/base.md#813-completion-state)
- [Section 8.4 Theme Application](../requirements/base.md#theme-application)

**Wireframe Usage:** Running, Paused, Completed states

---

### Component 6: Action Selection

**Purpose:** Display available actions based on timer state, including inline confirmation

**Component Type:** Inline (part of main window content)

**Properties:**
- **Format:** State-specific action prompts
- **Position:** Below status indicator, inside main window
- **Text Style:** Muted color, regular weight
- **Always Visible:** No fading, always displayed

**Display States (Content):**
- **Running:** "[p]ause  [s]top"
- **Paused:** "[r]esume  [s]top"
- **Confirmation:** "Stop timer and exit? [y]es / [n]o" (inline, timer paused)
- **Completed:** Not shown (exits immediately)

**State Transitions:**
- Timer Start → Running: Shows "[p]ause  [s]top"
- Running → Paused: Content changes to "[r]esume  [s]top"
- Paused → Running: Content changes to "[p]ause  [s]top"
- Running/Paused → Confirmation: Content changes to "Stop timer and exit? [y]es / [n]o" (timer paused automatically)
- Confirmation → Running/Paused: Content changes back to state-appropriate actions (timer resumes if was running)
- Any → Completed: Not shown (exits immediately)

**Keyboard Interactions:**
- **Running State:**
  - `p`: Pause timer (transitions to paused state)
  - `q` or `s`: Enter confirmation state (timer paused automatically)
- **Paused State:**
  - `r`: Resume timer (transitions to running state)
  - `q` or `s`: Enter confirmation state
- **Confirmation State:**
  - `y` or `Y`: Confirm stop, exit application
  - `n`, `N`, or `Esc`: Cancel, return to previous state (timer resumes if was running)
  - Other keys: Ignored (wait for y/n)

**Theme Integration:**
- Uses `theme.Colors.TextMuted` for text color
- Confirmation text may use `theme.Colors.Warning` for question
- Keys highlighted in brackets: `[p]`, `[q]`, `[y]es`, `[n]o`

**Implementation Notes:**
- **Inline Component:** Rendered inside main window, below status indicator
- **Position:** Inside window border, below status indicator
- **No Fade:** Always visible, no fade behavior
- **State-Based Content:** Content changes immediately on state transitions
- **Confirmation Behavior:** 
  - When user presses `q` or `s`, timer is automatically paused
  - Action selection content changes to confirmation prompt
  - User confirms or cancels within the same component
  - On cancel, timer resumes if it was running before confirmation
- Format examples:
  - Running: `"[p]ause  [s]top"`
  - Paused: `"[r]esume  [s]top"`
  - Confirmation: `"Stop timer and exit? [y]es / [n]o"`
- **Rationale:** 
  - Simpler than overlay approach
  - Always visible, no hidden controls
  - Inline confirmation reduces visual complexity
  - No fade timing complexity

**Requirements References:**
- [Section 8.2 Keyboard Controls](../requirements/base.md#keyboard-controls)
- [US-1.3](../requirements/base.md#us-13) - Pause and Resume Timer
- [US-1.4](../requirements/base.md#us-14) - Stop Timer Early
- [FR-TIMER-002](../requirements/base.md#fr-timer-002) - Timer Pause/Resume
- [FR-TIMER-003](../requirements/base.md#fr-timer-003) - Timer Stop

**Wireframe Usage:** Running, Paused, Confirmation states (not Completed)

---

### Component 7: Completion Message

**Purpose:** Show completion confirmation with countdown

**Properties:**
- **Text:** "Session saved! Closing in 3. 2. 1."
- **Position:** Below status indicator
- **Display Duration:** ~3 seconds (countdown from 3 to 1)
- **Text Style:** Success color, regular weight
- **Countdown:** Updates every second (3 → 2 → 1)

**States:**
- **Completed:** Visible with countdown (~3 seconds), then exit
- **Countdown States:**
  - Initial: "Session saved! Closing in 3."
  - After 1 second: "Session saved! Closing in 2."
  - After 2 seconds: "Session saved! Closing in 1."
  - After 3 seconds: Exit

**State Transitions:**
- Running → Completed: Message appears with "Closing in 3."
- After 1 second: Updates to "Closing in 2."
- After 2 seconds: Updates to "Closing in 1."
- After 3 seconds: Exit

**Keyboard Interactions:**
- None (display only, countdown proceeds automatically)

**Theme Integration:**
- Uses `theme.Colors.Success` for text color
- Success color for positive confirmation

**Implementation Notes:**
- Display after timer completes
- Show countdown: "Session saved! Closing in 3." → "Closing in 2." → "Closing in 1."
- Update every 1 second (4 tick cycles at 250ms per second)
- Total display duration: ~3 seconds
- Then exit TUI

**Requirements References:**
- [Section 8.1.3 Completion State](../requirements/base.md#813-completion-state)
- [UXDR: Auto-Exit Behavior](../uxdr/auto-exit-behavior.md) - Immediate exit

**Wireframe Usage:** Completed state only

---

### Component 8: Terminal Size Warning

**Purpose:** Warn user if terminal too small

**Properties:**
- **Text:** "Terminal too small!\nMinimum: 80x24\nCurrent: {width}x{height}"
- **Border:** Simple border
- **Position:** Centered in terminal
- **Display Condition:** Terminal size < 80x24

**States:**
- **Too Small:** Visible
- **Adequate Size:** Not shown

**State Transitions:**
- Adequate → Too Small: Warning appears
- Too Small → Adequate: Warning disappears

**Keyboard Interactions:**
- None (informational only)

**Theme Integration:**
- Uses `theme.Colors.Warning` for text color
- Warning color for alert

**Implementation Notes:**
- Check terminal size via `tea.WindowSizeMsg`
- Display if `width < 80 || height < 24`
- Show warning, continue with degraded layout OR exit

**Requirements References:**
- [Section 8.3 Terminal Resize Handling](../requirements/base.md#terminal-resize-handling)

**Wireframe Usage:** Error state (terminal too small)

---


---

## Component Relationships

### Component Hierarchy

**Main Window (Persistent Container):**
```
Window/Border (container)
├── Session Header
├── Progress Bar
├── Time Display
├── Status Indicator
├── Action Selection (always visible, state-based content)
└── Completion Message (conditional)
```

**Architecture:**
- All components are inline (part of main window)
- No overlay components
- Action Selection changes content based on timer state
- Confirmation handled inline within Action Selection component

### Component Dependencies

**Progress Bar:**
- Depends on: Timer state (for progress calculation)
- Affects: Visual feedback

**Time Display:**
- Depends on: Timer state (for remaining time)
- Affects: User awareness

**Status Indicator:**
- Depends on: Timer state (for status text/color)
- Affects: State communication

**Action Selection:**
- Depends on: Timer state (for available actions)
- Affects: User guidance and confirmation

**Completion Message:**
- Depends on: Timer state (completed)
- Affects: Confirmation feedback

### Component Composition Patterns

**Running State:**
- **Main Window:** Window + Session Header + Progress Bar + Time Display + Status Indicator + Action Selection (`[p]ause  [s]top`)

**Paused State:**
- **Main Window:** Window + Session Header + Progress Bar (frozen) + Time Display (frozen) + Status Indicator + Action Selection (`[r]esume  [s]top`)

**Confirmation State:**
- **Main Window:** Window + Session Header + Progress Bar (frozen) + Time Display (frozen) + Status Indicator (paused) + Action Selection (`Stop timer and exit? [y]es / [n]o`)

**Completed State:**
- **Main Window:** Window + Session Header + Progress Bar (100%) + Time Display (0:00) + Status Indicator + Completion Message
- **Action Selection:** Not shown (exits immediately)

**Error State (Terminal Too Small):**
- Terminal Size Warning (centered)

**Error State (Config Errors):**
- **Main Window:** Window + Config Error Banner + Session Header + Progress Bar + Time Display + Status Indicator + Action Selection

---

## References

### Requirements
- [Section 8 TUI Specification](../requirements/base.md#tui-specification)
- [Section 8.1 Screen Layouts](../requirements/base.md#screen-layouts)
- [Section 8.2 Keyboard Controls](../requirements/base.md#keyboard-controls)
- [Section 8.4 Theme Application](../requirements/base.md#theme-application)

### UXDRs
- [Keyboard Controls Design](../uxdr/keyboard-controls-design.md)
- [Session Label Defaults](../uxdr/session-label-defaults.md)
- [Auto-Exit Behavior](../uxdr/auto-exit-behavior.md)

### ADRs
- [TUI Framework Selection](../adr/tui-framework-selection.md) - Bubbles Components

### Screens
- [Screen: Running State](screens.md#running-state)
- [Screen: Paused State](screens.md#paused-state)
- [Screen: Completed State](screens.md#completed-state)

### Analysis

---

**Last Updated:** 2026-01-28
