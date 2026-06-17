package bucket

import (
	"fmt"

	"github.com/nitinchouhan1/cloudctl/internal/utils"
)

// resolveProvider returns the provider to act on: the explicit --provider
// flag value if given, otherwise the currently logged-in provider saved by
// `cloudctl auth login`.
func resolveProvider(flagValue string) (string, error) {

	if flagValue != "" {
		switch flagValue {
		case "aws", "gcp":
			return flagValue, nil
		default:
			return "", fmt.Errorf("unsupported provider: %s (expected aws or gcp)", flagValue)
		}
	}

	cfgFile, err := utils.LoadConfig()
	if err != nil {
		return "", err
	}

	if cfgFile.CurrentProvider == "" {
		return "", fmt.Errorf("no provider specified and no current provider set; pass --provider or run 'cloudctl auth login'")
	}

	return cfgFile.CurrentProvider, nil
}

// resolveGCPProject returns the GCP project to use: the explicit --project
// flag value if given, otherwise the default GCP project from config.
func resolveGCPProject(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}

	cfgFile, err := utils.LoadConfig()
	if err != nil {
		return "", err
	}

	if cfgFile.DefaultGCPProject == "" {
		return "", fmt.Errorf("no GCP project specified; pass --project or set default with 'cloudctl auth config set-project'")
	}

	return cfgFile.DefaultGCPProject, nil
}
