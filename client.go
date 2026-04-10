package ghrcooldown

import "time"

// Client represents a client for interacting with the GitHub API with cooldown support.
type Client struct {
	token string
}

// ClientParams contains the configuration parameters required to initialize a new Client.
type ClientParams struct {
	Token       string
	BaseURL     string
	CurrentTime *time.Time
}

// NewClient creates and returns a new Client instance using the provided parameters.
func NewClient(params *ClientParams) (*Client, error) {
	return &Client{}, nil
}
