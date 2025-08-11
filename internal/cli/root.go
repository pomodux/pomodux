package cli

import (
	"github.com/spf13/cobra"
	"github.com/pomodux/pomodux/internal/config"
)

var (
	// Used for flags.
	cfgFile string
	
	// Global config instance
	globalConfig *config.Config

	rootCmd = &cobra.Command{
		Use:   "pomodux",
		Short: "A terminal-based timer application with Pomodoro support",
		Long: `Pomodux is a powerful terminal timer application that helps you 
manage your time effectively with work sessions and breaks, including Pomodoro technique support.

Features:
  • Start work timers with custom durations
  • Pomodoro technique support (via 'start' and 'break' commands)
  • Break timer management
  • Session tracking and statistics
  • Rich terminal user interface
  • Plugin system for extensibility`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

// GetRootCmd returns the root command for external access
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pomodux/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// This function is called by cobra when the application starts
	// The actual config loading is handled in main.go
}
