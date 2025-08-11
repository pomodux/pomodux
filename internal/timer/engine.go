package timer

import (
	"context"
	"sync"
	"time"

	"github.com/pomodux/pomodux/internal/logger"
)

// TimerEventType represents different types of timer events
type TimerEventType string

const (
	EventTimerStart    TimerEventType = "timer_start"
	EventTimerPause    TimerEventType = "timer_pause"
	EventTimerResume   TimerEventType = "timer_resume"
	EventTimerComplete TimerEventType = "timer_complete"
	EventTimerStop     TimerEventType = "timer_stop"
	EventProgressTick  TimerEventType = "progress_tick" // For smooth animation
)

// TimerEvent represents a timer event with progress information
type TimerEvent struct {
	Type          TimerEventType `json:"type"`
	Timestamp     time.Time      `json:"timestamp"`
	Progress      float64        `json:"progress"` // 0.0 to 1.0
	TimeElapsed   time.Duration  `json:"time_elapsed"`
	TimeRemaining time.Duration  `json:"time_remaining"`
	Status        TimerStatus    `json:"status"`
	SessionName   string         `json:"session_name"`
}

// EventDrivenTimer wraps the existing Timer with event-driven capabilities
type EventDrivenTimer struct {
	*Timer // Embed existing timer

	subscribers []chan TimerEvent
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
}

// NewEventDrivenTimer creates a new event-driven timer
func NewEventDrivenTimer() *EventDrivenTimer {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventDrivenTimer{
		Timer:       NewTimer(),
		subscribers: make([]chan TimerEvent, 0),
		ctx:         ctx,
		cancel:      cancel,
		mu:          sync.RWMutex{},
	}
}

// NewEventDrivenTimerWithManagers creates a new event-driven timer with managers
func NewEventDrivenTimerWithManagers(stateManager *StateManager, historyManager *HistoryManager) *EventDrivenTimer {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventDrivenTimer{
		Timer:       NewTimerWithManagers(stateManager, historyManager),
		subscribers: make([]chan TimerEvent, 0),
		ctx:         ctx,
		cancel:      cancel,
		mu:          sync.RWMutex{},
	}
}

// Subscribe returns a channel that receives timer events
func (edt *EventDrivenTimer) Subscribe() <-chan TimerEvent {
	edt.mu.Lock()
	defer edt.mu.Unlock()

	ch := make(chan TimerEvent, 50) // Buffered for responsiveness
	edt.subscribers = append(edt.subscribers, ch)

	logger.Debug("Timer event subscriber added", map[string]interface{}{
		"subscribers_count": len(edt.subscribers),
	})

	return ch
}

// Unsubscribe removes a subscriber channel
func (edt *EventDrivenTimer) Unsubscribe(ch <-chan TimerEvent) {
	edt.mu.Lock()
	defer edt.mu.Unlock()

	// Find and remove the channel
	for i, sub := range edt.subscribers {
		if sub == ch {
			// Close the channel and remove from slice
			close(edt.subscribers[i])
			edt.subscribers = append(edt.subscribers[:i], edt.subscribers[i+1:]...)
			logger.Debug("Timer event subscriber removed", map[string]interface{}{
				"subscribers_count": len(edt.subscribers),
			})
			break
		}
	}
}

// emitEvent sends an event to all subscribers
func (edt *EventDrivenTimer) emitEvent(eventType TimerEventType) {
	event := TimerEvent{
		Type:          eventType,
		Timestamp:     time.Now(),
		Progress:      edt.GetProgress(),
		TimeElapsed:   edt.GetElapsed(),
		TimeRemaining: edt.GetDuration() - edt.GetElapsed(),
		Status:        edt.GetStatus(),
		SessionName:   edt.GetSessionName(),
	}

	edt.mu.RLock()
	subscribers := make([]chan TimerEvent, len(edt.subscribers))
	copy(subscribers, edt.subscribers)
	edt.mu.RUnlock()

	logger.Debug("Emitting timer event", map[string]interface{}{
		"event_type":        string(eventType),
		"progress":          event.Progress,
		"time_remaining":    event.TimeRemaining,
		"subscribers_count": len(subscribers),
	})

	// Non-blocking emit to all subscribers
	for _, ch := range subscribers {
		select {
		case ch <- event:
			// Event sent successfully
		default:
			// Channel full, subscriber is slow - skip this update
			logger.Warn("Timer event channel full, skipping update", map[string]interface{}{
				"event_type": string(eventType),
			})
		}
	}
}

// StartWithSessionName begins the timer and emits start event
func (edt *EventDrivenTimer) StartWithSessionName(duration time.Duration, sessionName string) error {
	if err := edt.Timer.StartWithSessionName(duration, sessionName); err != nil {
		return err
	}

	// Emit immediate start event
	edt.emitEvent(EventTimerStart)

	return nil
}

// Pause pauses the timer and emits pause event
func (edt *EventDrivenTimer) Pause() error {
	if err := edt.Timer.Pause(); err != nil {
		return err
	}

	// Emit pause event
	edt.emitEvent(EventTimerPause)
	return nil
}

// Resume resumes the timer and emits resume event
func (edt *EventDrivenTimer) Resume() error {
	if err := edt.Timer.Resume(); err != nil {
		return err
	}

	// Emit resume event
	edt.emitEvent(EventTimerResume)
	return nil
}

// Stop stops the timer and emits stop event
func (edt *EventDrivenTimer) Stop() error {
	if err := edt.Timer.Stop(); err != nil {
		return err
	}

	// Emit stop event
	edt.emitEvent(EventTimerStop)
	return nil
}

// CheckProgress checks timer status and emits events as needed (call from TUI tick)
func (edt *EventDrivenTimer) CheckProgress() {
	status := edt.GetStatus()

	// Check if timer completed
	if status == StatusCompleted {
		logger.Info("Timer completed, emitting completion event")
		edt.emitEvent(EventTimerComplete)
		return
	}

	// Emit progress tick for running timers
	if status == StatusRunning {
		edt.emitEvent(EventProgressTick)
	}
}

// Close cleans up the event-driven timer
func (edt *EventDrivenTimer) Close() {
	logger.Debug("Closing event-driven timer")

	// Cancel context
	edt.cancel()

	// Close all subscriber channels
	edt.mu.Lock()
	for _, ch := range edt.subscribers {
		close(ch)
	}
	edt.subscribers = nil
	edt.mu.Unlock()

	logger.Debug("Event-driven timer closed")
}
