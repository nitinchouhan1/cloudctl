package bucket

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	awsprovider "github.com/nitinchouhan1/cloudctl/internal/providers/aws"
	gcpprovider "github.com/nitinchouhan1/cloudctl/internal/providers/gcp"

	"github.com/spf13/cobra"
)

var (
	deleteProviderFlag string
	deleteForceFlag    bool
	deleteYesFlag      bool
	deleteRegionFlag   string
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [bucket-name]",
	Short: "Delete a bucket",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {

		name := args[0]

		provider, err := resolveProvider(deleteProviderFlag)
		if err != nil {
			return err
		}

		if !deleteYesFlag {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Are you sure you want to delete bucket %q? [y/N]: ", name)
			answer, _ := reader.ReadString('\n')
			answer = strings.ToLower(strings.TrimSpace(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		fmt.Printf("Deleting bucket %q in region %q...\n", name, deleteRegionFlag)
		switch provider {
		case "aws":
			return awsprovider.DeleteBucket(name, deleteRegionFlag, deleteForceFlag)

		case "gcp":
			return gcpprovider.DeleteBucket(name, deleteForceFlag)

		default:
			return fmt.Errorf("unsupported provider: %s", provider)
		}
	},
}

func init() {
	DeleteCmd.Flags().StringVar(&deleteProviderFlag, "provider", "", "Cloud provider (aws or gcp); defaults to the currently logged-in provider")
	DeleteCmd.Flags().BoolVar(&deleteForceFlag, "force", false, "Delete all objects in the bucket first")
	DeleteCmd.Flags().StringVar(&deleteRegionFlag, "region", "", "Bucket region/location; defaults to your saved AWS region for AWS, required for GCP")
	DeleteCmd.Flags().BoolVarP(&deleteYesFlag, "yes", "y", false, "Skip the confirmation prompt")

	BucketCmd.AddCommand(DeleteCmd)
}
