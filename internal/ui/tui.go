package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pippokairos/workflow-monitor/internal/analyzer"
	"github.com/pippokairos/workflow-monitor/internal/data"
)

type state int

const (
	stateLoading state = iota
	stateReady
	stateError
)

const noItemsFound = "No items found!"

type model struct {
	state    state
	err      error
	fetcher  *data.Fetcher
	insights *analyzer.Insights

	selectedView int // 0, 1, or 2 for the three views
	cursor       int // Selected item
}

type fetchCompleteMsg struct {
	insights *analyzer.Insights
	err      error
}

func InitialModel(fetcher *data.Fetcher) model {
	return model{
		state:        stateLoading,
		fetcher:      fetcher,
		selectedView: 0,
		cursor:       0,
	}
}

func (m model) Init() tea.Cmd {
	return fetchDataCmd(m.fetcher)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case fetchCompleteMsg:
		if msg.err != nil {
			m.state = stateError
			m.err = msg.err
			return m, tea.Quit
		}
		m.state = stateReady
		m.insights = msg.insights
		return m, nil

	case tea.KeyMsg:
		if m.state == stateLoading {
			return m, nil
		}

		if m.state == stateError {
			return m, tea.Quit
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			// Switch view
			m.selectedView = (m.selectedView + 1) % 3
			m.cursor = 0
			return m, nil

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			max := m.getMaxCursor()
			if m.cursor < max {
				m.cursor++
			}
			return m, nil

		case "enter":
			url := m.getSelectedURL()
			if url != "" {
				return m, openURLCmd(url)
			}
			return m, nil

		case "r":
			// Refresh
			m.state = stateLoading
			m.cursor = 0
			return m, fetchDataCmd(m.fetcher)
		}
	}

	return m, nil
}

func (m model) getMaxCursor() int {
	switch m.selectedView {
	case 0:
		return len(m.insights.DoneNotMergedPRs) - 1
	case 1:
		return len(m.insights.NeedReviewPRs) - 1
	case 2:
		return len(m.insights.ReviewedNotInQAPRs) - 1
	}
	return 0
}

func (m model) getSelectedURL() string {
	switch m.selectedView {
	case 0:
		if m.cursor < len(m.insights.DoneNotMergedPRs) {
			return m.insights.DoneNotMergedPRs[m.cursor].PullRequest.URL
		}
	case 1:
		if m.cursor < len(m.insights.NeedReviewPRs) {
			return m.insights.NeedReviewPRs[m.cursor].URL
		}
	case 2:
		if m.cursor < len(m.insights.ReviewedNotInQAPRs) {
			return m.insights.ReviewedNotInQAPRs[m.cursor].PullRequest.URL
		}
	}

	return ""
}

func (m model) View() string {
	if m.state == stateLoading {
		return "Loading data...\n\nFetching tickets and PRs, please wait..."
	}

	if m.state == stateError {
		return fmt.Sprintf("Error: %v\n\nPress any key to exit.", m.err)
	}

	// Header
	tabs := []string{
		fmt.Sprintf("Ticket done, PRs not merged (%d)", len(m.insights.DoneNotMergedPRs)),
		fmt.Sprintf("Need Review (%d)", len(m.insights.NeedReviewPRs)),
		fmt.Sprintf("Ready for QA (%d)", len(m.insights.ReviewedNotInQAPRs)),
	}

	header := ""
	for i, tab := range tabs {
		if i == m.selectedView {
			header += fmt.Sprintf("[%s] ", tab)
		} else {
			header += fmt.Sprintf(" %s  ", tab)
		}
	}
	header += "\n\n"

	// Content
	var content string
	switch m.selectedView {
	case 0:
		content = m.renderDoneNotMergedPRs()
	case 1:
		content = m.renderNeedReviewPRs()
	case 2:
		content = m.renderReviewedNotInQAPRs()
	}

	// Footer
	footer := "\n\nTab: switch view | ↑/↓ or j/k: navigate | Enter: open in browser | r: refresh | q: quit"

	return header + content + footer
}

func (m model) renderDoneNotMergedPRs() string {
	if len(m.insights.DoneNotMergedPRs) == 0 {
		return noItemsFound
	}

	s := ""
	for i, item := range m.insights.DoneNotMergedPRs {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, item.IssueID, item.PullRequest.Title)
		s += fmt.Sprintf("  PR #%d by %s (open)\n\n", item.PullRequest.Number, item.PullRequest.Author)
	}

	return s
}

func (m model) renderNeedReviewPRs() string {
	if len(m.insights.NeedReviewPRs) == 0 {
		return noItemsFound
	}

	s := ""
	for i, pr := range m.insights.NeedReviewPRs {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s PR #%d: %s\n", cursor, pr.Number, pr.Title)
		s += fmt.Sprintf("  by %s in %s\n\n", pr.Author, pr.Repo)
	}

	return s
}

func (m model) renderReviewedNotInQAPRs() string {
	if len(m.insights.ReviewedNotInQAPRs) == 0 {
		return noItemsFound
	}

	s := ""
	for i, item := range m.insights.ReviewedNotInQAPRs {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, item.IssueID, item.PullRequest.Title)
		s += fmt.Sprintf("  PR #%d approved by: %v\n", item.PullRequest.Number, item.PullRequest.Approvers)
		s += "  Ticket not in QA\n\n"
	}

	return s
}
