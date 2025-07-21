package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pomodux/pomodux/internal/config"
	"github.com/pomodux/pomodux/internal/logger"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	lua "github.com/yuin/gopher-lua"
)

// CancellationError represents a normal user cancellation
type CancellationError struct{}

func (e CancellationError) Error() string {
	return "timer setup canceled by user"
}

// IsCancellationError checks if an error is a cancellation error
func IsCancellationError(err error) bool {
	_, ok := err.(CancellationError)
	return ok
}

// EventType represents the type of timer event
type EventType string

const (
	EventTimerSetup     EventType = "timer_setup"
	EventTimerStarted   EventType = "timer_started"
	EventTimerPaused    EventType = "timer_paused"
	EventTimerResumed   EventType = "timer_resumed"
	EventTimerCompleted EventType = "timer_completed"
	EventTimerStopped   EventType = "timer_stopped"
)

// Event represents a timer event
type Event struct {
	Type      EventType
	Timestamp time.Time
	Data      map[string]interface{}
}

// Plugin represents a loaded plugin
type Plugin struct {
	Name        string
	Version     string
	Description string
	Author      string
	LState      *lua.LState
	Hooks       map[EventType][]lua.LValue
	Enabled     bool
	mu          sync.RWMutex
}

// PluginManager manages the plugin system
type PluginManager struct {
	plugins    map[string]*Plugin
	events     chan Event
	mu         sync.RWMutex
	done       chan struct{}
	pluginsDir string
	api        *PluginAPI
	config     *config.Config // Added config field
}

// PluginAPI provides the interface for plugins to register themselves
type PluginAPI struct {
	manager *PluginManager
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(pluginsDir string, config *config.Config) *PluginManager {
	pm := &PluginManager{
		plugins:    make(map[string]*Plugin),
		events:     make(chan Event, 100),
		done:       make(chan struct{}),
		pluginsDir: pluginsDir,
		config:     config,
	}

	pm.api = &PluginAPI{manager: pm}

	// Start event processing goroutine
	go pm.processEvents()

	return pm
}

// Helper: Warn about legacy plugins
func (pm *PluginManager) warnLegacyPlugins() error {
	entries, err := os.ReadDir(pm.pluginsDir)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".lua") {
			logger.Warn("Legacy plugin file found in root of plugins directory (ignored under new structure)", map[string]interface{}{"file": entry.Name()})
		}
	}
	return nil
}

// Helper: Load plugins from config
func (pm *PluginManager) loadPluginsFromConfig() error {
	pluginNames := make([]string, 0)
	if pm.config != nil && pm.config.PluginsRaw != nil {
		for key := range pm.config.PluginsRaw {
			if key == "directory" {
				continue
			}
			pluginNames = append(pluginNames, key)
		}
	}
	for _, pluginName := range pluginNames {
		pluginDir := filepath.Join(pm.pluginsDir, pluginName)
		pluginFile := filepath.Join(pluginDir, "plugin.lua")
		// Check if plugin is enabled in configuration BEFORE loading
		enabled := true // Default to enabled
		if pm.config != nil {
			// New style: plugin-specific sub-object
			if pm.config.PluginsRaw != nil {
				if pluginSection, ok := pm.config.PluginsRaw[pluginName]; ok {
					if pluginMap, ok := pluginSection.(map[string]interface{}); ok {
						if enabledVal, ok := pluginMap["enabled"]; ok {
							if b, ok := enabledVal.(bool); ok {
								enabled = b
							}
						}
					}
				}
			}
			// Old style: enabled map
			if pm.config.Plugins.Enabled != nil {
				if enabledState, exists := pm.config.Plugins.Enabled[pluginName]; exists {
					enabled = enabledState
				}
			}
		}
		if !enabled {
			logger.Info("Skipping disabled plugin", map[string]interface{}{"plugin": pluginName})
			continue
		}
		// Only load if pluginDir exists and plugin.lua exists
		if stat, err := os.Stat(pluginFile); err == nil && !stat.IsDir() {
			if err := pm.LoadPluginFromFile(pluginFile); err != nil {
				logger.Warn("Failed to load plugin", map[string]interface{}{"path": pluginFile, "error": err.Error()})
				continue // Continue loading other plugins
			}
		} else {
			logger.Warn("Plugin folder or plugin.lua missing for plugin (skipped)", map[string]interface{}{"plugin": pluginName, "expected_path": pluginFile})
		}
	}
	return nil
}

