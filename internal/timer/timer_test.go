package timer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTimer_Valid(t *testing.T) {
	timer, err := NewTimer(25*time.Minute, "Test timer", "work")
	assert.NoError(t, err)
	assert.NotNil(t, timer)
	assert.Equal(t, 25*time.Minute, timer.Duration())
	assert.Equal(t, "Test timer", timer.Label())
	assert.Equal(t, "work", timer.Preset())
	assert.Equal(t, StateIdle, timer.State())
}

func TestNewTimer_InvalidDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
	}{
		{"zero duration", 0},
		{"negative duration", -5 * time.Minute},
		{"too large", 25 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer, err := NewTimer(tt.duration, "Test", "")
			assert.Error(t, err)
			assert.Nil(t, timer)
		})
	}
}

func TestTimer_Start(t *testing.T) {
	timer, _ := NewTimer(25*time.Minute, "Test", "")
	
	err := timer.Start()
	assert.NoError(t, err)
	assert.Equal(t, StateRunning, timer.State())
	assert.False(t, timer.StartTime().IsZero())
}

func TestTimer_PauseResume(t *testing.T) {
	timer, _ := NewTimer(25*time.Minute, "Test", "")
	timer.Start()

	// Pause
	err := timer.Pause()
	assert.NoError(t, err)
	assert.Equal(t, StatePaused, timer.State())
	assert.Equal(t, 1, timer.PausedCount())

	// Resume
	err = timer.Resume()
	assert.NoError(t, err)
	assert.Equal(t, StateRunning, timer.State())
}

func TestTimer_Remaining(t *testing.T) {
	timer, _ := NewTimer(1*time.Second, "Test", "")
	timer.Start()

	remaining := timer.Remaining()
	assert.True(t, remaining > 0)
	assert.True(t, remaining <= 1*time.Second)
}

func TestTimer_Stop(t *testing.T) {
	timer, _ := NewTimer(25*time.Minute, "Test", "")
	timer.Start()

	err := timer.Stop()
	assert.NoError(t, err)
	assert.Equal(t, StateStopped, timer.State())
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		d        time.Duration
		expected string
	}{
		{"zero", 0, "0s"},
		{"seconds only", 45 * time.Second, "45s"},
		{"minutes", 25 * time.Minute, "25m"},
		{"hours and minutes", 90 * time.Minute, "1h30m"},
		{"hours only", 2 * time.Hour, "2h"},
		{"mixed", 2*time.Hour + 15*time.Minute, "2h15m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDuration(tt.d)
			assert.Equal(t, tt.expected, got)
		})
	}
}


