---
status: approved
---

# Plugin System Architecture for Pomodux

## 1. Context / Background

### 1.1 Motivation
Pomodux aims to be a highly extensible terminal timer application. Users and developers have requested the ability to extend Pomodux with custom features, notifications, integrations, and analytics. A plugin system is needed to:
- Enable user and community-driven feature development
- Allow for rapid prototyping and experimentation
- Support integrations with external tools and services
- Provide advanced customization (notifications, statistics, themes, etc.)

### 1.2 Requirements
- **Extensibility:** Allow third-party plugins to add/modify behavior
- **Security:** Sandbox plugin execution to prevent malicious actions
- **Performance:** Minimal impact on timer accuracy and resource usage
- **Cross-Platform:** Work on Linux, macOS, and Windows
- **Event-Driven:** Plugins should react to timer/session events
- **Ease of Use:** Simple API for plugin authors

## 2. Decision

**Selected Solution:** Lua-based plugin system using [gopher-lua](https://github.com/yuin/gopher-lua)

### 2.1 Rationale
- **Lua is lightweight, embeddable, and widely used for plugin systems** (e.g., Neovim)
- **gopher-lua** is a pure Go implementation, ensuring easy integration and cross-platform support
- **Event-driven architecture** allows plugins to subscribe to timer lifecycle events (start, pause, resume, complete, stop)
- **Sandboxing** is feasible with Lua, limiting plugin access to only the provided API
- **Proven approach:** Many successful projects use Lua for extensibility

### 2.2 Key Design Points
- **Plugin Lifecycle:** Plugins are loaded at startup or on demand, and can be enabled/disabled/unloaded at runtime
- **Event Hooks:** Plugins register Lua functions for specific timer events (e.g., `timer_started`, `timer_completed`)
- **API Surface:** Plugins can access timer/session data, send notifications, and log output
- **Isolation:** Each plugin runs in its own Lua state
- **Performance:** Event dispatch is asynchronous to avoid blocking the main timer loop

## 3. Alternatives Considered

### 3.1 No Plugin System
- **Pros:** Simpler codebase, less maintenance
- **Cons:** No extensibility, all features must be built-in, less community involvement

### 3.2 Go Plugins (native Go plugin package)
- **Pros:** Native performance, type safety
- **Cons:** Not cross-platform, complex build/distribution, unsafe for untrusted code

### 3.3 Scripting with Python/JS
- **Pros:** Popular languages, rich ecosystems
- **Cons:** Heavyweight, more dependencies, harder to sandbox, larger attack surface

### 3.4 External Process Plugins (IPC)
- **Pros:** Strong isolation, language agnostic
- **Cons:** Complex IPC, higher resource usage, more difficult for simple use cases

## 4. Consequences

### 4.1 Positive
- **Extensible:** Users can add new features without modifying core code
- **Safe:** Sandboxed Lua environment limits plugin capabilities
- **Cross-Platform:** Works on all supported OSes
- **Community-Friendly:** Low barrier for plugin authors
- **Maintainable:** Clear separation between core and extensions

### 4.2 Negative
- **Learning Curve:** Plugin authors must learn Lua
- **Debugging:** Errors in plugins may be harder to trace
- **Resource Usage:** Each plugin has its own Lua state (minimal, but nonzero overhead)
- **API Stability:** Need to maintain a stable plugin API

### 4.3 Risks
- **Security:** Plugins could attempt to escape sandbox (mitigated by strict API)
- **Performance:** Poorly written plugins could impact performance (mitigated by async event dispatch)
- **Complexity:** Plugin manager adds architectural complexity

## 5. Implementation Status

> **Status: MVP — Event type definitions only.**
>
> The full Lua-based plugin system described in this ADR has **not yet been implemented**. The `gopher-lua` dependency has not been added to `go.mod`. What exists today is the foundational event architecture that the plugin system will build on.

### 5.1 What Is Implemented (MVP)
- **Event type definitions** in `internal/plugin/events.go` — nine event types covering timer and application lifecycle (`TimerStarted`, `TimerPaused`, `TimerResumed`, `TimerStopped`, `TimerCompleted`, `ApplicationStarted`, `ApplicationStopping`, `ApplicationInterrupted`, `ConfigurationLoaded`)
- **Event emitter** struct with an `Emit()` method that logs events via the logger — serves as a placeholder for future plugin dispatch
- Events are emitted from `cmd/pomodux/main.go` at the appropriate lifecycle points

### 5.2 What Is NOT Yet Implemented
- **PluginManager** — no plugin loading, lifecycle management, or directory scanning
- **Lua integration** — `gopher-lua` is not a dependency; no Lua states are created
- **Plugin API surface** — no registration, hook management, or notification capabilities exposed to plugins
- **Sandboxing** — no isolation or security boundary exists
- **Example plugins** — no plugin files exist in the repository
- **Plugin configuration** — no plugin directory setting in the config schema
- **Async event dispatch** — `Emit()` logs synchronously

### 5.3 Next Steps
1. Add `github.com/yuin/gopher-lua` to `go.mod`
2. Implement `PluginManager` with directory scanning and plugin lifecycle
3. Create Lua sandbox with restricted API surface
4. Wire `Emit()` to dispatch events to loaded plugins
5. Write example plugins (notification, debug logging)
6. Add plugin directory to config schema
7. Document plugin authoring guide

## 6. References
- [gopher-lua](https://github.com/yuin/gopher-lua)
- [Neovim Plugin System](https://neovim.io/)