// LoadPlugins loads all plugins from the plugins directory
func (pm *PluginManager) LoadPlugins() error {
	logger.Debug("[PLUGIN] Loading plugins from directory: " + pm.pluginsDir)
	// Create plugins directory if it doesn't exist
	if err := os.MkdirAll(pm.pluginsDir, 0750); err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}
	if err := pm.warnLegacyPlugins(); err != nil {
		return err
	}
	if err := pm.loadPluginsFromConfig(); err != nil {
		return err
	}
	return nil
}

// LoadPluginFromFile loads a plugin from a Lua file
func (pm *PluginManager) LoadPluginFromFile(filePath string) error {
	logger.Debug("[PLUGIN] Loading plugin from file: " + filePath)
	// Validate file path for security
	if err := validateFilePath(filePath); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Read the plugin file
	content, err := os.ReadFile(filePath) // #nosec G304 -- filePath is validated by validateFilePath
	if err != nil {
		return fmt.Errorf("failed to read plugin file %s: %w", filePath, err)
	}

	// Extract plugin name from filename
	pluginName := filepath.Base(filePath)
	pluginName = pluginName[:len(pluginName)-4] // Remove .lua extension
	return pm.LoadPlugin(pluginName, string(content))
}

// Helper: Extract plugin info from Lua state
func extractPluginInfo(L *lua.LState, pluginInfo *lua.LTable) (version, description, author lua.LValue) {
	version = L.GetField(pluginInfo, "version")
	description = L.GetField(pluginInfo, "description")
	author = L.GetField(pluginInfo, "author")
	return
}

// Helper: Merge pending hooks from Lua state
func mergePendingHooks(L *lua.LState, plugin *Plugin) {
	if hooksTable := L.GetGlobal("__pomodux_pending_hooks"); hooksTable.Type() == lua.LTTable {
		hooks := hooksTable.(*lua.LTable)
		hooks.ForEach(func(key lua.LValue, value lua.LValue) {
			eventType := EventType(key.String())
			if value.Type() == lua.LTFunction {
				plugin.Hooks[eventType] = append(plugin.Hooks[eventType], value)
			}
		})
		L.SetGlobal("__pomodux_pending_hooks", lua.LNil)
	}
}

func (pm *PluginManager) LoadPlugin(name, code string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if plugin already exists
	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("plugin %s already loaded", name)
	}

	// Create new Lua state for the plugin
	L := lua.NewState()

	// Register the plugin API
	pm.registerPluginAPI(L)

	// Load and run the plugin code
	if err := L.DoString(code); err != nil {
		L.Close()
		return fmt.Errorf("failed to load plugin %s: %w", name, err)
	}

	// Check if plugin registered itself
	pluginTable := L.GetGlobal("plugin")
	if pluginTable.Type() != lua.LTTable {
		L.Close()
		return fmt.Errorf("plugin %s must register itself using pomodux.register_plugin()", name)
	}

	// Extract plugin info
	version, description, author := extractPluginInfo(L, pluginTable.(*lua.LTable))

	// Check if plugin is enabled in configuration
	enabled := true // Default to enabled if not specified in config
	if pm.config != nil {
		// New style: plugin-specific sub-object
		if pm.config.PluginsRaw != nil {
			if pluginSection, ok := pm.config.PluginsRaw[name]; ok {
				if pluginMap, ok := pluginSection.(map[string]interface{}); ok {
					if enabledVal, ok := pluginMap["enabled"]; ok {
						if b, ok := enabledVal.(bool); ok {
							enabled = b
						}
					}
				}
			}
		}
		// Old style: enabled map
		if pm.config.Plugins.Enabled != nil {
			if enabledState, exists := pm.config.Plugins.Enabled[name]; exists {
				enabled = enabledState
			}
		}
	}

	// Create plugin instance
	plugin := &Plugin{
		Name:        name,
		Version:     version.String(),
		Description: description.String(),
		Author:      author.String(),
		LState:      L,
		Hooks:       make(map[EventType][]lua.LValue),
		Enabled:     enabled,
	}

	// Merge any pending hooks from Lua state
	mergePendingHooks(L, plugin)

	// Store plugin
	pm.plugins[name] = plugin

	logger.Info("Loaded plugin", map[string]interface{}{"name": name, "version": version, "author": author})
	return nil
}

