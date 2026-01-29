package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "nonexistent.json")

	h, err := Load(path)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "1.0", h.Version)
	assert.Empty(t, h.Sessions)
}

func TestLoad_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "history.json")

	// Create valid history file
	content := `{
  "version": "1.0",
  "sessions": [
    {
      "id": "test-id-1",
      "started_at": "2025-01-15T14:00:00Z",
      "ended_at": "2025-01-15T14:25:00Z",
      "duration": "25m",
      "preset": "work",
      "label": "Test session",
      "end_status": "completed",
      "paused_count": 0,
      "paused_duration": "0s"
    }
  ]
}`
	err := os.WriteFile(path, []byte(content), 0600)
	require.NoError(t, err)

	h, err := Load(path)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "1.0", h.Version)
	require.Len(t, h.Sessions, 1)
	assert.Equal(t, "test-id-1", h.Sessions[0].ID)
	assert.Equal(t, "25m", h.Sessions[0].Duration)
	assert.Equal(t, "work", h.Sessions[0].Preset)
	assert.Equal(t, "Test session", h.Sessions[0].Label)
	assert.Equal(t, "completed", h.Sessions[0].EndStatus)
	assert.Equal(t, 0, h.Sessions[0].PausedCount)
	assert.Equal(t, "0s", h.Sessions[0].PausedDuration)
}

func TestLoad_UnreadablePath(t *testing.T) {
	tmpDir := t.TempDir()
	// Use a path that is a directory so ReadFile fails
	path := tmpDir

	h, err := Load(path)
	require.Error(t, err)
	assert.Nil(t, h)
	assert.Contains(t, err.Error(), "failed to read history file")
}

func TestLoad_CorruptFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "history.json")
	corruptContent := []byte(`{ invalid json }`)
	err := os.WriteFile(path, corruptContent, 0600)
	require.NoError(t, err)

	h, err := Load(path)
	require.Error(t, err)
	assert.Nil(t, h)
	assert.Contains(t, err.Error(), "failed to parse history file")
	assert.Contains(t, err.Error(), ".backup")

	// Backup file should exist
	backupPath := path + ".backup"
	_, err = os.Stat(backupPath)
	assert.NoError(t, err)
	backupData, err := os.ReadFile(backupPath)
	require.NoError(t, err)
	assert.Equal(t, corruptContent, backupData)
}

func TestSave_CreatesDirectoryAndFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "subdir", "history.json")

	h := &History{
		Version:  "1.0",
		Sessions: []Session{},
	}
	err := Save(h, path)
	require.NoError(t, err)

	_, err = os.Stat(path)
	require.NoError(t, err)

	// Verify content by loading
	loaded, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "1.0", loaded.Version)
	assert.Empty(t, loaded.Sessions)
}

func TestSave_AtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "history.json")

	h := &History{
		Version: "1.0",
		Sessions: []Session{
			{
				ID:             "atomic-test",
				StartedAt:      time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC),
				EndedAt:        time.Date(2025, 1, 15, 14, 25, 0, 0, time.UTC),
				Duration:       "25m",
				Preset:         "work",
				Label:          "Atomic test",
				EndStatus:      "completed",
				PausedCount:    0,
				PausedDuration: "0s",
			},
		},
	}
	err := Save(h, path)
	require.NoError(t, err)

	// Temp file should not exist after rename
	tmpPath := path + ".tmp"
	_, err = os.Stat(tmpPath)
	assert.True(t, os.IsNotExist(err))

	loaded, err := Load(path)
	require.NoError(t, err)
	require.Len(t, loaded.Sessions, 1)
	assert.Equal(t, "atomic-test", loaded.Sessions[0].ID)
}

