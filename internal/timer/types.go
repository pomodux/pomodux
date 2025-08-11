package timer

// TimerStatus represents the current status of a timer
type TimerStatus string

const (
	// StatusIdle indicates the timer is not running
	StatusIdle TimerStatus = "idle"
	
	// StatusRunning indicates the timer is currently running
	StatusRunning TimerStatus = "running"
	
	// StatusPaused indicates the timer is paused
	StatusPaused TimerStatus = "paused"
	
	// StatusCompleted indicates the timer has finished
	StatusCompleted TimerStatus = "completed"
)

// String returns the string representation of the timer status
func (ts TimerStatus) String() string {
	return string(ts)
}
