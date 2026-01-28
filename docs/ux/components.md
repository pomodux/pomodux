# Component Specifications

**Version:** 1.0  
**Date:** 2026-01-27

## Overview

This document provides detailed specifications for all UI components in the Pomodux TUI. Each component is specified with its purpose, properties, states, keyboard interactions, theme integration, and implementation notes.

**Source:** [Requirements TUI Specification](../requirements/base.md#tui-specification)

---

## Component Inventory

**Total Components:** 8

1. **Window/Border** - Visual containment
2. **Session Header** - Session type and label display
3. **Progress Bar** - Visual progress representation
4. **Time Display** - Remaining time (MM:SS)
5. **Status Indicator** - Current timer state
6. **Control Legend** - Available keyboard controls
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

### Component 6: Control Legend (Overlay)

**Purpose:** Display available keyboard controls initially, then fade away

**Component Type:** Overlay (not part of main window content)

**Properties:**
- **Format:** Minimal, state-specific controls
- **Position:** Below main window, centered
- **Display Duration:** Visible for 3-5 seconds after timer start, then fades
- **Fade Behavior:** Fades out after initial display period
- **Text Style:** Muted color, regular weight

**Action Classification:**
- **Displayed Actions:** Only essential, state-specific controls shown
- **Hidden but Available:** Aliases (`s` for stop) and universal conventions (Ctrl+C) work but not displayed

**States:**
- **Visible:** Fully opaque, displayed below window
- **Fading:** Gradually reducing opacity (fade out over ~1 second)
- **Hidden:** Not displayed (faded away)

**Display States (Content):**
- **Running:** "[p] pause [q] quit" (fades after 3-5 seconds)
- **Paused:** "[r] resume [q] quit" (always visible when paused)
- **Confirmation:** Not shown (confirmation dialog displayed instead)
- **Completed:** Not shown (exits immediately)

**State Transitions:**
- Timer Start → Visible: Legend appears below window
- Visible → Fading: After 3-5 seconds (running state), begins fade
- Fading → Hidden: After fade completes (~1 second), hidden
- Running → Paused: Legend reappears (becomes visible again) with updated content
- Paused → Running: Legend fades again after 3-5 seconds
- Running/Paused → Confirmation: Legend hidden (confirmation dialog shown instead)
- Confirmation → Running/Paused: Legend reappears (if was visible before)
- Any → Completed: Not shown (exits immediately)

**Fade Timing:**
- **Initial Display (Running):** 3-5 seconds after timer starts
- **Fade Duration:** ~1 second fade out
- **Paused State:** Legend always visible (reappears immediately when paused)
- **Resume Display (Running):** Legend visible for 3-5 seconds after resume, then fades

**Keyboard Interactions:**
- None (display only, shows available controls)
- Fade not interrupted by keypresses (fades regardless in running state)
- Pause action (`p`): Causes legend to reappear immediately
- Resume action (`r`): Legend remains visible, then fades after 3-5 seconds

**Theme Integration:**
- Uses `theme.Colors.TextMuted` for text color
- Fade implementation: Color-based dimming (TUI doesn't support true opacity)
- Single line, minimal display

**Implementation Notes:**
- **Overlay Component:** Rendered separately from main window, positioned below
- **Position:** Centered horizontally, positioned below main window border (Option 1: Below Window, Centered)
- **Layout Impact:** None - window stays fixed, legend fades independently
- Format: `fmt.Sprintf("[%s] %s [%s] %s", key1, action1, key2, action2)`
- Running state: `"[p] pause [q] quit"`
- Paused state: `"[r] resume [q] quit"` (if still visible)

**Fade Implementation (TUI Color-Based):**
- **Approach:** Color brightness reduction (terminals don't support true opacity)
- **Fade Steps:** 4-5 color transitions over ~1 second
- **Update Frequency:** Every ~200-250ms (4-5 steps)
- **Color Progression:**
  1. Start: `theme.Colors.TextMuted` (100% brightness)
  2. Step 1: Dimmed version (~70% brightness)
  3. Step 2: Very dim (~40% brightness)
  4. Step 3: Barely visible (~10% brightness)
  5. End: Stop rendering (invisible)
- **Implementation:** Use progressively dimmer color variants or reduce color intensity
- **Advantage:** No layout changes, smooth visual transition, independent fade behavior
- **Not displayed but still work:**
  - `s` key: Alias for stop, works but not shown (only `q` displayed)
  - `Ctrl+C`: Universal terminal convention, works but not shown
- **Positioning:** Centered below main window, not inside window border
- **Rationale:** 
  - Fades away during running state to reduce visual clutter
  - Reappears when paused to provide contextual help (controls change)
  - Paused state is interactive, so showing controls is helpful
  - Keeps UI minimal during active countdown, shows controls when user needs them

**Requirements References:**
- [Section 8.2 Keyboard Controls](../requirements/base.md#keyboard-controls)
- [UXDR: Keyboard Controls Design](../uxdr/keyboard-controls-design.md)

**Wireframe Usage:** Running, Paused states (not Completed)

---

### Component 7: Completion Message

**Purpose:** Show completion confirmation

**Properties:**
- **Text:** "Session saved!"
- **Position:** Below status indicator
- **Display Duration:** ~500ms before exit
- **Text Style:** Success color, regular weight

**States:**
- **Completed:** Visible briefly (~500ms), then exit

**State Transitions:**
- Running → Completed: Message appears
- Completed → Exit: Message disappears after ~500ms

**Keyboard Interactions:**
- None (display only, brief)

**Theme Integration:**
- Uses `theme.Colors.Success` for text color
- Success color for positive confirmation

**Implementation Notes:**
- Display after timer completes
- Show for ~500ms (2 tick cycles at 250ms)
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

### Component 9: Confirmation Dialog

**Purpose:** Confirm user intent before stopping timer

**Component Type:** Overlay (modal dialog)

**Properties:**
- **Format:** Question with yes/no options
- **Position:** Centered over main window (overlay)
- **Display Condition:** User presses `q` or `s` to stop timer
- **Text Style:** Warning/question color, regular weight

**Content:**
- **Question:** "Stop timer and exit? [y]es / [n]o"
- **Format:** Brief, clear question with keyboard shortcuts

**States:**
- **Visible:** Displayed when user requests stop
- **Hidden:** Not displayed (normal state)

**State Transitions:**
- Running/Paused → Confirmation: User presses `q` or `s` (timer paused automatically)
- Confirmation → Running: User presses `n` (cancel, timer resumes)
- Confirmation → Stopped: User presses `y` (confirm)

**Keyboard Interactions:**
- `y` or `Y`: Confirm stop, exit application
- `n` or `N`: Cancel, resume timer (timer unpauses and continues)
- `Esc`: Cancel (same as `n`, resume timer)
- Other keys: Ignored (wait for y/n)

**Theme Integration:**
- Uses `theme.Colors.Warning` or `theme.Colors.Text` for question text
- Keys highlighted: `[y]es` and `[n]o` in brackets
- Centered overlay, semi-transparent or dimmed background (if supported)

**Implementation Notes:**
- **Overlay Component:** Rendered on top of main window (takes over whole timer)
- Format: `"Stop timer and exit? [y]es / [n]o"`
- Position: Centered horizontally and vertically, overlays entire main window
- **Modal Behavior:** Blocks other input until confirmed or cancelled
- **Timer Behavior:** Timer is automatically paused when confirmation appears
- **Cancel Behavior:** Timer resumes (unpauses) when user cancels with `n` or `Esc`
- **Rationale:** Prevents accidental stops, gives user chance to reconsider, timer paused during decision

**Requirements References:**
- [US-1.4](../requirements/base.md#us-14) - Stop Timer Early
- [FR-TIMER-003](../requirements/base.md#fr-timer-003) - Timer Stop

**Wireframe Usage:** Confirmation state (overlay on Running/Paused)

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
└── Completion Message (conditional)
```

**Overlay Components (Transient):**
```
Control Legend (overlay, fades after 3-5s)
Confirmation Dialog (overlay, appears on stop request)
```

**Architecture:**
- Main window components are persistent (always visible)
- Overlay components are transient (appear, then fade)
- Control legend is separate overlay, positioned below main window
- Overlay architecture allows independent fade behavior

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

**Control Legend:**
- Depends on: Timer state (for available controls)
- Affects: User guidance

**Completion Message:**
- Depends on: Timer state (completed)
- Affects: Confirmation feedback

### Component Composition Patterns

**Running State:**
- **Main Window:** Window + Session Header + Progress Bar + Time Display + Status Indicator
- **Overlay:** Control Legend (visible initially, fades after 3-5s)

**Paused State:**
- **Main Window:** Window + Session Header + Progress Bar (frozen) + Time Display (frozen) + Status Indicator
- **Overlay:** Control Legend (always visible when paused, shows `[r] resume [q] quit`)

**Completed State:**
- **Main Window:** Window + Session Header + Progress Bar (100%) + Time Display (0:00) + Status Indicator + Completion Message
- **Overlay:** Control Legend (not shown, already faded)

**Error State (Terminal Too Small):**
- Terminal Size Warning (centered)

**Error State (Config Errors):**
- **Main Window:** Window + Config Error Banner + Session Header + Progress Bar + Time Display + Status Indicator
- **Overlay:** Control Legend (if still visible; typically already faded)

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

**Last Updated:** 2026-01-27
