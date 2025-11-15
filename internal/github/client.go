package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v79/github"
	"github.com/pippokairos/workflow-monitor/internal/config"
	"github.com/pippokairos/workflow-monitor/internal/debug"
)

type Client struct {
	github   *github.Client
	username string
	repos    []string
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		github:   github.NewClient(nil).WithAuthToken(cfg.GitHubToken),
		username: cfg.GitHubUsername,
		repos:    cfg.GitHubRepos,
	}
}

func (c *Client) FetchOpenPRs() ([]*github.PullRequest, []*github.PullRequest, error) {
	// TODO: handle multiple repos
	parts := strings.Split(c.repos[0], ":")
	owner := parts[0]
	repo := parts[1]

	options := &github.PullRequestListOptions{State: "open"}

	pullRequests, resp, err := c.github.PullRequests.List(context.Background(), owner, repo, options)
	debug.Printf("GitHub response: %+v", resp)
	if err != nil || len(pullRequests) == 0 {
		return nil, nil, err
	}

	myOpenPRs := make([]*github.PullRequest, 0, len(pullRequests))
	otherOpenPRs := make([]*github.PullRequest, 0, len(pullRequests))
	for _, pr := range pullRequests {
		if pr.User.Login != nil && *pr.User.Login == c.username {
			myOpenPRs = append(myOpenPRs, pr)
		} else {
			otherOpenPRs = append(otherOpenPRs, pr)
		}
	}

	return myOpenPRs, otherOpenPRs, err
}
