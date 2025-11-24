package plugin

import (
	"time"

	"github.com/pomodux/pomodux/internal/logger"
)

// EventType represents the type of event
type EventType string

const (
	EventTimerStarted        EventType = "TimerStarted"
	EventTimerPaused         EventType = "TimerPaused"
	EventTimerResumed        EventType = "TimerResumed"
	EventTimerStopped        EventType = "TimerStopped"
	EventTimerCompleted      EventType = "TimerCompleted"
	EventApplicationStarted  EventType = "ApplicationStarted"
	EventApplicationStopping EventType = "ApplicationStopping"
	EventApplicationInterrupted EventType = "ApplicationInterrupted"
	EventConfigurationLoaded EventType = "ConfigurationLoaded"
)

// Event represents an application event
type Event struct {
	Type      EventType            `json:"type"`
	Timestamp time.Time            `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// Emitter handles event emission (MVP: logging only)
type Emitter struct {
	// Future: plugin manager will be added here
}

// NewEmitter creates a new event emitter
func NewEmitter() *Emitter {
	return &Emitter{}
}

// Emit emits an event (MVP: logs the event)
func (e *Emitter) Emit(eventType EventType, data map[string]interface{}) {
	event := Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}

	// MVP: Just log the event
	// Future: Dispatch to loaded plugins
	logEvent(&event)
}

// logEvent logs an event (internal helper)
func logEvent(event *Event) {
	logger.WithFields(map[string]interface{}{
		"event_type": event.Type,
		"timestamp":  event.Timestamp,
		"data":       event.Data,
	}).Info("Event emitted")
}

