package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sue445/ghrcooldown"
	"github.com/urfave/cli/v3"
)

var (
	// Revision represents app revision (injected from ldflags)
	Revision string
)

// RepositoryPath represents GitHub repository's owner and repo
type RepositoryPath struct {
	Owner string
	Repo  string
}

func main() {
	var githubApiURL string
	var githubToken string
	var githubRepository string
	var cooldownDays int64

	commonFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "github-api-url",
			Usage:       "GitHub API Endpoint (e.g. https://<your-ghes-hostname>/api/v3). Required if using GitHub Enterprise Server",
			Sources:     cli.EnvVars("GITHUB_API_URL"),
			Required:    false,
			Destination: &githubApiURL,
		},
		&cli.StringFlag{
			Name:        "token",
			Usage:       "GitHub token",
			Sources:     cli.EnvVars("GITHUB_TOKEN"),
			Required:    false,
			Destination: &githubToken,
		},
		&cli.StringFlag{
			Name:        "repo",
			Usage:       "GitHub Repository Path (e.g. user/repo)",
			Required:    true,
			Destination: &githubRepository,
		},
		&cli.Int64Flag{
			Name:        "cooldown-days",
			Usage:       "Cooldown days",
			Required:    false,
			Destination: &cooldownDays,
			Value:       0,
		},
	}

	cmd := &cli.Command{
		Name:    "ghrcooldown",
		Version: fmt.Sprintf("%s (build. %s)", ghrcooldown.Version, Revision),
		Usage:   "Get the latest GitHub Releases respecting the cooldown period.",
		Commands: []*cli.Command{
			{
				Name:  "latest",
				Usage: "Print latest release version of the specified repository, respecting the provided cooldown period.",
				Flags: commonFlags,
				Action: func(ctx context.Context, _ *cli.Command) error {
					return CommandLatest(ctx, &CommandLatestParams{
						GithubApiURL:     githubApiURL,
						GithubToken:      githubToken,
						GithubRepository: githubRepository,
						CooldownDays:     cooldownDays,
					})
				},
			},
		},
	}

	// Sort commands
	sort.Slice(cmd.Commands, func(i, j int) bool {
		return cmd.Commands[i].Name < cmd.Commands[j].Name
	})

	// Sort sub-command flags
	for _, c := range cmd.Commands {
		sort.Sort(cli.FlagsByName(c.Flags))
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
}

// ParseRepositoryPath parses a GitHub repository path string (e.g., "owner/repo") and returns a [RepositoryPath]. It returns an error if the format is invalid.
func ParseRepositoryPath(path string) (*RepositoryPath, error) {
	path = strings.TrimSpace(path)
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return nil, errors.Errorf("invalid repository path: %s", path)
	}

	owner := strings.TrimSpace(parts[0])
	repo := strings.TrimSpace(parts[1])
	if owner == "" || repo == "" {
		return nil, errors.Errorf("invalid repository path: %s", path)
	}

	return &RepositoryPath{
		Owner: owner,
		Repo:  repo,
	}, nil
}
