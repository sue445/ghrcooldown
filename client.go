package ghrcooldown

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/go-github/v84/github"
)

// Day represents the duration of exactly 24 hours.
const Day = 24 * time.Hour

// Client represents a client for interacting with the GitHub API with cooldown support.
type Client struct {
	client      *github.Client
	currentTime *time.Time
}

// ClientParams contains the configuration parameters required to initialize a new Client.
type ClientParams struct {
	// Token is the personal access token used for authenticating with the GitHub API.
	// It is optional, but if omitted, the API request will be subject to IP-based rate limiting.
	// c.f. https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api
	Token string

	// BaseURL specifies the custom base URL for the GitHub API.
	// It is primarily used for GitHub Enterprise Server.
	BaseURL string

	// UserAgent specifies the User-Agent header used in API requests.
	// If omitted, a default User-Agent will be used.
	UserAgent string

	// CurrentTime is the reference time used to evaluate the cooldown period.
	// It is mainly used for mocking the current time in unit tests.
	CurrentTime *time.Time
}

// NewClient creates and returns a new Client instance using the provided parameters.
func NewClient(params *ClientParams) (*Client, error) {
	client := github.NewClient(nil)

	if params.UserAgent == "" {
		client.UserAgent = GetDefaultUserAgent()
	} else {
		client.UserAgent = params.UserAgent
	}

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

// GetDefaultUserAgent returns the default User-Agent.
func GetDefaultUserAgent() string {
	return fmt.Sprintf("ghrcooldown/%s (+https://github.com/sue445/ghrcooldown)", Version)
}

// HasCooldownPassed checks if the specified tag has passed the given cooldown period.
func (c *Client) HasCooldownPassed(ctx context.Context, owner string, repo string, tagName string, cooldown time.Duration) (bool, error) {
	release, _, err := c.client.Repositories.GetReleaseByTag(ctx, owner, repo, tagName)
	if err != nil {
		return false, errors.WithStack(err)
	}

	ts := release.GetPublishedAt()

	if ts.IsZero() {
		return false, nil
	}

	now := time.Now()
	if c.currentTime != nil {
		now = *c.currentTime
	}

	return now.Sub(ts.Time) >= cooldown, nil
}
