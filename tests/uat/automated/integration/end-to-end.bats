#!/usr/bin/env bats

# BATS integration test suite for Pomodux End-to-End Workflows
# Requires: bats-core (https://github.com/bats-core/bats-core)
# Run with: bats tests/uat/automated/integration/end-to-end.bats

setup() {
    # Setup test environment
    export APP_BINARY="${BATS_TEST_DIRNAME}/../../../bin/pomodux"
    export CONFIG_DIR="${HOME}/.config/pomodux"
    
    # Create test config directory
    mkdir -p "$CONFIG_DIR"
    
    # Backup existing config if it exists
    if [ -f "${CONFIG_DIR}/config.json" ]; then
        cp "${CONFIG_DIR}/config.json" "${CONFIG_DIR}/config.json.backup"
    fi
    
    # Create test configuration with short durations for testing
    cat > "${CONFIG_DIR}/config.json" << EOF
{
    "default_work_duration": "1m",
    "default_break_duration": "30s",
    "default_long_break_duration": "2m"
}
EOF
    
    # Ensure application is built
    cd "${BATS_TEST_DIRNAME}/../../../"
    make build > /dev/null 2>&1 || true
    
    # Clean up any existing history
    rm -f "${CONFIG_DIR}/session_history.json"
}

teardown() {
    # Clean up any running timer state (stop command removed from CLI)
    # Timer will auto-complete or can be stopped via TUI if needed
    
    # Restore original config
    if [ -f "${CONFIG_DIR}/config.json.backup" ]; then
        mv "${CONFIG_DIR}/config.json.backup" "${CONFIG_DIR}/config.json"
    else
        rm -f "${CONFIG_DIR}/config.json"
    fi
    
    # Clean up test artifacts
    rm -f "${CONFIG_DIR}/session_history.json"
}

@test "complete pomodoro workflow: work -> break -> work -> long break" {
    # Start first work session with short duration
    run "$APP_BINARY" start 2s
    [ "$status" -eq 0 ]
    
    # Verify work session is running
    run "$APP_BINARY" status
    [ "$status" -eq 0 ]
    [[ "$output" =~ "work" ]]
    [[ "$output" =~ "running" ]]
    
    # Wait for work session to complete naturally
    sleep 3
    
    # Start break session with short duration
    run "$APP_BINARY" break 1s
    [ "$status" -eq 0 ]
    
    # Verify break session is running
    run "$APP_BINARY" status
    [ "$status" -eq 0 ]
    [[ "$output" =~ "break" ]]
    [[ "$output" =~ "running" ]]
    
    # Wait for break session to complete naturally
    sleep 2
    
    # Start second work session with short duration
    run "$APP_BINARY" start 2s
    [ "$status" -eq 0 ]
    
    # Wait for second work session to complete naturally
    sleep 3
    
    # Start long break session with short duration
    run "$APP_BINARY" long-break 1s
    [ "$status" -eq 0 ]
    
    # Verify long break session is running
    run "$APP_BINARY" status
    [ "$status" -eq 0 ]
    [[ "$output" =~ "long break" ]]
    [[ "$output" =~ "running" ]]
    
    # Wait for long break session to complete naturally
    sleep 2
    
    # Verify all sessions are recorded in history
    run "$APP_BINARY" history
    [ "$status" -eq 0 ]
    [[ "$output" =~ "work" ]]
    [[ "$output" =~ "break" ]]
    [[ "$output" =~ "long break" ]]
}


@test "configuration management workflow" {
    # Show initial configuration
    run "$APP_BINARY" config show
    [ "$status" -eq 0 ]
    [[ "$output" =~ "1m" ]]
    [[ "$output" =~ "30s" ]]
    [[ "$output" =~ "2m" ]]
    
    # Update work duration
    run "$APP_BINARY" config set default_work_duration 25m
    [ "$status" -eq 0 ]
    
    # Update break duration
    run "$APP_BINARY" config set default_break_duration 5m
    [ "$status" -eq 0 ]
    
    # Update long break duration
    run "$APP_BINARY" config set default_long_break_duration 15m
    [ "$status" -eq 0 ]
    
    # Verify configuration changes
    run "$APP_BINARY" config show
    [ "$status" -eq 0 ]
    [[ "$output" =~ "25m" ]]
    [[ "$output" =~ "5m" ]]
    [[ "$output" =~ "15m" ]]
    
    # Test that new configuration is used with very short duration
    run "$APP_BINARY" start 1s
    [ "$status" -eq 0 ]
    
    run "$APP_BINARY" status
    [ "$status" -eq 0 ]
    [[ "$output" =~ "running" ]]
    
    # Wait for timer to complete naturally
    sleep 2
}

