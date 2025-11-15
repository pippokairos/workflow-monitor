package main

import (
	"flag"
	"log"

	"github.com/pippokairos/workflow-monitor/internal/analyzer"
	"github.com/pippokairos/workflow-monitor/internal/atlassian"
	"github.com/pippokairos/workflow-monitor/internal/config"
	"github.com/pippokairos/workflow-monitor/internal/debug"
	"github.com/pippokairos/workflow-monitor/internal/gh"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Enable debug output")
	flag.Parse()
	debug.Enabled = *debugFlag

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	debug.Printf("Config loaded successfully")
	debug.Printf("Config: %+v", cfg)
	debug.Printf("Atlassian URL: %s", cfg.AtlassianURL)
	debug.Printf("Project Keys: %v", cfg.AtlassianProjectKeys)

	atlassianClient, err := atlassian.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Atlassian client: %v", err)
		return
	}
	debug.Printf("Atlassian client created successfully")

	ghClient := gh.NewClient(cfg)
	debug.Printf("GitHub client created successfully")

	myIssues, err := atlassianClient.FetchMyIssuesInReviewOrDone()
	if err != nil {
		log.Fatalf("Failed to fetch my issues: %v", err)
		return
	}
	debug.Printf("Fetched %d issues of mine", len(myIssues))

	openPRs, err := ghClient.FetchOpenPRs()
	if err != nil {
		log.Fatalf("Failed to fetch my open PRs: %v", err)
		return
	}
	debug.Printf("Fetched %d open PRs: %+v", len(openPRs), openPRs)

	prsNeedingMyReview, err := ghClient.FetchPRsNeedingMyReview()
	if err != nil {
		log.Fatalf("Failed to fetch PRs needing my review: %v", err)
		return
	}
	debug.Printf("Fetched %d PRs needing my review: %+v", len(prsNeedingMyReview), prsNeedingMyReview)

	m := analyzer.NewMatcher(cfg.IssuePattern)
	issueIDToOpenPRs := m.IssueIDToPRs(openPRs)
	debug.Printf("issueIDToMyOpenPRs: %+v", issueIDToOpenPRs)

	insights, err := analyzer.GenerateInsights(myIssues, issueIDToOpenPRs, prsNeedingMyReview, cfg)
	if err != nil {
		log.Fatalf("Failed to generate insights: %v", err)
		return
	}

	debug.Printf("Insights: %v", insights)
}
