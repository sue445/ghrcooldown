package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sue445/ghrcooldown"
)

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