@test "force flag workflow" {
    # Start a work session
    run "$APP_BINARY" start 1m
    [ "$status" -eq 0 ]
    
    # Try to start another session without force (should fail)
    run "$APP_BINARY" start 1m
    [ "$status" -eq 1 ]
    [[ "$output" =~ "error" ]]
    
    # Start another session with force (should succeed)
    run "$APP_BINARY" start --force 1m
    [ "$status" -eq 0 ]
    
    # Verify the new session is running
    run "$APP_BINARY" status
    [ "$status" -eq 0 ]
    [[ "$output" =~ "running" ]]
    
    # Wait for the session to complete naturally
    sleep 2
}

@test "error handling workflow" {
    # Test invalid commands
    run "$APP_BINARY" invalid-command
    [ "$status" -eq 1 ]
    
    # Test invalid durations
    run "$APP_BINARY" start invalid
    [ "$status" -eq 1 ]
    [[ "$output" =~ "error" ]]
    
    # Test invalid command
    run "$APP_BINARY" invalid-operation
    [ "$status" -eq 1 ]
    [[ "$output" =~ "error" ]]
    
    # Test invalid configuration
    run "$APP_BINARY" config set default_work_duration invalid
    [ "$status" -eq 1 ]
    [[ "$output" =~ "error" ]]
}

@test "real-time timer completion workflow" {
    # Start a very short session
    run "$APP_BINARY" start 1s
    [ "$status" -eq 0 ]
    
    # Wait for auto-completion
    sleep 3
    
    # Check status (should be idle)
    run "$APP_BINARY" status
    [ "$status" -eq 0 ]
    [[ "$output" =~ "idle" ]]
    
    # Verify session was recorded in history
    run "$APP_BINARY" history
    [ "$status" -eq 0 ]
    [[ "$output" =~ "work" ]]
    [[ "$output" =~ "completed" ]]
}

@test "session history workflow" {
    # Create multiple sessions with very short durations
    "$APP_BINARY" start 1s > /dev/null 2>&1
    sleep 2  # Wait for completion
    
    "$APP_BINARY" break 1s > /dev/null 2>&1
    sleep 2  # Wait for completion
    
    "$APP_BINARY" long-break 1s > /dev/null 2>&1
    sleep 2  # Wait for completion
    
    # Check history shows all sessions
    run "$APP_BINARY" history
    [ "$status" -eq 0 ]
    
    # Verify all session types are present
    work_count=$(echo "$output" | grep -c "work" || echo "0")
    break_count=$(echo "$output" | grep -c "break" || echo "0")
    long_break_count=$(echo "$output" | grep -c "long break" || echo "0")
    
    [ "$work_count" -ge 1 ]
    [ "$break_count" -ge 1 ]
    [ "$long_break_count" -ge 1 ]
}

@test "completion commands workflow" {
    # Test bash completion
    run "$APP_BINARY" completion bash
    [ "$status" -eq 0 ]
    [[ "$output" =~ "complete" ]]
    
    # Test zsh completion
    run "$APP_BINARY" completion zsh
    [ "$status" -eq 0 ]
    [[ "$output" =~ "compdef" ]]
}

@test "version and help workflow" {
    # Test version command
    run "$APP_BINARY" version
    [ "$status" -eq 0 ]
    [[ "$output" =~ "pomodux" ]]
    
    # Test help command
    run "$APP_BINARY" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "start" ]]
    [[ "$output" =~ "status" ]]
    [[ "$output" =~ "break" ]]
    [[ "$output" =~ "long-break" ]]
    [[ "$output" =~ "config" ]]
    [[ "$output" =~ "history" ]]
    [[ "$output" =~ "completion" ]]
    [[ "$output" =~ "plugin" ]]
}

@test "state persistence workflow" {
    # Start a session with short duration
    run "$APP_BINARY" start 3s
    [ "$status" -eq 0 ]
    
    # Verify session is running
    run "$APP_BINARY" status
    [ "$status" -eq 0 ]
    [[ "$output" =~ "running" ]]
    
    # Simulate application restart by checking status again
    run "$APP_BINARY" status
    [ "$status" -eq 0 ]
    [[ "$output" =~ "running" ]]
    
    # Wait for the session to complete naturally
    sleep 4
    
    # Verify session is recorded
    run "$APP_BINARY" history
    [ "$status" -eq 0 ]
    [[ "$output" =~ "work" ]]
}

 