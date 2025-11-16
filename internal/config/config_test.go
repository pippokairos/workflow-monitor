package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	os.Setenv("TEST_TOKEN", "secret-token")
	defer os.Unsetenv("TEST_TOKEN")

	content := `
atlassian_url: https://test.atlassian.net
atlassian_email: test@example.com
atlassian_token: ${TEST_TOKEN}
atlassian_project_keys:
  - PROJ
  - TEST

github_token: gh-token
github_username: testuser
github_repos:
  - owner/repo1
  - owner/repo2

issue_pattern: '([A-Z]+-\d+)'
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify config values
	if cfg.AtlassianURL != "https://test.atlassian.net" {
		t.Errorf("Expected URL 'https://test.atlassian.net', got '%s'", cfg.AtlassianURL)
	}

	if cfg.AtlassianToken != "secret-token" {
		t.Errorf("Expected token 'secret-token', got '%s'", cfg.AtlassianToken)
	}

	if len(cfg.AtlassianProjectKeys) != 2 {
		t.Errorf("Expected 2 project keys, got %d", len(cfg.AtlassianProjectKeys))
	}

	if cfg.GitHubUsername != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", cfg.GitHubUsername)
	}

	if len(cfg.GitHubRepos) != 2 {
		t.Errorf("Expected 2 repos, got %d", len(cfg.GitHubRepos))
	}
}

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		cfg     Config
		wantErr bool
		errMsg  string
	}{
		"valid config": {
			cfg: Config{
				AtlassianURL:         "https://test.atlassian.net",
				AtlassianEmail:       "test@example.com",
				AtlassianToken:       "token",
				AtlassianProjectKeys: []string{"PROJ"},
				GitHubToken:          "gh-token",
				GitHubUsername:       "user",
				GitHubRepos:          []string{"owner/repo"},
				IssuePattern:         `([A-Z]+-\d+)`,
			},
			wantErr: false,
		},
		"missing atlassian url": {
			cfg: Config{
				AtlassianEmail:       "test@example.com",
				AtlassianToken:       "token",
				AtlassianProjectKeys: []string{"PROJ"},
				GitHubToken:          "gh-token",
				GitHubUsername:       "user",
				GitHubRepos:          []string{"owner/repo"},
				IssuePattern:         `([A-Z]+-\d+)`,
			},
			wantErr: true,
			errMsg:  "atlassian_url is required",
		},
		"missing atlassian email": {
			cfg: Config{
				AtlassianURL:         "https://test.atlassian.net",
				AtlassianToken:       "token",
				AtlassianProjectKeys: []string{"PROJ"},
				GitHubToken:          "gh-token",
				GitHubUsername:       "user",
				GitHubRepos:          []string{"owner/repo"},
				IssuePattern:         `([A-Z]+-\d+)`,
			},
			wantErr: true,
			errMsg:  "atlassian_email is required",
		},
		"missing atlassian token": {
			cfg: Config{
				AtlassianURL:         "https://test.atlassian.net",
				AtlassianEmail:       "test@example.com",
				AtlassianProjectKeys: []string{"PROJ"},
				GitHubToken:          "gh-token",
				GitHubUsername:       "user",
				GitHubRepos:          []string{"owner/repo"},
				IssuePattern:         `([A-Z]+-\d+)`,
			},
			wantErr: true,
			errMsg:  "atlassian_token is required",
		},
		"missing github username": {
			cfg: Config{
				AtlassianURL:         "https://test.atlassian.net",
				AtlassianEmail:       "test@example.com",
				AtlassianToken:       "token",
				AtlassianProjectKeys: []string{"PROJ"},
				GitHubToken:          "gh-token",
				GitHubRepos:          []string{"owner/repo"},
				IssuePattern:         `([A-Z]+-\d+)`,
			},
			wantErr: true,
			errMsg:  "github_username is required",
		},
		"missing github token": {
			cfg: Config{
				AtlassianURL:         "https://test.atlassian.net",
				AtlassianEmail:       "test@example.com",
				AtlassianToken:       "token",
				AtlassianProjectKeys: []string{"PROJ"},
				GitHubUsername:       "user",
				GitHubRepos:          []string{"owner/repo"},
				IssuePattern:         `([A-Z]+-\d+)`,
			},
			wantErr: true,
			errMsg:  "github_token is required",
		},
		"missing github repos": {
			cfg: Config{
				AtlassianURL:         "https://test.atlassian.net",
				AtlassianEmail:       "test@example.com",
				AtlassianToken:       "token",
				AtlassianProjectKeys: []string{"PROJ"},
				GitHubToken:          "gh-token",
				GitHubUsername:       "user",
				GitHubRepos:          []string{},
				IssuePattern:         `([A-Z]+-\d+)`,
			},
			wantErr: true,
			errMsg:  "at least one github_repo is required",
		},
		"missing issue pattern": {
			cfg: Config{
				AtlassianURL:         "https://test.atlassian.net",
				AtlassianEmail:       "test@example.com",
				AtlassianToken:       "token",
				AtlassianProjectKeys: []string{},
				GitHubToken:          "gh-token",
				GitHubUsername:       "user",
				GitHubRepos:          []string{"owner/repo"},
				IssuePattern:         "",
			},
			wantErr: true,
			errMsg:  "issue_pattern is required",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Expected error '%s', got '%s'", tt.errMsg, err.Error())
			}
		})
	}
}
