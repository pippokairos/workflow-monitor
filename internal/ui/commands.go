package ui

import (
	"context"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pippokairos/workflow-monitor/internal/data"
)

func fetchDataCmd(fetcher *data.Fetcher) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		insights, err := fetcher.FetchAll(ctx)

		return fetchCompleteMsg{
			insights: insights,
			err:      err,
		}
	}
}

func openURLCmd(url string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", url)
		default:
			return nil
		}

		_ = cmd.Start()
		return nil
	}
}
