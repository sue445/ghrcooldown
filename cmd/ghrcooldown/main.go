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
	cmd := &cli.Command{
		Name:    "ghrcooldown",
		Version: fmt.Sprintf("%s (build. %s)", ghrcooldown.Version, Revision),
		Usage:   "Fetch the latest GitHub Releases respecting the cooldown period.",
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
