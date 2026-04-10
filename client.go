package ghrcooldown

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/go-github/v84/github"
)

// Client represents a client for interacting with the GitHub API with cooldown support.
type Client struct {
	client      *github.Client
	currentTime *time.Time
}

// ClientParams contains the configuration parameters required to initialize a new Client.
type ClientParams struct {
	Token       string
	BaseURL     string
	CurrentTime *time.Time
}

// NewClient creates and returns a new Client instance using the provided parameters.
func NewClient(params *ClientParams) (*Client, error) {
	client := github.NewClient(nil)

	if params.Token != "" {
		client = client.WithAuthToken(params.Token)
	}

	if params.BaseURL != "" {
		var err error
		client, err = client.WithEnterpriseURLs(params.BaseURL, "")
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &Client{client: client, currentTime: params.CurrentTime}, nil
}

// GetLatestTagName retrieves the latest release version of the specified repository, respecting the provided cooldown period.
func (c *Client) GetLatestTagName(ctx context.Context, owner string, repo string, cooldown time.Duration) (string, error) {
	currentTime := time.Now()
	if c.currentTime != nil {
		currentTime = *c.currentTime
	}

	opt := &github.ListOptions{
		PerPage: 10,
	}

	for {
		releases, resp, err := c.client.Repositories.ListReleases(ctx, owner, repo, opt)
		if err != nil {
			return "", errors.WithStack(err)
		}

		for _, release := range releases {
			if release.GetDraft() || release.GetPrerelease() {
				continue
			}

			ts := release.GetPublishedAt()
			if ts.IsZero() {
				continue
			}

			publishedAt := ts.Time

			if currentTime.Sub(publishedAt) >= cooldown {
				return release.GetTagName(), nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return "", fmt.Errorf("no release found that respects the cooldown period of %v", cooldown)
}