// registerPluginAPI registers the plugin API functions in the Lua state
func (pm *PluginManager) registerPluginAPI(L *lua.LState) {
	// Create the pomodux table
	pomoduxTable := L.CreateTable(0, 2)
	L.SetGlobal("pomodux", pomoduxTable)

	// Register plugin registration function
	registerFn := L.NewFunction(func(L *lua.LState) int {
		pluginInfo := L.CheckTable(1)
		pm.api.registerPlugin(L, pluginInfo)
		return 0
	})
	pomoduxTable.RawSetString("register_plugin", registerFn)

	// Register hook registration function
	hookFn := L.NewFunction(func(L *lua.LState) int {
		eventType := L.CheckString(1)
		callback := L.CheckFunction(2)
		pm.api.registerHook(L, eventType, callback)
		return 0
	})
	pomoduxTable.RawSetString("register_hook", hookFn)

	// Register utility functions
	pm.registerUtilityFunctions(L, pomoduxTable)
}

// Helper: Register log function in Lua
func registerLogFn(L *lua.LState, pomoduxTable *lua.LTable) {
	logFn := L.NewFunction(func(L *lua.LState) int {
		message := L.CheckString(1)
		logger.Debug("[PLUGIN] " + message)
		return 0
	})
	pomoduxTable.RawSetString("log", logFn)
}

// Helper: Register get_config function in Lua
func registerGetConfigFn(L *lua.LState, pomoduxTable *lua.LTable) {
	getConfigFn := L.NewFunction(func(L *lua.LState) int {
		_ = L.CheckString(1) // key parameter (not used yet)
		defaultValue := L.OptString(2, "")
		// TODO: Implement config retrieval
		L.Push(lua.LString(defaultValue))
		return 1
	})
	pomoduxTable.RawSetString("get_config", getConfigFn)
}

// Helper: Register show_notification function in Lua
func registerShowNotificationFn(L *lua.LState, pomoduxTable *lua.LTable) {
	showNotificationFn := L.NewFunction(func(L *lua.LState) int {
		message := L.CheckString(1)
		logger.Debug("[PLUGIN] Lua called pomodux.show_notification with message: " + message)
		if ShowNotification(message) {
			L.Push(lua.LTrue)
		} else {
			L.Push(lua.LFalse)
		}
		return 1
	})
	pomoduxTable.RawSetString("show_notification", showNotificationFn)
}

// Helper: Register select_from_list function in Lua
func registerSelectFromListFn(L *lua.LState, pomoduxTable *lua.LTable) {
	selectFromListFn := L.NewFunction(func(L *lua.LState) int {
		title := L.CheckString(1)
		optionsTable := L.CheckTable(2)
		options := []string{}
		optionsTable.ForEach(func(_, v lua.LValue) {
			options = append(options, v.String())
		})
		idx, ok := ShowListSelection(title, options)
		if ok {
			L.Push(lua.LNumber(idx + 1)) // Lua is 1-based
			L.Push(lua.LTrue)
		} else {
			L.Push(lua.LNil)
			L.Push(lua.LFalse)
		}
		return 2
	})
	pomoduxTable.RawSetString("select_from_list", selectFromListFn)
}

