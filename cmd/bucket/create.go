package bucket

import (
	"fmt"

	awsprovider "github.com/nitinchouhan1/cloudctl/internal/providers/aws"
	gcpprovider "github.com/nitinchouhan1/cloudctl/internal/providers/gcp"

	"github.com/spf13/cobra"
)

var (
	createProviderFlag string
	createProjectFlag  string
	createRegionFlag   string
)

var CreateCmd = &cobra.Command{
	Use:   "create [bucket-name]",
	Short: "Create a bucket",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {

		name := args[0]

		provider, err := resolveProvider(createProviderFlag)
		if err != nil {
			return err
		}

		switch provider {
		case "aws":
			return awsprovider.CreateBucket(name, createRegionFlag)

		case "gcp":
			project, err := resolveGCPProject(createProjectFlag)
			if err != nil {
				return err
			}
			return gcpprovider.CreateBucket(project, name, createRegionFlag)

		default:
			return fmt.Errorf("unsupported provider: %s", provider)
		}
	},
}

func init() {
	CreateCmd.Flags().StringVar(&createProviderFlag, "provider", "", "Cloud provider (aws or gcp); defaults to the currently logged-in provider")
	CreateCmd.Flags().StringVar(&createProjectFlag, "project", "", "GCP project ID; defaults to configured default project")
	CreateCmd.Flags().StringVar(&createRegionFlag, "region", "", "Bucket region/location; defaults to your saved AWS region for AWS, required for GCP")

	BucketCmd.AddCommand(CreateCmd)
}
