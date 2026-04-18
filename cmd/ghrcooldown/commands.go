package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sue445/ghrcooldown"
	"github.com/urfave/cli/v3"
)

var errExitCooldownNotPassed = cli.Exit("", 1)

type commandLatestParams struct {
	githubApiURL     string
	githubToken      string
	githubRepository string
	cooldownDays     int64
}

func commandLatest(ctx context.Context, params *commandLatestParams) error {
	if params.cooldownDays < 0 {
		params.cooldownDays = 0
	}

	client, err := ghrcooldown.NewClient(&ghrcooldown.ClientParams{
		Token:   params.githubToken,
		BaseURL: params.githubApiURL,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	repositoryPath, err := parseRepositoryPath(params.githubRepository)
	if err != nil {
		return errors.WithStack(err)
	}

	tagName, err := client.GetLatestTagName(ctx, repositoryPath.Owner, repositoryPath.Repo, time.Duration(params.cooldownDays)*ghrcooldown.Day)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println(tagName)

	return nil
}

type commandHasPassedParams struct {
	githubApiURL     string
	githubToken      string
	githubRepository string
	githubTagName    string
	cooldownDays     int64
	isExitCode       bool
	currentTime      *time.Time
}

func commandHasPassed(ctx context.Context, params *commandHasPassedParams) error {
	if params.cooldownDays < 0 {
		params.cooldownDays = 0
	}

	client, err := ghrcooldown.NewClient(&ghrcooldown.ClientParams{
		Token:       params.githubToken,
		BaseURL:     params.githubApiURL,
		CurrentTime: params.currentTime,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	repositoryPath, err := parseRepositoryPath(params.githubRepository)
	if err != nil {
		return errors.WithStack(err)
	}

	hasPassed, err := client.HasCooldownPassed(ctx, repositoryPath.Owner, repositoryPath.Repo, params.githubTagName, time.Duration(params.cooldownDays)*ghrcooldown.Day)
	if err != nil {
		return errors.WithStack(err)
	}

	if params.isExitCode {
		if !hasPassed {
			return errExitCooldownNotPassed
		}
	} else {
		if hasPassed {
			fmt.Println("Cooldown has passed.")
		} else {
			fmt.Println("Cooldown has not passed yet.")
		}
	}

	return nil
}
