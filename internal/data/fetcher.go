package data

import (
	"context"
	"fmt"

	"github.com/pippokairos/workflow-monitor/internal/analyzer"
	"github.com/pippokairos/workflow-monitor/internal/atlassian"
	"github.com/pippokairos/workflow-monitor/internal/config"
	"github.com/pippokairos/workflow-monitor/internal/debug"
	"github.com/pippokairos/workflow-monitor/internal/gh"
)

type Fetcher struct {
	atlassianClient *atlassian.Client
	ghClient        *gh.Client
	matcher         *analyzer.Matcher
	cfg             *config.Config
}

func NewFetcher(cfg *config.Config) (*Fetcher, error) {
	atlassianClient, err := atlassian.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Atlassian client: %v", err)
	}
	debug.Printf("Atlassian client created successfully")

	ghClient := gh.NewClient(cfg)
	debug.Printf("GitHub client created successfully")

	matcher := analyzer.NewMatcher(cfg.IssuePattern)

	return &Fetcher{
		atlassianClient: atlassianClient,
		ghClient:        ghClient,
		matcher:         matcher,
		cfg:             cfg,
	}, nil
}

func (f *Fetcher) FetchAll(ctx context.Context) (*analyzer.Insights, error) {
	myIssues, err := f.atlassianClient.FetchMyIssuesInReviewOrDone()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch my issues: %w", err)
	}
	debug.Printf("Fetched %d issues of mine", len(myIssues))

	openPRs, err := f.ghClient.FetchOpenPRs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch my open PRs: %w", err)
	}
	debug.Printf("Fetched %d open PRs: %+v", len(openPRs), openPRs)

	prsNeedingMyReview, err := f.ghClient.FetchPRsNeedingMyReview(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PRs needing my review: %w", err)
	}
	debug.Printf("Fetched %d PRs needing my review: %+v", len(prsNeedingMyReview), prsNeedingMyReview)

	issueIDToOpenPRs := f.matcher.IssueIDToPRs(openPRs)
	debug.Printf("issueIDToMyOpenPRs: %+v", issueIDToOpenPRs)

	return analyzer.GenerateInsights(myIssues, issueIDToOpenPRs, prsNeedingMyReview, f.cfg)
}
