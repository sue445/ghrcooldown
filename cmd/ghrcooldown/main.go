package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/cockroachdb/errors"
	"github.com/sue445/ghrcooldown"
	"github.com/urfave/cli/v3"
)

var (
	// Revision represents app revision (injected from ldflags)
	Revision string
)

func main() {
	var githubApiURL string
	var githubToken string
	var githubRepository string
	var cooldownDays int64

	commonFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "github-api-url",
			Usage:       "GitHub API Endpoint (e.g. https://<your-ghes-hostname>/api/v3)",
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
					return nil
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
		log.Fatal(errors.WithStack(err))
	}
}
