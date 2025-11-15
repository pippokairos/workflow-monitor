package atlassian

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/pippokairos/workflow-monitor/internal/config"
	"github.com/pippokairos/workflow-monitor/internal/debug"
)

type Client struct {
	jira        *jira.Client
	projectKeys []string
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
		jira:        client,
		projectKeys: cfg.AtlassianProjectKeys,
	}, nil
}

func (c *Client) FetchDoneTickets() ([]jira.Issue, error) {
	jql := "assignee = currentUser() AND status = Done AND updated >= -30d"
	if len(c.projectKeys) > 0 {
		jql += fmt.Sprintf(" AND project IN (%s)", strings.Join(c.projectKeys, ","))
	}
	debug.Printf("jql: %+v", jql)

	options := &jira.SearchOptionsV2{Fields: []string{"*all"}}

	// TODO: handle pagination?
	issues, resp, err := c.jira.Issue.SearchV2JQL(jql, options)

	debug.Printf("Jira response: %+v", resp)
	debug.Printf("Response: %+v", resp.Response)

	return issues, err
}