// Helper: Register select_from_list_enhanced function in Lua
func registerSelectFromListEnhancedFn(L *lua.LState, pomoduxTable *lua.LTable) {
	selectFromListEnhancedFn := L.NewFunction(func(L *lua.LState) int {
		title := L.CheckString(1)
		itemsTable := L.CheckTable(2)
		items := []ListItem{}

		logger.Debug(fmt.Sprintf("[LUA] Enhanced list selection called with title: %s", title))
		logger.Debug(fmt.Sprintf("[LUA] Items table length: %d", itemsTable.Len()))

		itemsTable.ForEach(func(_, v lua.LValue) {
			if v.Type() == lua.LTTable {
				itemTable := v.(*lua.LTable)
				item := ListItem{}

				if text := L.GetField(itemTable, "text"); text.Type() == lua.LTString {
					item.Text = text.String()
				}
				if subtext := L.GetField(itemTable, "subtext"); subtext.Type() == lua.LTString {
					item.Subtext = subtext.String()
				}
				if isRecent := L.GetField(itemTable, "is_recent"); isRecent.Type() == lua.LTBool {
					item.IsRecent = bool(isRecent.(lua.LBool))
				}
				if isSeparator := L.GetField(itemTable, "is_separator"); isSeparator.Type() == lua.LTBool {
					item.IsSeparator = bool(isSeparator.(lua.LBool))
				}
				if usageCount := L.GetField(itemTable, "usage_count"); usageCount.Type() == lua.LTNumber {
					item.UsageCount = int(usageCount.(lua.LNumber))
				}

				items = append(items, item)
				logger.Debug(fmt.Sprintf("[LUA] Added item: text='%s', subtext='%s', is_separator=%v",
					item.Text, item.Subtext, item.IsSeparator))
			}
		})

		logger.Debug(fmt.Sprintf("[LUA] Calling ShowEnhancedListSelection with %d items", len(items)))
		idx, ok := ShowEnhancedListSelection(title, items)
		logger.Debug(fmt.Sprintf("[LUA] ShowEnhancedListSelection returned: idx=%d, ok=%v", idx, ok))

		if ok {
			L.Push(lua.LNumber(idx + 1)) // Lua is 1-based
			L.Push(lua.LTrue)
		} else {
			L.Push(lua.LNil)
			L.Push(lua.LFalse)
		}
		return 2
	})
	pomoduxTable.RawSetString("select_from_list_enhanced", selectFromListEnhancedFn)
}

// Helper: Register input_prompt function in Lua
func registerInputPromptFn(L *lua.LState, pomoduxTable *lua.LTable) {
	inputPromptFn := L.NewFunction(func(L *lua.LState) int {
		title := L.CheckString(1)
		defaultValue := L.OptString(2, "")
		placeholder := L.OptString(3, "")

		result, ok := ShowInputPrompt(title, defaultValue, placeholder)
		if ok {
			L.Push(lua.LString(result))
			L.Push(lua.LTrue)
		} else {
			L.Push(lua.LString(""))
			L.Push(lua.LFalse)
		}
		return 2
	})
	pomoduxTable.RawSetString("input_prompt", inputPromptFn)
}

func (pm *PluginManager) registerUtilityFunctions(L *lua.LState, pomoduxTable *lua.LTable) {
	registerLogFn(L, pomoduxTable)
	registerGetConfigFn(L, pomoduxTable)
	registerShowNotificationFn(L, pomoduxTable)
	registerSelectFromListFn(L, pomoduxTable)
	registerSelectFromListEnhancedFn(L, pomoduxTable)
	registerInputPromptFn(L, pomoduxTable)
}

