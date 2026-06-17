package bucket

import (
	"github.com/spf13/cobra"
)

// BucketCmd is the parent command for all bucket-related subcommands:
//
//	cloudctl bucket list
//	cloudctl bucket create
//	cloudctl bucket delete
//	cloudctl bucket list-objects
var BucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Manage cloud storage buckets (S3 / GCS)",
}
