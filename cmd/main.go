package main

import (
	"flag"
	"log"

	"github.com/pippokairos/workflow-monitor/internal/atlassian"
	"github.com/pippokairos/workflow-monitor/internal/config"
	"github.com/pippokairos/workflow-monitor/internal/debug"
	"github.com/pippokairos/workflow-monitor/internal/github"
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

	githubClient := github.NewClient(cfg)
	debug.Printf("GitHub client created successfully")

	doneTickets, err := atlassianClient.FetchDoneTickets()
	if err != nil {
		log.Fatalf("Failed to fetch done tickets: %v", err)
		return
	}
	debug.Printf("Fetched %d done tickets", len(doneTickets))

	myOpenPRs, otherOpenPRs, err := githubClient.FetchOpenPRs()
	if err != nil {
		log.Fatalf("Failed to fetch open PRs: %v", err)
		return
	}
	debug.Printf("Fetched %d open PRs of mine", len(myOpenPRs))
	debug.Printf("Fetched %d open PRs of others", len(otherOpenPRs))
}