// registerPlugin registers a plugin with the manager
func (api *PluginAPI) registerPlugin(L *lua.LState, pluginInfo *lua.LTable) {
	// Store the plugin info globally so the manager can access it
	L.SetGlobal("plugin", pluginInfo)

	// If there are any hooks registered before plugin registration, move them to the plugin.Hooks map
	if hooksTable := L.GetGlobal("__pomodux_pending_hooks"); hooksTable.Type() == lua.LTTable {
		// Find the plugin in the manager
		for _, plugin := range api.manager.plugins {
			if plugin.LState == L {
				plugin.mu.Lock()
				hooks := hooksTable.(*lua.LTable)
				hooks.ForEach(func(key lua.LValue, value lua.LValue) {
					eventType := EventType(key.String())
					if value.Type() == lua.LTFunction {
						plugin.Hooks[eventType] = append(plugin.Hooks[eventType], value)
					}
				})
				plugin.mu.Unlock()
			}
		}
		L.SetGlobal("__pomodux_pending_hooks", lua.LNil)
	}
}

// registerHook registers a hook for a specific event type
func (api *PluginAPI) registerHook(L *lua.LState, eventType string, callback *lua.LFunction) {
	// Find the current plugin
	pluginName := ""
	for name, plugin := range api.manager.plugins {
		if plugin.LState == L {
			pluginName = name
			break
		}
	}

	if pluginName == "" {
		// Plugin not yet registered, store in a temporary table in Lua state
		hooksTable := L.GetGlobal("__pomodux_pending_hooks")
		if hooksTable.Type() != lua.LTTable {
			hooksTable = L.CreateTable(0, 5)
			L.SetGlobal("__pomodux_pending_hooks", hooksTable)
		}
		hooksTable.(*lua.LTable).RawSetString(eventType, callback)
		return
	}

	// Register the hook
	plugin := api.manager.plugins[pluginName]
	eventTypeEnum := EventType(eventType)
	plugin.mu.Lock()
	plugin.Hooks[eventTypeEnum] = append(plugin.Hooks[eventTypeEnum], callback)
	plugin.mu.Unlock()

	logger.Info("Registered hook", map[string]interface{}{"event_type": eventType, "plugin": pluginName})
}

// EmitEvent emits an event to all registered plugins
func (pm *PluginManager) EmitEvent(event Event) {
	logger.Debug("[PLUGIN] Emitting event: " + string(event.Type))
	select {
	case pm.events <- event:
		logger.Debug("PLUGIN: Event queued for processing", map[string]interface{}{"event": event.Type})
	default:
		logger.Warn("PLUGIN: Event channel full, dropping event", map[string]interface{}{"event": event.Type})
	}
}

// EmitEventSync emits an event and calls hooks synchronously (for blocking events like timer_setup)
func (pm *PluginManager) EmitEventSync(event Event) error {
	return pm.callPluginHooks(event)
}

// processEvents processes events and calls plugin hooks
func (pm *PluginManager) processEvents() {
	for {
		select {
		case event := <-pm.events:
			if err := pm.callPluginHooks(event); err != nil {
				logger.Warn("Plugin hook error", map[string]interface{}{"error": err.Error()})
			}
		case <-pm.done:
			return
		}
	}
}

// callPluginHooks calls all registered hooks for an event and returns error if any fail
func (pm *PluginManager) callPluginHooks(event Event) error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	for _, plugin := range pm.plugins {
		// Skip disabled plugins
		if !plugin.Enabled {
			continue
		}
		plugin.mu.RLock()
		hooks := plugin.Hooks[event.Type]
		plugin.mu.RUnlock()
		for _, hook := range hooks {
			if err := pm.callHook(plugin, hook, event); err != nil {
				return err
			}
		}
	}
	return nil
}

