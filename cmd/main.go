package main

import (
	"flag"
	"log"

	"github.com/pippokairos/workflow-monitor/internal/atlassian"
	"github.com/pippokairos/workflow-monitor/internal/config"
	"github.com/pippokairos/workflow-monitor/internal/debug"
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

	debug.Printf("Atlassian client created successfully!")

	doneTickets, err := atlassianClient.FetchDoneTickets(cfg.AtlassianProjectKeys)
	if err != nil {
		log.Fatalf("Failed to fetch done tickets: %v", err)
		return
	}

	debug.Printf("Fetched %d done tickets", len(doneTickets))
}
