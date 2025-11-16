package atlassian

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/pippokairos/workflow-monitor/internal/config"
	"github.com/pippokairos/workflow-monitor/internal/debug"
)

type Client struct {
	jira         *jira.Client
	statusReview string
	statusDone   string
	projectKeys  []string
}

func NewClient(cfg *config.Config) (*Client, error) {
	tp := jira.BasicAuthTransport{
		Username: cfg.AtlassianEmail,
		Password: cfg.AtlassianToken,
	}

	client, err := jira.NewClient(tp.Client(), cfg.AtlassianURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		jira:         client,
		statusReview: cfg.AtlassianStatusReview,
		statusDone:   cfg.AtlassianStatusDone,
		projectKeys:  cfg.AtlassianProjectKeys,
	}, nil
}

func (c *Client) FetchMyIssuesInReviewOrDone() ([]jira.Issue, error) {
	jql := fmt.Sprintf("assignee = currentUser() AND updated >= -14d AND status IN (%s, %s)", c.statusReview, c.statusDone)
	if len(c.projectKeys) > 0 {
		jql += fmt.Sprintf(" AND project IN (%s)", strings.Join(c.projectKeys, ","))
	}
	jql += " ORDER BY updated DESC"

	options := &jira.SearchOptionsV2{Fields: []string{"*all"}}

	// TODO: handle pagination?
	issues, resp, err := c.jira.Issue.SearchV2JQL(jql, options)
	debug.Printf("Jira response: %+v", resp)
	debug.Printf("Response: %+v", resp.Response)
	if err != nil {
		return nil, err
	}

	return issues, err
}