// callHook calls a single plugin hook and returns error if the hook fails or returns false
func (pm *PluginManager) callHook(plugin *Plugin, hook lua.LValue, event Event) error {
	L := plugin.LState
	// Push the hook function and event argument
	L.Push(hook)
	L.Push(pm.eventToLuaTable(L, event))
	// Call the function with 1 argument, 1 result
	err := L.PCall(1, 1, nil)
	if err != nil {
		return fmt.Errorf("failed to call hook: %w", err)
	}
	ret := L.Get(-1)
	L.Pop(1)

	// Check for cancellation (false return)
	if ret.Type() == lua.LTBool && ret == lua.LFalse {
		return CancellationError{}
	}

	// Check for session data modifications (table return)
	if ret.Type() == lua.LTTable {
		modifications := ret.(*lua.LTable)
		modifications.ForEach(func(key lua.LValue, value lua.LValue) {
			keyStr := key.String()
			switch value.Type() {
			case lua.LTString:
				event.Data[keyStr] = value.String()
			case lua.LTNumber:
				event.Data[keyStr] = int(value.(lua.LNumber))
			case lua.LTBool:
				event.Data[keyStr] = bool(value.(lua.LBool))
			case lua.LTNil, lua.LTFunction, lua.LTUserData, lua.LTThread, lua.LTTable, lua.LTChannel:
				// No-op for these types, but handled for exhaustiveness
			default:
				// No-op
			}
		})
	}

	return nil
}

// GetPlugin returns a plugin by name
func (pm *PluginManager) GetPlugin(name string) (*Plugin, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	plugin, exists := pm.plugins[name]
	return plugin, exists
}

// ListPlugins returns all loaded plugins
func (pm *PluginManager) ListPlugins() []*Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugins := make([]*Plugin, 0, len(pm.plugins))
	for _, plugin := range pm.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// EnablePlugin enables or disables a plugin
func (pm *PluginManager) EnablePlugin(name string, enabled bool) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	plugin.Enabled = enabled
	return nil
}

// UnloadPlugin unloads a plugin
func (pm *PluginManager) UnloadPlugin(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Close the Lua state
	plugin.LState.Close()

	// Remove from plugins map
	delete(pm.plugins, name)

	return nil
}

// Shutdown shuts down the plugin manager
func (pm *PluginManager) Shutdown() {
	close(pm.done)

	// Unload all plugins
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for name, plugin := range pm.plugins {
		plugin.LState.Close()
		delete(pm.plugins, name)
	}
}

// validateFilePath validates a file path for security
func validateFilePath(filePath string) error {
	// Check for path traversal attempts
	if strings.Contains(filePath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	// Check for dangerous characters
	dangerousChars := []string{"|", "&", ";", "`", "$", "(", ")", "<", ">", "*", "?"}
	for _, char := range dangerousChars {
		if strings.Contains(filePath, char) {
			return fmt.Errorf("dangerous character '%s' not allowed in file path", char)
		}
	}

	// For absolute paths, check if they're in a safe location
	if filepath.IsAbs(filePath) {
		// Allow paths in /tmp for testing
		if strings.HasPrefix(filePath, "/tmp/") {
			return nil
		}

		// Allow paths in the user's home directory
		homeDir, err := os.UserHomeDir()
		if err == nil {
			if strings.HasPrefix(filePath, homeDir) {
				return nil
			}
		}

		// For now, be restrictive with other absolute paths
		return fmt.Errorf("absolute paths not allowed for security")
	}

	return nil
}

// ShowNotification displays a modal notification using tview and returns true if OK, false if Cancel
func ShowNotification(message string) bool {
	logger.Debug("[PLUGIN] ShowNotification called with message: " + message)
	result := false
	app := tview.NewApplication()
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				result = true
			}
			app.Stop()
		})
	_ = app.SetRoot(modal, false).Run()
	return result
}

// ShowListSelection displays a modal list selection dialog using tview
func ShowListSelection(title string, options []string) (int, bool) {
	selected := -1
	confirmed := false
	app := tview.NewApplication()
	list := tview.NewList()
	for i, opt := range options {
		list.AddItem(opt, "", 0, func(idx int) func() {
			return func() {
				selected = idx
				confirmed = true
				app.Stop()
			}
		}(i))
	}
	list.SetTitle(title).SetBorder(true)
	list.SetDoneFunc(func() {
		app.Stop()
	})
	_ = app.SetRoot(list, true).Run()
	return selected, confirmed
}

// ListItem represents an item in an enhanced list selection
type ListItem struct {
	Text        string
	Subtext     string
	IsRecent    bool
	IsSeparator bool
	UsageCount  int
}

