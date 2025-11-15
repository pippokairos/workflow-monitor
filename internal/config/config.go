// internal/config/config.go
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// Atlassian
	AtlassianURL         string   `yaml:"atlassian_url"`
	AtlassianEmail       string   `yaml:"atlassian_email"`
	AtlassianToken       string   `yaml:"atlassian_token"`
	AtlassianProjectKeys []string `yaml:"atlassian_project_keys"`

	// GitHub
	GitHubToken    string   `yaml:"github_token"`
	GitHubRepos    []string `yaml:"github_repos"`
	GitHubUsername string   `yaml:"github_username"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	expanded := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
