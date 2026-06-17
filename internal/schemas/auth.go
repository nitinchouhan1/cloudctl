package schemas

type Config struct {
	CurrentProvider   string              `yaml:"current_provider"`
	DefaultGCPProject string              `yaml:"default_gcp_project,omitempty"`
	Providers         map[string]Provider `yaml:"providers"`
}

type Provider struct {
	AccessKeyID     string `yaml:"access_key_id,omitempty"`
	SecretAccessKey string `yaml:"secret_access_key,omitempty"`
	AccountID       string `yaml:"account_id,omitempty"`
	ARN             string `yaml:"arn,omitempty"`

	ClientID     string `yaml:"client_id,omitempty"`
	RefreshToken string `yaml:"refresh_token,omitempty"`
	ProjectID    string `yaml:"project_id,omitempty"`

	Region string `yaml:"region,omitempty"`
}
