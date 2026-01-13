#!/bin/bash
# Test script for singleton timer functionality
# Run this in an interactive terminal

set -e

echo "=== Pomodux Timer Test Script ==="
echo ""

# Clean up any existing state
rm -f ~/.local/state/pomodux/timer_state.json
echo "✅ Cleaned up existing state"

echo ""
echo "Test 1: Start a timer (will run for 5 seconds)"
echo "Press Ctrl+C to interrupt, or wait for completion"
echo ""
./bin/pomodux start 5s "Test timer 1"

echo ""
echo "Test 2: Try to start a second timer while first is running"
echo "In another terminal, run: ./bin/pomodux start 10s 'Second timer'"
echo "Expected: Error 'timer already running in process X'"
echo ""
read -p "Press Enter after testing in second terminal..."

echo ""
echo "Test 3: Crash recovery"
echo "Starting a timer, then kill it with: kill <PID>"
echo "Then start again - should auto-resume"
./bin/pomodux start 30s "Crash test timer" &
TIMER_PID=$!
echo "Timer started with PID: $TIMER_PID"
echo "Kill it with: kill $TIMER_PID"
read -p "Press Enter after killing the process..."
echo "Now starting again - should auto-resume:"
./bin/pomodux start 10s "Should resume"

echo ""
echo "✅ All tests completed!"
