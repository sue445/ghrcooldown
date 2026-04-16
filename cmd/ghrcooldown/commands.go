package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sue445/ghrcooldown"
)

type CommandLatestParams struct {
	GithubApiURL     string
	GithubToken      string
	GithubRepository string
	CooldownDays     int64
}

func CommandLatest(ctx context.Context, params *CommandLatestParams) error {
	if params.CooldownDays < 0 {
		params.CooldownDays = 0
	}

	client, err := ghrcooldown.NewClient(&ghrcooldown.ClientParams{
		Token:   params.GithubToken,
		BaseURL: params.GithubApiURL,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	repositoryPath, err := ParseRepositoryPath(params.GithubRepository)
	if err != nil {
		return errors.WithStack(err)
	}

	tagName, err := client.GetLatestTagName(ctx, repositoryPath.Owner, repositoryPath.Repo, time.Duration(params.CooldownDays)*ghrcooldown.Day)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println(tagName)

	return nil
}
