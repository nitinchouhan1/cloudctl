package auth

import "github.com/spf13/cobra"

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long: `Manage authentication for cloud providers.

Examples:

  cloudctl auth login aws
  cloudctl auth login gcp

  cloudctl auth current

  cloudctl auth switch aws
`,
}

func init() {
	AuthCmd.AddCommand(LoginCmd)
}
