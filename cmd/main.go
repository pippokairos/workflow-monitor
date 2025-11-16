package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pippokairos/workflow-monitor/internal/config"
	"github.com/pippokairos/workflow-monitor/internal/data"
	"github.com/pippokairos/workflow-monitor/internal/debug"
	"github.com/pippokairos/workflow-monitor/internal/ui"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Enable debug output")
	flag.Parse()
	debug.Enabled = *debugFlag

	cfg, err := config.Load("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	debug.Printf("Config loaded successfully\n")
	debug.Printf("Config: %+v", cfg)
	debug.Printf("Atlassian URL: %s", cfg.AtlassianURL)
	debug.Printf("Project Keys: %v", cfg.AtlassianProjectKeys)

	fetcher, err := data.NewFetcher(cfg)
	if err != nil {
		log.Fatalf("Failed to create data fetcher: %v", err)
	}

	if debug.Enabled {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		insights, err := fetcher.FetchAll(ctx)
		if err != nil {
			log.Fatalf("Failed to fetch data: %v", err)
		}

		debug.Printf("Insights: %v", insights)
		return
	}

	p := tea.NewProgram(ui.InitialModel(fetcher), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
