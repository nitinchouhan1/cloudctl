package auth

import (
	"fmt"

	"github.com/nitinchouhan1/cloudctl/internal/utils"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage cloudctl configuration",
	Long:  `Manage cloudctl configuration settings like default GCP project.`,
}

var setProjectCmd = &cobra.Command{
	Use:   "set-project [project-id]",
	Short: "Set the default GCP project",
	Long: `Set the default GCP project ID to use for GCP operations.

This project will be used by default when running GCP commands 
(like bucket operations) if the --project flag is not specified.

Example:
  cloudctl auth config set-project flowvoice
  cloudctl auth config set-project my-gcp-project-123`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		projectID := args[0]

		// Load existing config
		cfg, err := utils.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Set the default GCP project
		cfg.DefaultGCPProject = projectID

		// Save the updated config
		if err := utils.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✓ Default GCP project set to: %s\n", projectID)
		return nil
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current cloudctl configuration settings.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := utils.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		fmt.Println("Current Configuration:")
		fmt.Printf("  Provider: %s\n", cfg.CurrentProvider)

		if cfg.DefaultGCPProject != "" {
			fmt.Printf("  Default GCP Project: %s\n", cfg.DefaultGCPProject)
		} else {
			fmt.Println("  Default GCP Project: (not set)")
		}

		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(setProjectCmd)
	ConfigCmd.AddCommand(showConfigCmd)
	AuthCmd.AddCommand(ConfigCmd)
}