func TestAddSession(t *testing.T) {
	h := &History{
		Version:  "1.0",
		Sessions: []Session{},
	}
	session := Session{
		ID:             "new-session",
		StartedAt:      time.Now().Add(-25 * time.Minute),
		EndedAt:        time.Now(),
		Duration:       "25m",
		Label:          "New session",
		EndStatus:      "completed",
		PausedCount:    1,
		PausedDuration: "2m",
	}

	h.AddSession(session)
	require.Len(t, h.Sessions, 1)
	assert.Equal(t, "new-session", h.Sessions[0].ID)

	h.AddSession(Session{ID: "second", Label: "Second", EndStatus: "stopped", PausedCount: 0, PausedDuration: "0s"})
	require.Len(t, h.Sessions, 2)
	assert.Equal(t, "second", h.Sessions[1].ID)
}

func TestGetRecent_EmptyHistory(t *testing.T) {
	h := &History{Version: "1.0", Sessions: []Session{}}
	result := h.GetRecent(5)
	assert.Empty(t, result)
}

func TestGetRecent_WithinLimit(t *testing.T) {
	h := &History{
		Version: "1.0",
		Sessions: []Session{
			{ID: "oldest", Label: "A", EndStatus: "completed", PausedCount: 0, PausedDuration: "0s"},
			{ID: "middle", Label: "B", EndStatus: "completed", PausedCount: 0, PausedDuration: "0s"},
			{ID: "newest", Label: "C", EndStatus: "completed", PausedCount: 0, PausedDuration: "0s"},
		},
	}
	result := h.GetRecent(2)
	require.Len(t, result, 2)
	// Newest first: C then B
	assert.Equal(t, "newest", result[0].ID)
	assert.Equal(t, "middle", result[1].ID)
}

func TestGetRecent_MoreThanLimit(t *testing.T) {
	h := &History{
		Version: "1.0",
		Sessions: []Session{
			{ID: "1", Label: "A", EndStatus: "completed", PausedCount: 0, PausedDuration: "0s"},
			{ID: "2", Label: "B", EndStatus: "completed", PausedCount: 0, PausedDuration: "0s"},
			{ID: "3", Label: "C", EndStatus: "completed", PausedCount: 0, PausedDuration: "0s"},
		},
	}
	result := h.GetRecent(10)
	require.Len(t, result, 3)
	assert.Equal(t, "3", result[0].ID)
	assert.Equal(t, "2", result[1].ID)
	assert.Equal(t, "1", result[2].ID)
}

func TestGetRecent_ZeroOrNegativeLimit(t *testing.T) {
	h := &History{
		Version: "1.0",
		Sessions: []Session{
			{ID: "a", Label: "A", EndStatus: "completed", PausedCount: 0, PausedDuration: "0s"},
		},
	}
	// limit <= 0 means use len(h.Sessions) per implementation
	result := h.GetRecent(0)
	require.Len(t, result, 1)
	assert.Equal(t, "a", result[0].ID)
}

func TestRoundTrip_LoadSaveLoad(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "history.json")

	started := time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC)
	ended := time.Date(2025, 1, 15, 14, 25, 0, 0, time.UTC)
	original := &History{
		Version: "1.0",
		Sessions: []Session{
			{
				ID:             "roundtrip-id",
				StartedAt:      started,
				EndedAt:        ended,
				Duration:       "25m",
				Preset:         "work",
				Label:          "Roundtrip session",
				EndStatus:      "stopped",
				PausedCount:    2,
				PausedDuration: "3m",
			},
		},
	}
	err := Save(original, path)
	require.NoError(t, err)

	loaded, err := Load(path)
	require.NoError(t, err)
	require.Len(t, loaded.Sessions, 1)
	s := loaded.Sessions[0]
	assert.Equal(t, "roundtrip-id", s.ID)
	assert.True(t, s.StartedAt.Equal(started))
	assert.True(t, s.EndedAt.Equal(ended))
	assert.Equal(t, "25m", s.Duration)
	assert.Equal(t, "work", s.Preset)
	assert.Equal(t, "Roundtrip session", s.Label)
	assert.Equal(t, "stopped", s.EndStatus)
	assert.Equal(t, 2, s.PausedCount)
	assert.Equal(t, "3m", s.PausedDuration)
}
