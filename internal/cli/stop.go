package cli

import (
	"fmt"

	"github.com/pomodux/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the currently running timer",
	Long: `Stop the currently running timer.
	
This command will stop the timer and record the session as interrupted.
The timer must be running (not paused or completed) for this command to work.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the global timer
		t := timer.GetGlobalTimer()

		// Check if timer is running
		if t.GetStatus() != timer.StatusRunning {
			return fmt.Errorf("timer is not running (current status: %s)", t.GetStatus())
		}

		// Stop the timer
		if err := t.Stop(); err != nil {
			return fmt.Errorf("failed to stop timer: %v", err)
		}

		fmt.Println("⏹️  Timer stopped.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
