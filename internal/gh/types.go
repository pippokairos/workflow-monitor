package gh

import (
	"strings"

	"github.com/google/go-github/v79/github"
)

type PullRequest struct {
	URL        string
	Number     int
	Title      string
	State      string
	BranchName string
	Author     string
	Repo       string
	Approvers  []string
}

func ToInternalPullRequest(pr *github.PullRequest, approvers []string) *PullRequest {
	var author string
	if pr.User.Login != nil {
		author = *pr.User.Login
	}

	return &PullRequest{
		URL:        pr.GetHTMLURL(),
		Number:     pr.GetNumber(),
		Title:      pr.GetTitle(),
		State:      pr.GetState(),
		BranchName: pr.GetHead().GetRef(),
		Author:     author,
		Repo:       pr.GetBase().GetRepo().GetFullName(),
		Approvers:  approvers,
	}
}

func ToPullRequest(issue *github.Issue) *PullRequest {
	var author, title, repo string
	if issue.User.Login != nil {
		author = *issue.User.Login
	}
	if issue.Title != nil {
		title = *issue.Title
	}
	if issue.RepositoryURL != nil {
		parts := strings.Split(*issue.RepositoryURL, "/")
		repo = parts[len(parts)-1]
	}

	return &PullRequest{
		URL:        issue.GetHTMLURL(),
		Number:     issue.GetNumber(),
		Title:      title,
		State:      issue.GetState(),
		BranchName: "", // N/A
		Author:     author,
		Repo:       repo,
		Approvers:  []string{}, // N/A
	}
}
