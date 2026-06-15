package aws

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/nitinchouhan1/cloudctl/internal/schemas"
	"github.com/nitinchouhan1/cloudctl/internal/utils"
)

func Login() error {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("AWS Access Key ID: ")
	accessKey, _ := reader.ReadString('\n')

	fmt.Print("AWS Secret Access Key: ")
	secretKey, _ := reader.ReadString('\n')

	fmt.Print("AWS Region: ")
	region, _ := reader.ReadString('\n')

	accessKey = strings.TrimSpace(accessKey)
	secretKey = strings.TrimSpace(secretKey)
	region = strings.TrimSpace(region)

	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     accessKey,
					SecretAccessKey: secretKey,
				}, nil
			}),
		),
	)
	if err != nil {
		return err
	}

	stsClient := sts.NewFromConfig(cfg)

	identity, err := stsClient.GetCallerIdentity(
		context.Background(),
		&sts.GetCallerIdentityInput{},
	)
	if err != nil {
		return fmt.Errorf("invalid AWS credentials: %w", err)
	}

	cfgFile, err := utils.LoadConfig()
	if err != nil {
		return err
	}

	if cfgFile.Providers == nil {
		cfgFile.Providers = make(map[string]schemas.Provider)
	}

	cfgFile.Providers["aws"] = schemas.Provider{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		Region:          region,
		AccountID:       aws.ToString(identity.Account),
		ARN:             aws.ToString(identity.Arn),
	}

	cfgFile.CurrentProvider = "aws"

	if err := utils.SaveConfig(cfgFile); err != nil {
		return err
	}

	fmt.Printf(
		"✓ Logged into AWS Account %s\n",
		aws.ToString(identity.Account),
	)

	return nil
}
