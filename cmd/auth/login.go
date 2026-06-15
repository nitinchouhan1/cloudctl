package auth

import (
	"fmt"

	"github.com/nitinchouhan1/cloudctl/internal/providers/aws"
	"github.com/nitinchouhan1/cloudctl/internal/providers/gcp"

	"github.com/spf13/cobra"
)

var LoginCmd = &cobra.Command{
	Use:   "login [provider]",
	Short: "Login to a cloud provider",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		provider := args[0]

		switch provider {
		case "aws":
			return aws.Login()

		case "gcp":
			return gcp.Login()

		default:
			return fmt.Errorf(
				"unsupported provider: %s",
				provider,
			)
		}
	},
}

func init() {
	AuthCmd.AddCommand(LoginCmd)
}
