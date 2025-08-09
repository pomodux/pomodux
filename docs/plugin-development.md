# Pomodux Plugin Development Guide

This guide explains how to create plugins for Pomodux using the current plugin API.

> **📋 Note**: Plugin system is currently stable. Future API enhancements may be planned but no breaking changes are currently scheduled.

---

## Getting Started

Pomodux supports plugins written in Lua. Plugins can react to timer events, display notifications, prompt for user input, and more.

### Plugin Location
- Place your `.lua` plugin files in your Pomodux plugins directory (e.g., `config/pomodux/plugins/`).
- Enable plugins in your Pomodux `config.yaml`.

---

## Registering a Plugin

Every plugin must register itself:

```lua
pomodux.register_plugin({
    name = "my_plugin",
    version = "1.0.0",
    description = "My plugin description",
    author = "Your Name"
})
```

---

## Registering Hooks

Plugins can react to timer events by registering hooks:

```lua
pomodux.register_hook("timer_setup", function(event)
    -- Your code here
end)
```

**Common event types:**
- `"timer_setup"`
- `"timer_started"`
- `"timer_paused"`
- `"timer_resumed"`
- `"timer_completed"`
- `"timer_stopped"`

The `event` argument is a table with:
- `event.type` (string)
- `event.timestamp` (number, unix time)
- `event.data` (table, event-specific data)

---

## Plugin API Functions

Pomodux exposes several TUI functions to plugins via the `pomodux` table:

### 1. `pomodux.show_notification(message: string) -> boolean`
Displays a modal notification dialog with the given message and OK/Cancel buttons.

### 2. `pomodux.select_from_list(title: string, options: table) -> (number|nil, boolean)`
Shows a modal list selection dialog. Returns the 1-based index of the selected option and `true` if confirmed, or `nil, false` if cancelled.

### 3. `pomodux.input_prompt(title: string, default_value: string, placeholder: string) -> (string, boolean)`
Shows a modal input prompt dialog. Returns the entered string and `true` if confirmed, or `"", false` if cancelled.

### 4. `pomodux.log(message: string)`
Logs a debug message to the Pomodux log (for plugin debugging).

### 5. `pomodux.get_config(key: string, default_value: string) -> string`
(Stub) Intended to retrieve configuration values. Currently always returns the default value.

---

## Example: Plugin Implementation

This example demonstrates how to use the plugin API:

```lua
-- Example Plugin for Pomodux
-- Demonstrates usage of the Plugin API

pomodux.register_plugin({
    name = "example_plugin",
    version = "1.0.0",
    description = "Shows how to use Pomodux Plugin API.",
    author = "Pomodux Team"
})

-- Register a hook for timer setup (runs before timer starts)
pomodux.register_hook("timer_setup", function(event)
    -- Show a notification dialog
    local confirmed = pomodux.show_notification("Welcome to your timer session!")
    
    -- Log debug information
    pomodux.log("Timer setup completed")
end)

-- Register a hook for timer completion
pomodux.register_hook("timer_completed", function(event)
    -- Show completion notification
    pomodux.show_notification("Great job! Session completed.")
end)

-- Register a hook for timer start
pomodux.register_hook("timer_started", function(event)
    -- Example of selecting from list
    local option, confirmed = pomodux.select_from_list("Break Type", {"Short Break", "Long Break", "Continue Working"})
    if confirmed then
        pomodux.log("User selected option: " .. option)
    end
end)

print("✅ Example Plugin loaded!")
```

---

## Example: TUI Plugin

This example demonstrates how to use the TUI API in a Lua plugin for Pomodux. It shows a notification, prompts the user to select from a list, and asks for input.

```lua
-- Example TUI Plugin for Pomodux
-- Demonstrates usage of the TUI API (since 0.4.0)

pomodux.register_plugin({
    name = "example_tui_plugin",
    version = "1.0.0",
    description = "Shows how to use Pomodux TUI API in a plugin.",
    author = "Pomodux Team"
})

-- Register a hook for timer setup (runs before timer starts)
pomodux.register_hook("timer_setup", function(event)
    -- 1. Show a notification
    pomodux.show_notification("👋 Hello from the Example TUI Plugin!")

    -- 2. List selection
    local options = {"Red", "Green", "Blue"}
    local idx, ok = pomodux.select_from_list("Pick a color", options)
    if not ok or not idx then
        pomodux.show_notification("❌ No color selected. Cancelling timer setup.")
        error("timer setup cancelled by user")
    end
    local color = options[idx]
    pomodux.show_notification("✅ You picked: " .. color)

    -- 3. Input prompt
    local name, ok = pomodux.input_prompt("Enter your name", "", "Name")
    if not ok or name == "" then
        pomodux.show_notification("❌ No name entered. Cancelling timer setup.")
        error("timer setup cancelled by user")
    end
    pomodux.show_notification("👋 Welcome, " .. name .. "! Timer will start now.")

    -- (Optional) Log a debug message
    pomodux.log("User selected color: " .. color .. ", name: " .. name)
end)

print("✅ Example TUI Plugin loaded!")
```

**How it works:**
- When a timer is about to start, the plugin:
  1. Shows a greeting notification.
  2. Asks the user to pick a color from a list.
  3. Asks the user to enter their name.
  4. Cancels the timer setup if the user cancels any dialog.
  5. Logs the user's choices for debugging.

---

## Enabling Plugins

Enable your plugin in `config.yaml` using either style:

**New style:**
```yaml
plugins:
  example_tui_plugin:
    enabled: true
  directory: /path/to/plugins
```
**Old style (still supported):**
```yaml
plugins:
  enabled:
    example_tui_plugin: true
  directory: /path/to/plugins
```

---

## More Resources
- See `config/pomodux/plugins/` for more plugin examples.
- See [ADR 004](adr/004-plugin-system-architecture.md) for plugin architecture decisions. 