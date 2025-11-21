package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pippokairos/workflow-monitor/internal/analyzer"
	"github.com/pippokairos/workflow-monitor/internal/data"
)

var (
	primaryColor   = lipgloss.Color("#FFB86C") // Orange
	secondaryColor = lipgloss.Color("#00FF87") // Bright green
	tertiaryColor  = lipgloss.Color("#000000") // Black
	errorColor     = lipgloss.Color("#FF5555") // Red
	subtleColor    = lipgloss.Color("#6272A4") // Gray

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(tertiaryColor).
			Background(primaryColor).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(subtleColor).
				Padding(0, 2)

	cursorStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	numberStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	successBadgeStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true)

	footerStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(subtleColor).
			Padding(1, 0, 0, 0)

	statsStyle = lipgloss.NewStyle().
			Italic(true)
)

type state int

const (
	stateLoading state = iota
	stateReady
	stateError
)

type model struct {
	state     state
	err       error
	fetcher   *data.Fetcher
	insights  *analyzer.Insights
	spinner   spinner.Model
	startTime time.Time
	loadTime  time.Duration

	selectedView int // 0, 1, or 2 for the three views
	cursor       int // Selected item
}

type fetchCompleteMsg struct {
	insights *analyzer.Insights
	err      error
	duration time.Duration
}

func InitialModel(fetcher *data.Fetcher) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	return model{
		state:        stateLoading,
		fetcher:      fetcher,
		selectedView: 0,
		cursor:       0,
		spinner:      s,
		startTime:    time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchDataCmd(m.fetcher, m.startTime),
	)
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
		m.loadTime = msg.duration
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

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
			m.startTime = time.Now()
			return m, tea.Batch(
				m.spinner.Tick,
				fetchDataCmd(m.fetcher, m.startTime),
			)
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
		elapsed := time.Since(m.startTime).Round(100 * time.Millisecond)
		return fmt.Sprintf("\n  %s Fetching tickets and PRs from Jira and GitHub... %s\n\n",
			m.spinner.View(),
			subtitleStyle.Render(fmt.Sprintf("(%s)", elapsed)))
	}

	if m.state == stateError {
		errorMsg := lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Render(fmt.Sprintf("Error: %v", m.err))
		return fmt.Sprintf("\n%s\n\n%s\n", errorMsg, subtitleStyle.Render("Press any key to exit."))
	}

	// Header with tabs
	tabs := []string{
		fmt.Sprintf("Ticket done, PRs not merged (%d)", len(m.insights.DoneNotMergedPRs)),
		fmt.Sprintf("Need Review (%d)", len(m.insights.NeedReviewPRs)),
		fmt.Sprintf("Ready for QA (%d)", len(m.insights.ReviewedNotInQAPRs)),
	}

	// Add load time
	header := "\n" + statsStyle.Render(fmt.Sprintf("  Data loaded in %s", m.loadTime.Round(10*time.Millisecond))) + "\n\n"

	for i, tab := range tabs {
		if i == m.selectedView {
			header += activeTabStyle.Render(tab)
		} else {
			header += inactiveTabStyle.Render(tab)
		}
		header += " "
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
	footer := footerStyle.Render(
		"Tab: switch view | ↑/↓ or j/k: navigate | Enter: open in browser | r: refresh | q: quit",
	)

	return header + content + footer
}

func (m model) renderDoneNotMergedPRs() string {
	if len(m.insights.DoneNotMergedPRs) == 0 {
		return renderNoItemsFoundMessage()
	}

	s := ""
	for i, item := range m.insights.DoneNotMergedPRs {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▸ ")
		}

		ticketBadge := numberStyle.Render(fmt.Sprintf("%s", item.IssueID))
		title := titleStyle.Render(item.PullRequest.Title)
		prInfo := subtitleStyle.Render(fmt.Sprintf("PR #%d by %s (open)", item.PullRequest.Number, item.PullRequest.Author))

		s += fmt.Sprintf("%s%s %s\n    %s\n\n", cursor, ticketBadge, title, prInfo)
	}

	return s
}

func (m model) renderNeedReviewPRs() string {
	if len(m.insights.NeedReviewPRs) == 0 {
		return renderNoItemsFoundMessage()
	}

	s := ""
	for i, pr := range m.insights.NeedReviewPRs {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▸ ")
		}

		prBadge := numberStyle.Render(fmt.Sprintf("#%d", pr.Number))
		title := titleStyle.Render(pr.Title)
		info := subtitleStyle.Render(fmt.Sprintf("PR by %s in %s", pr.Author, pr.Repo))

		s += fmt.Sprintf("%s%s %s\n    %s\n\n", cursor, prBadge, title, info)
	}

	return s
}

func (m model) renderReviewedNotInQAPRs() string {
	if len(m.insights.ReviewedNotInQAPRs) == 0 {
		return renderNoItemsFoundMessage()
	}

	s := ""
	for i, item := range m.insights.ReviewedNotInQAPRs {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▸ ")
		}

		ticketBadge := successBadgeStyle.Render(fmt.Sprintf("%s", item.IssueID))
		title := titleStyle.Render(item.PullRequest.Title)

		approvers := fmt.Sprintf("PR #%d", item.PullRequest.Number)
		if len(item.PullRequest.Approvers) > 0 {
			approvers = subtitleStyle.Render(fmt.Sprintf("%s approved by: %s", approvers, strings.Join(item.PullRequest.Approvers, ", ")))
		} else {
			approvers = subtitleStyle.Render(fmt.Sprintf("%s - no approvals yet", approvers))
		}

		s += fmt.Sprintf("%s%s %s\n    %s\n\n", cursor, ticketBadge, title, approvers)
	}

	return s
}

func renderNoItemsFoundMessage() string {
	return successBadgeStyle.Render("No items found!") + "\n\n\n"
}
