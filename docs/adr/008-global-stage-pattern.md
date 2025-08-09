---
status: approved
type: technical
---

# ADR 008: Global Stage Pattern for TUI Component Management

## 1. Context / Background

Future TUI architecture will introduce a TUI-only design that eliminates cross-process synchronization by running all timer operations within a single Bubbletea application. This requires a unified approach to managing multiple UI components, including:

- **Timer Window**: Primary timer display and controls
- **Plugin Modals**: Plugin-driven dialog windows
- **Notifications**: Auto-dismissing notification messages
- **Status Panel**: Plugin information display areas

The current TUI implementation uses a simple single-component model that doesn't support multiple windows or component management. A more sophisticated architecture is needed to handle:

- Window layering and z-order management
- Event distribution across components
- Component lifecycle management
- Responsive layout coordination

## 2. Decision

**Pomodux will implement a Global Stage Pattern for TUI component management, using a singleton stage manager to coordinate all UI components within the Bubbletea application.**

### Core Architecture
- **Global Stage Singleton**: Central manager for all UI components
- **Component Registration**: Components register with the stage for lifecycle management
- **Event Distribution**: Stage distributes events to appropriate components
- **Window Management**: Stage manages window layering, visibility, and z-order
- **Responsive Layout**: Stage coordinates layout changes across all components

### Component Types
1. **Timer Window**: Primary component, always visible
2. **Plugin Modals**: Overlay windows for plugin dialogs
3. **Notifications**: Temporary overlay messages
4. **Status Panel**: Information display areas

### Implementation Details
```go
type GlobalStage struct {
    timerWindow     *TimerWindow
    pluginModals    []*PluginModal
    notifications   []*Notification
    statusPanel     *StatusPanel
    mu              sync.RWMutex
    eventBus        chan StageEvent
}
```

## 3. Rationale

### **Unified Component Management**
- Single point of control for all UI components
- Consistent event handling and lifecycle management
- Simplified component coordination and communication

### **Scalable Architecture**
- Easy to add new component types
- Flexible window management system
- Support for complex UI interactions

### **Bubbletea Integration**
- Leverages Bubbletea's reactive model
- Maintains single-threaded event processing
- Compatible with existing Bubbletea patterns

### **Plugin System Support**
- Enables plugin-driven UI components
- Supports modal window spawning rules
- Provides consistent plugin UI experience

### **Performance Benefits**
- Efficient event distribution
- Minimal component re-rendering
- Optimized layout calculations

## 4. Alternatives Considered

### **Component Tree Pattern**
- **Description**: Hierarchical component tree with parent-child relationships
- **Rejected**: Overly complex for Pomodux's needs, difficult to manage plugin components

### **Event Bus Only**
- **Description**: Pure event-driven architecture without central management
- **Rejected**: Lacks coordination for window management and layout

### **Multiple Bubbletea Programs**
- **Description**: Separate Bubbletea programs for different UI areas
- **Rejected**: Violates single-process architecture, complex inter-process communication

### **Simple Component Array**
- **Description**: Basic array of components without central management
- **Rejected**: Insufficient for window management and event coordination

## 5. Consequences

### **Positive Consequences**
- **Unified UI Management**: Single point of control for all components
- **Plugin Integration**: Seamless plugin UI component integration
- **Consistent UX**: Uniform window management and interaction patterns
- **Maintainability**: Clear separation of concerns and component responsibilities
- **Extensibility**: Easy to add new component types and features

### **Negative Consequences**
- **Complexity**: More complex than simple single-component model
- **Learning Curve**: Developers must understand stage management patterns
- **Testing Complexity**: More complex testing scenarios for component interactions
- **Performance Overhead**: Small overhead for stage management operations

### **Mitigation Strategies**
- **Comprehensive Documentation**: Clear documentation of stage patterns and usage
- **Testing Framework**: Robust testing using teatest framework
- **Performance Monitoring**: Monitor stage operations for performance impact
- **Incremental Implementation**: Implement stage pattern incrementally

## 6. Implementation Plan

### **Phase 1: Core Stage Implementation**
- Implement GlobalStage singleton with basic component management
- Add TimerWindow as primary component
- Implement basic event distribution system

### **Phase 2: Window Management**
- Add plugin modal support with window layering
- Implement notification system with auto-dismiss
- Add status panel for plugin information

### **Phase 3: Advanced Features**
- Implement responsive layout coordination
- Add component lifecycle management
- Optimize performance and event handling

### **Testing Strategy**
- Unit tests for stage management operations
- Integration tests for component interactions
- teatest framework for end-to-end UI testing
- Performance tests for stage operations

## 7. Status

- **Approved** (2025-01-27)
- **Alternative Implementation** (2025-01-09): TUI Timer Feature implemented using simpler Bubbletea model approach
- **Recommendation**: Global Stage Pattern deferred for future multi-component UI needs

### Implementation Notes
The TUI Timer Feature (2025-01-09) successfully implemented timer functionality using a simplified Bubbletea model that meets current requirements without the full Global Stage Pattern complexity. The existing implementation provides:
- Centered, responsive timer window
- Session name display and progress tracking  
- Keyboard controls and visual feedback
- Plugin integration within single process

The Global Stage Pattern remains valuable for future enhancements requiring multiple UI components (plugin modals, notifications, status panels) but is not required for current timer functionality.

## 8. References

- [ADR 007: TUI Standardization](007-tui-standardization.md)
- [Bubbletea Documentation](https://github.com/charmbracelet/bubbletea)
- [teatest Framework](https://github.com/charmbracelet/x/exp/teatest)

## 9. Related Components

### **Affected Components**
- `internal/tui/stage.go` - Global stage implementation
- `internal/tui/timer_window.go` - Timer window component
- `internal/tui/plugin_modal.go` - Plugin modal component
- `internal/tui/notification.go` - Notification component
- `internal/tui/status_panel.go` - Status panel component

### **Integration Points**
- Plugin API for modal window spawning
- Event system for component communication
- Layout system for responsive design
- Testing framework for component validation 