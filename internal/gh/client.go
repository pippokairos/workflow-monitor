package gh

import (
	"context"
	"fmt"
	"slices"
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

func (c *Client) FetchOpenPRs() ([]PullRequest, error) {
	var openPRs []PullRequest

	for i := range c.repos {
		owner, repo := getOwnerAndRepo(c.repos[i])
		options := &github.PullRequestListOptions{State: "open"}
		githubPRs, resp, err := c.github.PullRequests.List(context.Background(), owner, repo, options)
		debug.Printf("GitHub PullRequests List response: %+v", resp)
		if err != nil || len(githubPRs) == 0 {
			return nil, err
		}

		openPRs := make([]PullRequest, 0, len(githubPRs))
		for _, githubPR := range githubPRs {
			reviews, resp, err := c.github.PullRequests.ListReviews(context.Background(), owner, repo, githubPR.GetNumber(), nil)
			debug.Printf("GitHub PullRequest ListReviews response: %+v", resp)
			if err != nil {
				return nil, err
			}

			var approvers []string
			for _, review := range reviews {
				if review.GetState() == "APPROVED" {
					reviewer := review.GetUser().GetLogin()
					if !slices.Contains(approvers, reviewer) {
						approvers = append(approvers, reviewer)
					}
				}
			}

			openPRs = append(openPRs, *ToPullRequest(githubPR, approvers))
		}
	}

	return openPRs, nil
}

// This needs to be a separate call, because the PullRequests.List method does not support filtering by review requested.
func (c *Client) FetchPRsNeedingMyReview() ([]PullRequest, error) {
	query := fmt.Sprintf("is:pr is:open review-requested:%s", c.username)
	for i := range c.repos {
		query += fmt.Sprintf(" repo:%s", c.repos[i])
	}

	opts := &github.SearchOptions{}
	result, _, err := c.github.Search.Issues(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}

	var prs []PullRequest
	for _, issue := range result.Issues {
		prs = append(prs, PullRequest{
			Username:  issue.GetUser().GetLogin(),
			Title:     issue.GetTitle(),
			Approvers: []string{}, // They are not available in the search result, but they are not needed in this context.
		})
	}

	return prs, nil
}

func getOwnerAndRepo(repoCfg string) (string, string) {
	parts := strings.Split(repoCfg, "/")
	return parts[0], parts[1]
}
