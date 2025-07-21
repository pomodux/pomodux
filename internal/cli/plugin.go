package cli

import (
	"fmt"
	"os"

	"github.com/pomodux/pomodux/internal/config"
	"github.com/pomodux/pomodux/internal/logger"
	"github.com/pomodux/pomodux/internal/plugin"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage plugins",
	Long:  `Manage Pomodux plugins - enable, disable, or list plugins.`,
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all loaded plugins",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Create plugin manager
		pm := plugin.NewPluginManager(cfg.Plugins.Directory, cfg)
		if err := pm.LoadPlugins(); err != nil {
			fmt.Printf("Error loading plugins: %v\n", err)
			os.Exit(1)
		}

		// List plugins
		plugins := pm.ListPlugins()
		if len(plugins) == 0 {
			fmt.Println("No plugins loaded.")
			return
		}

		fmt.Println("Loaded plugins:")
		fmt.Println("===============")
		for _, p := range plugins {
			status := "Enabled"
			if !p.Enabled {
				status = "Disabled"
			}
			fmt.Printf("• %s (v%s) - %s [%s]\n", p.Name, p.Version, p.Description, status)
		}
	},
}

var pluginDisableCmd = &cobra.Command{
	Use:   "disable [plugin-name]",
	Short: "Disable a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Initialize plugins.Enabled map if it doesn't exist
		if cfg.Plugins.Enabled == nil {
			cfg.Plugins.Enabled = make(map[string]bool)
		}

		// Disable plugin in config
		cfg.Plugins.Enabled[pluginName] = false

		// Save configuration
		if err := config.Save(cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Plugin '%s' disabled successfully.\n", pluginName)
		fmt.Printf("   The plugin will remain disabled until you enable it again.\n")
		logger.Info("Plugin disabled", map[string]interface{}{"plugin": pluginName})
	},
}

var pluginEnableCmd = &cobra.Command{
	Use:   "enable [plugin-name]",
	Short: "Enable a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Initialize plugins.Enabled map if it doesn't exist
		if cfg.Plugins.Enabled == nil {
			cfg.Plugins.Enabled = make(map[string]bool)
		}

		// Enable plugin in config
		cfg.Plugins.Enabled[pluginName] = true

		// Save configuration
		if err := config.Save(cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Plugin '%s' enabled successfully.\n", pluginName)
		fmt.Printf("   The plugin will be loaded on the next timer start.\n")
		logger.Info("Plugin enabled", map[string]interface{}{"plugin": pluginName})
	},
}

func init() {
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginDisableCmd)
	pluginCmd.AddCommand(pluginEnableCmd)
	GetRootCmd().AddCommand(pluginCmd)
}
