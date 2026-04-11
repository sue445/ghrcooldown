package ghrcooldown_test

import (
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/sue445/ghrcooldown"
)

func TestClient_GetLatestTagName(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/repos/hashicorp/terraform/releases?per_page=10",
		httpmock.NewStringResponder(200, readTestData("testdata/terraform-releases.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"https://github.example.com/api/v3/repos/hashicorp/terraform/releases?per_page=10",
		httpmock.NewStringResponder(200, readTestData("testdata/terraform-releases.json")),
	)

	type args struct {
		owner       string
		repo        string
		cooldown    time.Duration
		currentTime *time.Time
		baseURL     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "cooldown 7 days",
			args: args{
				owner:       "hashicorp",
				repo:        "terraform",
				currentTime: new(time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)),
				cooldown:    days(7),
			},
			want: "v1.14.7",
		},
		{
			name: "cooldown 0 days",
			args: args{
				owner:       "hashicorp",
				repo:        "terraform",
				currentTime: new(time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)),
				cooldown:    days(0),
			},
			want: "v1.14.8",
		},
		{
			name: "with BaseURL",
			args: args{
				owner:       "hashicorp",
				repo:        "terraform",
				currentTime: new(time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)),
				cooldown:    days(7),
				baseURL:     "https://github.example.com/api/v3/",
			},
			want: "v1.14.7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ghrcooldown.NewClient(&ghrcooldown.ClientParams{
				Token:       "DUMMY",
				CurrentTime: tt.args.currentTime,
				BaseURL:     tt.args.baseURL,
			})

			if assert.NoError(t, err) {
				got, err := c.GetLatestTagName(t.Context(), tt.args.owner, tt.args.repo, tt.args.cooldown)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func TestClient_HasCooldownPassed(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/repos/hashicorp/terraform/releases/tags/v1.14.8",
		httpmock.NewStringResponder(200, readTestData("testdata/terraform-tags-v1.14.8.json")),
	)

	type args struct {
		owner       string
		repo        string
		cooldown    time.Duration
		tagName     string
		currentTime *time.Time
		baseURL     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "cooldown 7 days",
			args: args{
				owner:       "hashicorp",
				repo:        "terraform",
				currentTime: new(time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)),
				cooldown:    days(7),
				tagName:     "v1.14.8",
			},
			want: false,
		},
		{
			name: "cooldown 0 days",
			args: args{
				owner:       "hashicorp",
				repo:        "terraform",
				currentTime: new(time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)),
				cooldown:    days(0),
				tagName:     "v1.14.8",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ghrcooldown.NewClient(&ghrcooldown.ClientParams{
				Token:       "DUMMY",
				CurrentTime: tt.args.currentTime,
				BaseURL:     tt.args.baseURL,
			})

			if assert.NoError(t, err) {
				got, err := c.HasCooldownPassed(t.Context(), tt.args.owner, tt.args.repo, tt.args.cooldown, tt.args.tagName)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}
