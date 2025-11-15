package analyzer

import (
	"github.com/andygrunwald/go-jira"
	"github.com/pippokairos/workflow-monitor/internal/config"
	"github.com/pippokairos/workflow-monitor/internal/gh"
)

type Insights struct {
	DoneNotMergedPRs   []DoneNotMergedPR
	NeedReviewPRs      []ReviewNeededPR
	ReviewedNotInQAPRs []ReviewedNotInQAPR
}

type DoneNotMergedPR struct {
	IssueID     string
	PullRequest gh.PullRequest
}

type ReviewNeededPR gh.PullRequest

type ReviewedNotInQAPR struct {
	IssueID     string
	PullRequest gh.PullRequest
	Approvers   []string
}

func GenerateInsights(
	issues []jira.Issue,
	issueIDToOpenPRs map[string][]gh.PullRequest,
	prsNeedingMyReview []gh.PullRequest,
	cfg *config.Config,
) (*Insights, error) {
	return &Insights{
		DoneNotMergedPRs:   GetDoneNotMergedPRs(issues, issueIDToOpenPRs, cfg.AtlassianStatusDone),
		NeedReviewPRs:      GetReviewNeededPRs(prsNeedingMyReview),
		ReviewedNotInQAPRs: GetReviewedNotInQAPRs(issues, issueIDToOpenPRs, cfg.GitHubRequiredApprovers),
	}, nil
}

func GetDoneNotMergedPRs(issues []jira.Issue, issueIDToOpenPRs map[string][]gh.PullRequest, statusDone string) []DoneNotMergedPR {
	doneNotMergedPRs := make([]DoneNotMergedPR, 0, len(issues))

	for i := range issues {
		if issues[i].Fields == nil ||
			issues[i].Fields.Status == nil ||
			issues[i].Fields.Status.Name != statusDone {
			continue
		}

		issueID := issues[i].Key
		prs, ok := issueIDToOpenPRs[issueID]
		if !ok || len(prs) == 0 {
			continue
		}

		for j := range prs {
			doneNotMergedPRs = append(doneNotMergedPRs, DoneNotMergedPR{
				IssueID:     issueID,
				PullRequest: prs[j],
			})
		}
	}

	return doneNotMergedPRs
}

func GetReviewNeededPRs(prsNeedingMyReview []gh.PullRequest) []ReviewNeededPR {
	reviewNeededPRs := make([]ReviewNeededPR, len(prsNeedingMyReview))
	for i := range prsNeedingMyReview {
		reviewNeededPRs[i] = ReviewNeededPR(prsNeedingMyReview[i])
	}

	return reviewNeededPRs
}

func GetReviewedNotInQAPRs(issues []jira.Issue, issueIDToOpenPRs map[string][]gh.PullRequest, requiredApprovers int) []ReviewedNotInQAPR {
	reviewedNotInQAPRs := make([]ReviewedNotInQAPR, 0)
	for i := range issues {
		issueID := issues[i].Key
		prs, ok := issueIDToOpenPRs[issueID]
		if !ok || len(prs) == 0 {
			continue
		}

		for j := range prs {
			if len(prs[j].Approvers) >= requiredApprovers {
				reviewedNotInQAPRs = append(reviewedNotInQAPRs, ReviewedNotInQAPR{})
			}
		}
	}

	return reviewedNotInQAPRs
}
