---
status: approved
type: technical
---

# ADR 007: TUI Standardization for Pomodux

## 1. Context / Background

Pomodux has historically used two different TUI libraries:
- **Bubbletea** for the main timer interface (from release 0.5.0)
- **rivo/tview** for plugin-driven dialogs (modals, lists, prompts)

This mixed approach led to inconsistencies in user experience, increased maintenance complexity, and confusion for contributors. As the project matures, a unified TUI stack is needed for maintainability and a consistent look and feel.

## 2. Decision

**Pomodux will standardize on Bubbletea for all user-facing TUI, including the main timer interface and all plugin dialogs.**

- All new TUI features must use Bubbletea and its ecosystem (e.g., bubbles components).
- Existing plugin dialogs implemented with tview will be migrated to Bubbletea in release 0.5.1.
- **Theming (colors, styles, UI elements) will be standardized across all Bubbletea-based UI, including the main timer and plugin dialogs, to ensure a consistent look and feel.**
- **Lipgloss will be used as the standard library for theming and styling all Bubbletea-based UI. A shared theme (defined with Lipgloss) will be used across the application.**

## 3. Rationale

- **Consistency:** Ensures a unified user experience across all UI components, including consistent theming and visual style.
- **Maintainability:** Reduces dependencies and cognitive load for contributors.
- **Modern Architecture:** Bubbletea and Lipgloss together offer a declarative, reactive, and easily styled TUI model well-suited for Pomodux's needs.
- **Community Support:** Bubbletea and Lipgloss have active ecosystems and are widely adopted for modern terminal UIs.

## 4. Alternatives Considered

- **Continue using both Bubbletea and tview:**
  - Rejected due to inconsistent UX and increased maintenance burden.
- **Standardize on tview:**
  - Rejected because the main timer UI is already implemented in Bubbletea, and Bubbletea is better suited for responsive, modern TUI design.
- **Custom abstraction layer:**
  - Considered unnecessary complexity for current project needs.

## 5. Consequences

- **Short-term:**
  - Requires migration of plugin dialogs from tview to Bubbletea (planned for 0.5.1).
  - Contributors must learn Bubbletea and Lipgloss patterns for all TUI work.
- **Long-term:**
  - Unified, modern, and maintainable TUI codebase.
  - Easier onboarding and future feature development.

## 6. Migration Plan

- Refactor all plugin dialog utilities (modals, lists, prompts) to use Bubbletea.
- Update the plugin API to call Bubbletea-based dialogs.
- **Apply a unified theming system using Lipgloss (colors, styles, UI elements) to all Bubbletea-based UI, including the main timer and plugin dialogs.**
- Remove tview as a dependency after migration is complete.
- Test all plugin dialog flows for parity and UX consistency.

## 7. Status

- **Approved** (2025-07-21)
- Migration scheduled for release 0.5.1

## 8. References

- [Release 0.5.0 Notes](../releases/release-0.5.0.md)
- [Bubbletea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [tview Documentation](https://github.com/rivo/tview) 