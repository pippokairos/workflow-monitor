package gh

import "github.com/google/go-github/v79/github"

type PullRequest struct {
	Username  string
	Title     string
	Approvers []string
}

func ToPullRequest(pr *github.PullRequest, approvers []string) *PullRequest {
	var username, title string
	if pr.User.Login != nil {
		username = *pr.User.Login
	}
	if pr.Head != nil && pr.Head.Label != nil {
		title = *pr.Head.Label
	}

	return &PullRequest{
		Username:  username,
		Title:     title,
		Approvers: approvers,
	}
}