// ShowEnhancedListSelection displays a modal list selection dialog with enhanced items
func ShowEnhancedListSelection(title string, items []ListItem) (int, bool) {
	selected := -1
	confirmed := false
	app := tview.NewApplication()
	list := tview.NewList()

	// Add debugging
	logger.Debug(fmt.Sprintf("[TUI] Enhanced list selection: %s with %d items", title, len(items)))

	// Debug: Print all items
	for i, item := range items {
		logger.Debug(fmt.Sprintf("[TUI] Item %d: text='%s', subtext='%s', is_separator=%v",
			i, item.Text, item.Subtext, item.IsSeparator))
	}

	for i, item := range items {
		if item.IsSeparator {
			// Add separator item
			separatorText := "─ " + item.Text + " ─"
			list.AddItem(separatorText, "", 0, nil)
			logger.Debug(fmt.Sprintf("[TUI] Added separator item %d: %s", i, separatorText))
		} else {
			// Add regular item with subtext
			displayText := item.Text
			if item.Subtext != "" {
				displayText = displayText + " (" + item.Subtext + ")"
			}
			list.AddItem(displayText, "", 0, func(idx int) func() {
				return func() {
					logger.Debug(fmt.Sprintf("[TUI] User clicked on item %d: %s", idx, items[idx].Text))
					selected = idx
					confirmed = true
					logger.Debug(fmt.Sprintf("[TUI] Setting selected=%d, confirmed=%v", selected, confirmed))
					app.Stop()
				}
			}(i))
			logger.Debug(fmt.Sprintf("[TUI] Added selectable item %d: %s", i, displayText))
		}
	}

	list.SetTitle(title).SetBorder(true)

	// Add keyboard shortcuts for better UX
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		logger.Debug(fmt.Sprintf("[TUI] Key pressed: %v", event.Key()))
		if event.Key() == tcell.KeyEscape {
			logger.Debug("[TUI] User pressed Escape - cancelling selection")
			app.Stop()
			return nil
		}
		return event
	})

	logger.Debug("[TUI] Starting enhanced list selection dialog")
	_ = app.SetRoot(list, true).Run()

	logger.Debug(fmt.Sprintf("[TUI] Enhanced list selection result: selected=%d, confirmed=%v", selected, confirmed))
	return selected, confirmed
}

// ShowInputPrompt displays a modal input prompt dialog
func ShowInputPrompt(title, defaultValue, placeholder string) (string, bool) {
	result := ""
	confirmed := false
	app := tview.NewApplication()

	// Create a simple form with title, input, and buttons
	form := tview.NewForm().
		AddTextView("Title", title, 0, 1, true, false).
		AddInputField("Input", defaultValue, 50, nil, func(text string) {
			result = text
		}).
		AddButton("OK", func() {
			confirmed = true
			app.Stop()
		}).
		AddButton("Cancel", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Input Prompt")
	_ = app.SetRoot(form, true).Run()
	return result, confirmed
}

// eventToLuaTable converts a Go Event to a Lua table for Lua hooks
func (pm *PluginManager) eventToLuaTable(L *lua.LState, event Event) *lua.LTable {
	eventTable := L.CreateTable(0, 3)
	eventTable.RawSetString("type", lua.LString(event.Type))
	eventTable.RawSetString("timestamp", lua.LNumber(event.Timestamp.Unix()))
	dataTable := L.CreateTable(0, len(event.Data))
	for k, v := range event.Data {
		switch val := v.(type) {
		case string:
			dataTable.RawSetString(k, lua.LString(val))
		case int:
			dataTable.RawSetString(k, lua.LNumber(val))
		case float64:
			dataTable.RawSetString(k, lua.LNumber(val))
		case bool:
			dataTable.RawSetString(k, lua.LBool(val))
		default:
			dataTable.RawSetString(k, lua.LString(fmt.Sprintf("%v", val)))
		}
	}
	eventTable.RawSetString("data", dataTable)
	return eventTable
}
