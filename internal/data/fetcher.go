package data

import (
	"context"
	"fmt"
	"sync"

	"github.com/andygrunwald/go-jira"
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
	var myIssues []jira.Issue
	var openPRs, prsNeedingMyReview []gh.PullRequest
	var myIssuesErr, openPRsErr, prsNeedingMyReviewErr error

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		myIssues, myIssuesErr = f.atlassianClient.FetchMyIssuesInReviewOrDone()
		debug.Printf("Fetched %d issues of mine", len(myIssues))
		debug.Printf("My issues: %+v", myIssues)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		openPRs, openPRsErr = f.ghClient.FetchOpenPRs(ctx)
		debug.Printf("Fetched %d open PRs: %+v", len(openPRs), openPRs)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		prsNeedingMyReview, prsNeedingMyReviewErr = f.ghClient.FetchPRsNeedingMyReview(ctx)
		debug.Printf("Fetched %d PRs needing my review: %+v", len(prsNeedingMyReview), prsNeedingMyReview)
	}()

	wg.Wait()

	if myIssuesErr != nil {
		return nil, fmt.Errorf("failed to fetch my issues: %v", myIssuesErr)
	}
	if openPRsErr != nil {
		return nil, fmt.Errorf("failed to fetch open PRs: %v", openPRsErr)
	}
	if prsNeedingMyReviewErr != nil {
		return nil, fmt.Errorf("failed to fetch PRs needing my review: %v", prsNeedingMyReviewErr)
	}

	issueIDToOpenPRs := f.matcher.IssueIDToPRs(openPRs)
	debug.Printf("issueIDToOpenPRs: %+v", issueIDToOpenPRs)

	return analyzer.GenerateInsights(myIssues, issueIDToOpenPRs, prsNeedingMyReview, f.cfg)
}
