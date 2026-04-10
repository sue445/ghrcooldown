package ghrcooldown

import (
	"context"
	"net/http"
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
	HTTPClient  *http.Client
}

// NewClient creates and returns a new Client instance using the provided parameters.
func NewClient(params *ClientParams) (*Client, error) {
	client := github.NewClient(params.HTTPClient)

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

// GetLatestVersion retrieves the latest release version of the specified repository, respecting the provided cooldown period.
func (c *Client) GetLatestVersion(ctx context.Context, owner string, repo string, cooldown time.Duration) (string, error) {
	// TODO: Impl
	return "", nil
}
