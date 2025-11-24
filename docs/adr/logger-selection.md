---
status: approved
---

# Structured Logger Architecture for Pomodux

## 1. Context / Background

Robust logging is essential for maintainability, debugging, and user support.

## 2. Requirements
- **Log Levels:** DEBUG, INFO, WARN, ERROR
- **Structured Fields:** Attach context (component, plugin, event, etc.)
- **Configurable Output:** Console, file, or both
- **Format:** Human-readable text and machine-parsable JSON
- **Configurable via config file** (XDG-compliant)
- **Minimal Overhead:** No impact on timer accuracy or performance
- **Cross-Platform:** No OS-specific dependencies
- **Easy Integration:** Usable in all parts of the codebase

## 3. Decision

**Selected Solution:**
- Use a mature Go logging library (**logrus**)
- Add a new `internal/logger` package
- Expose a global logger instance with helper functions for each log level
- Add a `logging` section to the config file for user control
- Replace all `fmt.Printf`/`fmt.Fprintf` debug/info/warn/error calls with logger calls
- Provide a simple API for plugins (future work)

## 4. Alternatives Considered
- **Use fmt.Printf:** Not scalable, no log levels, no structure
- **Use Go’s standard log package:** Lacks structured fields and log levels
- **Write a custom logger:** Reinvents the wheel, more maintenance

## 5. Risks
- **Migration:** Need to update all existing debug/info/error output
- **User confusion:** Must clearly separate user-facing output from logs
- **Performance:** Must ensure logging does not block timer or CLI responsiveness

