# Release 0.4.1: Plugin Loader Overhaul & Kimai Integration Improvements

> **Release**: 0.4.1 - Plugin Loader Overhaul & Kimai Integration Improvements
> **Status**: ✅ **RELEASED**
> **Release Date**: 2025-07-21

## Release Overview

This release delivers a major overhaul of the plugin loading system and significant improvements to the Kimai integration, along with enhanced logging and documentation.

---

## 🎯 Key Features & Fixes

### 1. Plugin Loader Refactor
- Only loads plugins from subfolders listed in the config.
- Only loads `plugin.lua` from each subfolder.
- Ignores and logs a warning for any `.lua` files in the root of the plugins directory (legacy plugins).

### 2. Kimai Plugin Robustness
- Project/activity selection now supports `< Back>` navigation for user-friendly flow.
- Timer is only started after both project and activity are selected and confirmed.
- Kimai timer start/stop is reliably synchronized with Pomodux timer events.
- All user cancellations are silent in the CLI but logged for audit/debugging.

### 3. Logging & Output
- All backend/plugin status, warning, and error output is routed through the logger (not print).
- Timer UI output remains user-friendly and unchanged.

### 4. Documentation & Migration
- Plugin development guide updated for new structure, cancellation, and TUI API patterns.
- Migration instructions: Move plugins to subfolders, ensure only `plugin.lua` is loaded, update any legacy plugin locations.

---

## 🛠 Migration Instructions
- Move all plugins to their own subfolders in the plugins directory.
- Ensure each plugin subfolder contains a `plugin.lua` as the entry point.
- Remove or archive any `.lua` files in the root of the plugins directory.
- Update any custom plugins to use the new cancellation and logging patterns.

---

## ✅ Quality Assurance
- All tests updated for new plugin structure and warnings.
- Manual and automated UAT completed for Kimai integration and plugin loader changes.

---

## 📢 Announcement
- The Pomodux plugin system is now more robust, maintainable, and user-friendly.
- Kimai integration is seamless and reliable for time tracking workflows. 