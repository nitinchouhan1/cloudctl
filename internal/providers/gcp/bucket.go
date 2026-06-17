package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/nitinchouhan1/cloudctl/internal/model"
	"github.com/nitinchouhan1/cloudctl/internal/utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// loadGCSClient builds a Cloud Storage client from the refresh token saved
// by gcp.Login(), the same way Login() builds its oauth2 token.
func loadGCSClient(ctx context.Context) (*storage.Client, error) {

	cfgFile, err := utils.LoadConfig()
	if err != nil {
		return nil, err
	}

	provider, ok := cfgFile.Providers["gcp"]
	if !ok {
		return nil, fmt.Errorf("not logged into GCP: run 'cloudctl auth login gcp' first")
	}

	oauthConfig := &oauth2.Config{
		ClientID:     model.GCP_CLIENT_ID,
		ClientSecret: model.GCP_CLIENT_SECRET,
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	tokenSource := oauthConfig.TokenSource(
		ctx,
		&oauth2.Token{RefreshToken: provider.RefreshToken},
	)

	client, err := storage.NewClient(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return client, nil
}

// resolveProjectID returns the explicit projectID if given, otherwise
// returns an error asking the user to supply one. GCS requires a project ID
// for bucket-level operations like listing and creating buckets (not for
// reading/deleting individual buckets by name).
func resolveProjectID(projectID string) (string, error) {
	if projectID != "" {
		return projectID, nil
	}
	return "", fmt.Errorf("GCP project ID is required: pass --project")
}

// ListBuckets lists every GCS bucket in the given project.
func ListBuckets(projectID string) error {

	projectID, err := resolveProjectID(projectID)
	if err != nil {
		return err
	}

	ctx := context.Background()

	client, err := loadGCSClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	it := client.Buckets(ctx, projectID)

	fmt.Printf("%-40s %s\n", "NAME", "LOCATION")

	found := false
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list buckets: %w", err)
		}
		found = true
		fmt.Printf("%-40s %s\n", attrs.Name, attrs.Location)
	}

	if !found {
		fmt.Println("No GCS buckets found.")
	}

	return nil
}

// CreateBucket creates a new GCS bucket with the given name in the given
// project. location follows GCS conventions, e.g. "US", "asia-south1".
func CreateBucket(projectID string, name string, location string) error {

	if name == "" {
		return fmt.Errorf("bucket name is required")
	}

	projectID, err := resolveProjectID(projectID)
	if err != nil {
		return err
	}

	ctx := context.Background()

	client, err := loadGCSClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(name)

	attrs := &storage.BucketAttrs{}
	if location != "" {
		attrs.Location = location
	}

	if err := bucket.Create(ctx, projectID, attrs); err != nil {
		return fmt.Errorf("failed to create bucket %q: %w", name, err)
	}

	fmt.Printf("✓ Created GCS bucket %q\n", name)

	return nil
}

// DeleteBucket deletes a GCS bucket. The bucket must already be empty;
// pass force=true to empty it of objects first.
func DeleteBucket(name string, force bool) error {

	if name == "" {
		return fmt.Errorf("bucket name is required")
	}

	ctx := context.Background()

	client, err := loadGCSClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(name)

	if force {
		if err := emptyBucket(ctx, bucket); err != nil {
			return fmt.Errorf("failed to empty bucket %q before delete: %w", name, err)
		}
	}

	if err := bucket.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete bucket %q: %w", name, err)
	}

	fmt.Printf("✓ Deleted GCS bucket %q\n", name)

	return nil
}

// emptyBucket deletes every object in the bucket so it can be removed by
// DeleteBucket.
func emptyBucket(ctx context.Context, bucket *storage.BucketHandle) error {

	it := bucket.Objects(ctx, nil)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		if err := bucket.Object(attrs.Name).Delete(ctx); err != nil {
			return err
		}
	}

	return nil
}

// ListObjects lists objects in a bucket, optionally filtered by prefix.
func ListObjects(bucketName string, prefix string) error {

	if bucketName == "" {
		return fmt.Errorf("bucket name is required")
	}

	ctx := context.Background()

	client, err := loadGCSClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	query := &storage.Query{}
	if prefix != "" {
		query.Prefix = prefix
	}

	it := client.Bucket(bucketName).Objects(ctx, query)

	fmt.Printf("%-60s %12s %s\n", "NAME", "SIZE", "UPDATED")

	count := 0
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list objects in %q: %w", bucketName, err)
		}

		fmt.Printf(
			"%-60s %12d %s\n",
			attrs.Name,
			attrs.Size,
			attrs.Updated.Local().Format("2006-01-02 15:04:05"),
		)
		count++
	}

	if count == 0 {
		fmt.Println("No objects found.")
	}

	return nil
}
