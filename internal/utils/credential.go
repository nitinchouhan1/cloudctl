package utils

import (
	"os"
	"path/filepath"

	"github.com/nitinchouhan1/cloudctl/internal/schemas"
	"gopkg.in/yaml.v2"
)

func ConfigPath() string {
	home, _ := os.UserHomeDir()

	return filepath.Join(
		home,
		".cloudctl",
		"config.yaml",
	)
}

func SaveConfig(cfg *schemas.Config) error {

	path := ConfigPath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func LoadConfig() (*schemas.Config, error) {

	path := ConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return &schemas.Config{}, nil
	}

	var cfg schemas.Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
