package bucket

import (
	"fmt"

	awsprovider "github.com/nitinchouhan1/cloudctl/internal/providers/aws"
	gcpprovider "github.com/nitinchouhan1/cloudctl/internal/providers/gcp"

	"github.com/spf13/cobra"
)

var (
	listObjectsProviderFlag string
	listObjectsPrefixFlag   string
)

var ListObjectsCmd = &cobra.Command{
	Use:   "list-objects [bucket-name]",
	Short: "List objects inside a bucket",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {

		bucketName := args[0]

		provider, err := resolveProvider(listObjectsProviderFlag)
		if err != nil {
			return err
		}

		switch provider {
		case "aws":
			return awsprovider.ListObjects(bucketName, listObjectsPrefixFlag)

		case "gcp":
			return gcpprovider.ListObjects(bucketName, listObjectsPrefixFlag)

		default:
			return fmt.Errorf("unsupported provider: %s", provider)
		}
	},
}

func init() {
	ListObjectsCmd.Flags().StringVar(&listObjectsProviderFlag, "provider", "", "Cloud provider (aws or gcp); defaults to the currently logged-in provider")
	ListObjectsCmd.Flags().StringVar(&listObjectsPrefixFlag, "prefix", "", "Only list objects whose key starts with this prefix")

	BucketCmd.AddCommand(ListObjectsCmd)
}
