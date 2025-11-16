# Workflow Monitor

A terminal-based tool to monitor the sync status between Jira tickets and GitHub pull requests.

## Features

- **Ticket Done, PRs Not Merged**: Find Jira tickets marked as "Done" but with open PRs
- **Need Review**: See PRs waiting for your review
- **Ready for QA**: Identify approved PRs whose tickets haven't been moved to QA
- **Real-time Refresh**: Update data on demand

## Prerequisites

- Go 1.24
- Jira account with API access
- GitHub account with personal access token
- Access to your organization's Jira projects and GitHub repositories

## Installation

### Clone the repository

```bash
git clone https://github.com/pippokairos/workflow-monitor.git
cd workflow-monitor
```

### Install dependencies

```bash
go mod download
```

### Build the binary

```bash
go build -o workflow-monitor cmd/main.go
```

Or run directly:

```bash
go run cmd/main.go
```

## Configuration

### 1. Create API tokens

#### Jira API Token

1. Go to [Atlassian Account Settings](https://id.atlassian.com/manage-profile/security/api-tokens)
2. Click "Create API token"
3. Give it a name (e.g., "Workflow Monitor")
4. Copy the token and save it.

#### GitHub Personal Access Token

1. Go to [GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)](https://github.com/settings/tokens)
2. Click "Generate new token (classic)"
3. Select scopes:
   - `repo` (Full control of private repositories)
   - `read:org` (Read org and team membership) - optional
4. Generate and copy the token.

### 2. Create config.yaml

Create a `config.yaml` file in the project root following the [example](https://github.com/pippokairos/workflow-monitor/blob/main/config.yml.example)

**Important:** GitHub repos must be in `owner/repo` format (e.g., `myorg/api`, not just `api`)

## Usage

### Basic usage

```bash
./workflow-monitor
```

Or with Go:

```bash
go run cmd/main.go
```

### With debug mode

See detailed API calls and responses:

```bash
./workflow-monitor -debug
```

### Keyboard shortcuts

Once the TUI is running:

- **Tab** - Switch between views
- **↑/↓** or **j/k** - Navigate through items
- **Enter** - Open selected PR/ticket in browser
- **r** - Refresh data
- **q** or **Ctrl+C** - Quit

## How It Works

### 1. Ticket Done, PRs Not Merged

Finds Jira tickets that are:

- Assigned to you
- Status = "Done" (or as defined)
- Have matching PRs that are still open

**Use case:** Remind you to merge PRs after tickets are completed.

### 2. Need Review

Shows GitHub PRs where:

- You're requested as a reviewer
- You haven't submitted a review yet
- PR is still open

**Use case:** Don't miss review requests.

### 3. Ready for QA

Identifies your PRs that:

- Have the minimum amount of approvals
- The linked Jira ticket is NOT in "QA" status

**Use case:** Remember to move tickets to QA after PRs are approved.

### Ticket-PR Matching

The tool matches Jira tickets to GitHub PRs by:

1. Extracting the ticket ID from the PR's branch name using the regex pattern in config
2. Example: branch `feature/PROJ-123-add-login` → ticket `PROJ-123`

### Rate limiting

GitHub API has generous rate limits:

- Authenticated: 5,000 requests/hour
- Unauthenticated: 60 requests/hour

If you hit the limit, wait an hour or reduce the number of repos in your config.

## Development

### Project structure

```
workflow-monitor/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── analyzer/            # Matching and insights generation
│   │   ├── matcher.go
│   │   └── insights.go
│   ├── atlassian/           # Jira client
│   │   └── client.go
│   ├── config/              # Configuration loading
│   │   ├── config.go
│   │   └── config_test.go
│   ├── data/                # Data orchestration layer
│   │   └── fetcher.go
│   ├── debug/               # Debug utilities
│   │   └── debug.go
│   ├── gh/                  # GitHub client
│   │   ├── client.go
│   │   ├── client_test.go
│   │   └── types.go
│   └── ui/                  # Terminal UI
│       ├── commands.go
│       └── tui.go
├── config.yml               # Your configuration
├── go.mod
├── go.sum
└── README.md
```

### Running tests

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see the [LICENSE](https://github.com/pippokairos/workflow-monitor/blob/main/LICENSE) file for details

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling library
- [go-jira](https://github.com/andygrunwald/go-jira) - Jira API client
- [go-github](https://github.com/google/go-github) - GitHub API client
