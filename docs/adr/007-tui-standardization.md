---
status: approved
type: technical
---

# ADR 007: TUI Standardization for Pomodux

## 1. Context / Background

Pomodux uses **Bubbletea** as the exclusive TUI framework for all user interfaces. This decision ensures consistency in user experience, reduces maintenance complexity, and provides a unified development approach for contributors.

## 2. Decision

**Pomodux uses Bubbletea exclusively for all user-facing TUI components.**

- All TUI features use Bubbletea and its ecosystem (e.g., bubbles components)
- **Lipgloss** is the standard library for theming and styling all Bubbletea-based UI
- **teatest** is the exclusive testing framework for all TUI components
- A shared theme (defined with Lipgloss) ensures consistent look and feel across the application

## 3. Rationale

- **Consistency:** Ensures a unified user experience across all UI components, including consistent theming and visual style.
- **Maintainability:** Reduces dependencies and cognitive load for contributors.
- **Modern Architecture:** Bubbletea and Lipgloss together offer a declarative, reactive, and easily styled TUI model well-suited for Pomodux's needs.
- **Community Support:** Bubbletea and Lipgloss have active ecosystems and are widely adopted for modern terminal UIs.

## 4. Alternatives Considered

- **Custom TUI framework:**
  - Rejected due to unnecessary complexity and maintenance overhead
- **Multiple TUI libraries:**
  - Rejected due to inconsistent UX and increased maintenance burden
- **Custom abstraction layer:**
  - Rejected as unnecessary complexity for current project needs

## 5. Consequences

- **Short-term:**
  - Contributors must learn Bubbletea, Lipgloss, and teatest patterns for all TUI work
  - Consistent development patterns across all UI components
- **Long-term:**
  - Unified, modern, and maintainable TUI codebase
  - Easier onboarding and future feature development
  - Simplified dependency management

## 6. Implementation Standards

- All TUI components must use Bubbletea framework
- All styling must use Lipgloss for consistency
- All TUI testing must use teatest framework
- Apply unified theming system using Lipgloss across all UI components
- Follow Bubbletea Model-View-Update (MVU) architecture pattern

## 7. Status

- **Approved** (2025-07-21)
- **Implemented** - Bubbletea ecosystem is the exclusive TUI standard

## 8. References

- [Bubbletea Documentation](https://github.com/charmbracelet/bubbletea) - Primary TUI framework
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss) - Styling and theming
- [teatest Documentation](https://github.com/charmbracelet/x/exp/teatest) - TUI testing framework 

## TUI Testing Strategy

- **teatest** ([github.com/charmbracelet/x/exp/teatest](https://github.com/charmbracelet/x/exp/teatest)) is the exclusive framework for testing Bubbletea-based TUIs. It provides comprehensive testing capabilities including simulation of keypresses, window resizes, and assertions on output and model state.

### Rationale
- teatest is part of the Bubbletea ecosystem and provides native integration with Bubbletea models
- All Bubbletea-based TUI code must be tested using teatest to ensure consistency and reliability
- No alternative TUI testing frameworks are needed or should be used 