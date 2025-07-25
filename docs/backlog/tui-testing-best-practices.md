# Bubbletea TUI Testing Best Practices for Pomodux

## 1. teatest (Charmbracelet’s Official Experimental Library)
- **Repository:** https://github.com/charmbracelet/x/exp/teatest
- **Purpose:** Provides utilities for simulating user input, window resizing, and asserting on TUI output and model state.
- **Key Features:**
  - Simulate keypresses (e.g., "p", "q", "ctrl+c").
  - Simulate window resizes for responsive layouts.
  - Assert on the final model state and output (including golden file testing).
  - Useful for both unit and integration testing of Bubbletea programs.
- **Example Usage:**
  ```go
  import "github.com/charmbracelet/x/exp/teatest"

  func TestPauseResume(t *testing.T) {
      p := tea.NewProgram(NewModel(...))
      teatest.NewTestProgram(t, p).
          SendKey("p").
          AssertModel(func(m tea.Model) {
              // assert paused state
          }).
          SendKey("r").
          AssertModel(func(m tea.Model) {
              // assert running state
          })
  }
  ```
- **Status:** Experimental, but widely used in the Bubbletea community.

---

## 2. Golden File Testing
- **Strategy:** Capture the output of your TUI (as a string) and compare it to a "golden" reference file.
- **Benefits:** Detects regressions in UI rendering and layout.
- **How to Use:** Run your Bubbletea program with test inputs, capture the output, and compare to a stored file using Go’s `os` and `io/ioutil` packages or with `teatest`.

---

## 3. Standard Go Testing
- **Approach:** Test your Bubbletea model’s logic as plain Go structs and methods.
- **Best for:** State transitions, timer logic, and business rules that are independent of the UI rendering.

---

## 4. Integration with CI
- **Recommendation:** Add Bubbletea TUI tests to your CI pipeline, especially golden file and keypress simulation tests, to catch regressions early.

---

## 5. Other Community Tools
- **bubbles/testing:** Some Bubbletea component libraries (like `bubbles`) provide their own test helpers for specific widgets.
- **Testcontainers:** For plugin or integration testing, use [testcontainers-go](https://github.com/testcontainers/testcontainers-go) to spin up dependencies.

---

## Actionable Recommendations for Pomodux
1. **Add `teatest` as a test dependency:**
   ```sh
   go get github.com/charmbracelet/x/exp/teatest
   ```
2. **Write tests for your TUI:**
   - Simulate keypresses for pause, resume, stop, and quit.
   - Assert on the model’s state after each action.
   - Optionally, use golden file testing for the rendered output.
3. **Test edge cases:**
   - Window resizing.
   - Timer completion.
   - Rapid keypresses.
4. **Document your TUI testing approach** in your developer docs for future contributors.

---

This document serves as a reference for future TUI test development and review in Pomodux. 