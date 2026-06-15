package cmd

import (
	"os"

	"github.com/nitinchouhan1/cloudctl/cmd/auth"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use: "cloudctl",
}

func init() {
	rootCmd.AddCommand(auth.AuthCmd)
}
