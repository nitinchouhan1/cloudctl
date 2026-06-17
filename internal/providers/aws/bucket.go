package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/nitinchouhan1/cloudctl/internal/utils"
)

// loadS3Client builds an S3 client from the credentials already saved by
// `cloudctl auth login aws`, the same way aws.Login() builds its STS client.
func loadS3Client(region string) (*s3.Client, string, error) {

	cfgFile, err := utils.LoadConfig()
	if err != nil {
		return nil, "", err
	}

	provider, ok := cfgFile.Providers["aws"]
	if !ok {
		return nil, "", fmt.Errorf("not logged into AWS: run 'cloudctl auth login aws' first")
	}
	selectedRegion := region
	if selectedRegion == "" {
		selectedRegion = provider.Region
	}

	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(selectedRegion),
		awsconfig.WithCredentialsProvider(
			aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     provider.AccessKeyID,
					SecretAccessKey: provider.SecretAccessKey,
				}, nil
			}),
		),
	)
	if err != nil {
		return nil, "", err
	}

	return s3.NewFromConfig(cfg), provider.Region, nil
}

// ListBuckets lists every S3 bucket owned by the logged-in AWS account.
func ListBuckets() error {

	client, _, err := loadS3Client("")
	if err != nil {
		return err
	}

	out, err := client.ListBuckets(
		context.Background(),
		&s3.ListBucketsInput{},
	)
	if err != nil {
		return fmt.Errorf("failed to list buckets: %w", err)
	}

	if len(out.Buckets) == 0 {
		fmt.Println("No S3 buckets found.")
		return nil
	}

	fmt.Printf("%-40s %s\n", "NAME", "CREATED")
	for _, b := range out.Buckets {
		created := ""
		if b.CreationDate != nil {
			created = b.CreationDate.Local().Format("2006-01-02 15:04:05")
		}
		fmt.Printf("%-40s %s\n", aws.ToString(b.Name), created)
	}

	return nil
}

// CreateBucket creates a new S3 bucket with the given name.
// If region is empty, the region from the saved AWS provider config is used.
func CreateBucket(name string, region string) error {

	if name == "" {
		return fmt.Errorf("bucket name is required")
	}

	client, providerRegion, err := loadS3Client(region)
	if err != nil {
		return err
	}

	if region == "" {
		region = providerRegion
	}
	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	// us-east-1 is AWS's default region and must NOT be passed as a
	// LocationConstraint, or the API returns InvalidLocationConstraint.
	if region != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		}
	}

	_, err = client.CreateBucket(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to create bucket %q: %w", name, err)
	}

	fmt.Printf("✓ Created S3 bucket %q\n", name)

	return nil
}

// DeleteBucket deletes an S3 bucket. The bucket must already be empty;
// pass force=true to empty it of objects first.
func DeleteBucket(name string, region string, force bool) error {

	if name == "" {
		return fmt.Errorf("bucket name is required")
	}

	client, _, err := loadS3Client(region)
	if err != nil {
		return err
	}

	if force {
		if err := emptyBucket(client, name); err != nil {
			return fmt.Errorf("failed to empty bucket %q before delete: %w", name, err)
		}
	}

	_, err = client.DeleteBucket(
		context.Background(),
		&s3.DeleteBucketInput{Bucket: aws.String(name)},
	)
	if err != nil {
		return fmt.Errorf("failed to delete bucket %q: %w", name, err)
	}

	fmt.Printf("✓ Deleted S3 bucket %q\n", name)

	return nil
}

// emptyBucket deletes every object (and version, if versioning is enabled)
// in the bucket so it can be removed by DeleteBucket.
func emptyBucket(client *s3.Client, name string) error {

	paginator := s3.NewListObjectsV2Paginator(
		client,
		&s3.ListObjectsV2Input{Bucket: aws.String(name)},
	)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return err
		}

		for _, obj := range page.Contents {
			_, err := client.DeleteObject(
				context.Background(),
				&s3.DeleteObjectInput{
					Bucket: aws.String(name),
					Key:    obj.Key,
				},
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ListObjects lists objects in a bucket, optionally filtered by prefix.
func ListObjects(bucket string, prefix string) error {

	if bucket == "" {
		return fmt.Errorf("bucket name is required")
	}

	client, _, err := loadS3Client("")
	if err != nil {
		return err
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	paginator := s3.NewListObjectsV2Paginator(client, input)

	count := 0
	fmt.Printf("%-60s %12s %s\n", "KEY", "SIZE", "LAST MODIFIED")

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return fmt.Errorf("failed to list objects in %q: %w", bucket, err)
		}

		for _, obj := range page.Contents {
			modified := ""
			if obj.LastModified != nil {
				modified = obj.LastModified.Local().Format("2006-01-02 15:04:05")
			}
			fmt.Printf(
				"%-60s %12d %s\n",
				aws.ToString(obj.Key),
				aws.ToInt64(obj.Size),
				modified,
			)
			count++
		}
	}

	if count == 0 {
		fmt.Println("No objects found.")
	}

	return nil
}
