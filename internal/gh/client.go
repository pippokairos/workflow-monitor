package gh

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

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

func (c *Client) FetchOpenPRs(ctx context.Context) ([]PullRequest, error) {
	var allOpenPRs []PullRequest
	var mu sync.Mutex

	for i := range c.repos {
		owner, repo, err := getOwnerAndRepo(c.repos[i])
		if err != nil {
			debug.Printf("Error parsing repo %s: %v", c.repos[i], err)
			continue
		}

		options := &github.PullRequestListOptions{State: "open"}
		githubPRs, resp, err := c.github.PullRequests.List(ctx, owner, repo, options)
		debug.Printf("GitHub PullRequests List response: %+v", resp)
		if err != nil {
			return nil, err
		}

		if len(githubPRs) == 0 {
			continue
		}

		var wg sync.WaitGroup

		for _, githubPR := range githubPRs {
			wg.Add(1)
			go func(githubPR *github.PullRequest) {
				defer wg.Done()

				approvers, err := c.FetchApprovers(ctx, owner, repo, githubPR)
				if err != nil {
					debug.Printf("Error fetching approvers: %v", err)
					return
				}

				mu.Lock()
				allOpenPRs = append(allOpenPRs, *ToInternalPullRequest(githubPR, approvers))
				mu.Unlock()
			}(githubPR)
		}

		wg.Wait()
	}

	return allOpenPRs, nil
}

func (c *Client) FetchApprovers(ctx context.Context, owner, repo string, githubPR *github.PullRequest) ([]string, error) {
	reviews, resp, err := c.github.PullRequests.ListReviews(ctx, owner, repo, githubPR.GetNumber(), nil)
	debug.Printf("GitHub PullRequest ListReviews response: %+v", resp)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews for PR #%d: %w", githubPR.GetNumber(), err)
	}

	var approvers []string
	for _, review := range reviews {
		approver := review.GetUser().GetLogin()
		if review.GetState() == "APPROVED" && !slices.Contains(approvers, approver) {
			approvers = append(approvers, approver)
		}
	}

	return approvers, nil
}

// This needs to be a separate call, because the PullRequests.List method does not support filtering by review requested.
func (c *Client) FetchPRsNeedingMyReview(ctx context.Context) ([]PullRequest, error) {
	query := fmt.Sprintf("is:pr is:open review-requested:%s", c.username)
	for i := range c.repos {
		query += fmt.Sprintf(" repo:%s", c.repos[i])
	}

	opts := &github.SearchOptions{}
	result, _, err := c.github.Search.Issues(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	var prs []PullRequest
	for _, issue := range result.Issues {
		prs = append(prs, *ToPullRequest(issue))
	}

	return prs, nil
}

func getOwnerAndRepo(repoCfg string) (string, string, error) {
	parts := strings.Split(repoCfg, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repo format: %s (expected: owner/repo)", repoCfg)
	}

	return parts[0], parts[1], nil
}
