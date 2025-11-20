package analyzer

import (
	"regexp"

	"github.com/pippokairos/workflow-monitor/internal/debug"
	"github.com/pippokairos/workflow-monitor/internal/gh"
)

type Matcher struct {
	issuePattern regexp.Regexp
}

func NewMatcher(pattern string) *Matcher {
	return &Matcher{
		issuePattern: *regexp.MustCompile(pattern),
	}
}

func (m *Matcher) IssueIDToPRs(prs []gh.PullRequest) map[string][]gh.PullRequest {
	issueIDToPRs := make(map[string][]gh.PullRequest, len(prs))

	for _, pr := range prs {
		issueID := m.getIssueID(pr.BranchName)
		if issueID == nil {
			debug.Printf("No issue ID found in branch name: %s", pr.BranchName)
			continue
		}

		issueIDToPRs[*issueID] = append(issueIDToPRs[*issueID], pr)
	}

	return issueIDToPRs
}

func (m *Matcher) getIssueID(branchName string) *string {
	match := m.issuePattern.Find([]byte(branchName))
	if match == nil {
		return nil
	}

	id := string(match)
	return &id
}
