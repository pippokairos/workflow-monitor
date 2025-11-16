// internal/config/config.go
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// Atlassian
	AtlassianURL          string   `yaml:"atlassian_url"`
	AtlassianEmail        string   `yaml:"atlassian_email"`
	AtlassianToken        string   `yaml:"atlassian_token"`
	AtlassianStatusReview string   `yaml:"atlassian_status_review"`
	AtlassianStatusDone   string   `yaml:"atlassian_status_done"`
	AtlassianProjectKeys  []string `yaml:"atlassian_project_keys"`

	// GitHub
	GitHubUsername          string   `yaml:"github_username"`
	GitHubToken             string   `yaml:"github_token"`
	GitHubRequiredApprovers int      `yaml:"github_required_approvers"`
	GitHubRepos             []string `yaml:"github_repos"`

	// Matching
	IssuePattern string `yaml:"issue_pattern"`
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

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg.AtlassianURL == "" {
		return fmt.Errorf("atlassian_url is required")
	}

	if cfg.AtlassianEmail == "" {
		return fmt.Errorf("atlassian_email is required")
	}

	if cfg.AtlassianToken == "" {
		return fmt.Errorf("atlassian_token is required")
	}

	if cfg.GitHubUsername == "" {
		return fmt.Errorf("github_username is required")
	}

	if cfg.GitHubToken == "" {
		return fmt.Errorf("github_token is required")
	}

	if cfg.IssuePattern == "" {
		return fmt.Errorf("issue_pattern is required")
	}

	if len(cfg.GitHubRepos) == 0 {
		return fmt.Errorf("at least one github_repo is required")
	}

	return nil
}
