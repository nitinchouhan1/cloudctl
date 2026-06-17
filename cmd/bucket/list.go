package bucket

import (
	"fmt"

	awsprovider "github.com/nitinchouhan1/cloudctl/internal/providers/aws"
	gcpprovider "github.com/nitinchouhan1/cloudctl/internal/providers/gcp"

	"github.com/spf13/cobra"
)

var (
	listProviderFlag string
	listProjectFlag  string
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List buckets",

	RunE: func(cmd *cobra.Command, args []string) error {

		provider, err := resolveProvider(listProviderFlag)
		if err != nil {
			return err
		}

		switch provider {
		case "aws":
			return awsprovider.ListBuckets()

		case "gcp":
			project, err := resolveGCPProject(listProjectFlag)
			if err != nil {
				return err
			}
			return gcpprovider.ListBuckets(project)

		default:
			return fmt.Errorf("unsupported provider: %s", provider)
		}
	},
}

func init() {
	ListCmd.Flags().StringVar(&listProviderFlag, "provider", "", "Cloud provider (aws or gcp); defaults to the currently logged-in provider")
	ListCmd.Flags().StringVar(&listProjectFlag, "project", "", "GCP project ID; defaults to configured default project")

	BucketCmd.AddCommand(ListCmd)
}
