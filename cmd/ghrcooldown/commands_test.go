package main

import (
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_commandLatest(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/repos/hashicorp/terraform/releases?per_page=10",
		httpmock.NewStringResponder(200, readTestData("../../testdata/terraform-releases.json")),
	)

	got := captureStdout(t, func() {
		err := commandLatest(t.Context(), &commandLatestParams{
			githubToken:      "DUMMY",
			githubRepository: "hashicorp/terraform",
			cooldownDays:     0,
		})
		require.NoError(t, err)
	})

	assert.Equal(t, "v1.14.8\n", got)
}

func Test_commandHasPassed(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/repos/hashicorp/terraform/releases/tags/v1.14.8",
		httpmock.NewStringResponder(200, readTestData("../../testdata/terraform-tags-v1.14.8.json")),
	)

	tests := []struct {
		name       string
		params     *commandHasPassedParams
		wantErr    bool
		wantStdout string
	}{
		{
			name: "cooldown has passed without --exit-code",
			params: &commandHasPassedParams{
				githubToken:      "DUMMY",
				githubRepository: "hashicorp/terraform",
				githubTagName:    "v1.14.8",
				cooldownDays:     0,
				currentTime:      new(time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)),
				isExitCode:       false,
			},
			wantStdout: "Cooldown has passed.\n",
		},
		{
			name: "cooldown has not passed without --exit-code",
			params: &commandHasPassedParams{
				githubToken:      "DUMMY",
				githubRepository: "hashicorp/terraform",
				githubTagName:    "v1.14.8",
				cooldownDays:     7,
				currentTime:      new(time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)),
				isExitCode:       false,
			},
			wantStdout: "Cooldown has not passed yet.\n",
		},
		{
			name: "cooldown has passed with --exit-code",
			params: &commandHasPassedParams{
				githubToken:      "DUMMY",
				githubRepository: "hashicorp/terraform",
				githubTagName:    "v1.14.8",
				cooldownDays:     0,
				currentTime:      new(time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)),
				isExitCode:       true,
			},
		},
		{
			name: "cooldown has not passed with --exit-code",
			params: &commandHasPassedParams{
				githubToken:      "DUMMY",
				githubRepository: "hashicorp/terraform",
				githubTagName:    "v1.14.8",
				cooldownDays:     7,
				currentTime:      new(time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)),
				isExitCode:       true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStdout := captureStdout(t, func() {
				gotErr := commandHasPassed(t.Context(), tt.params)

				if tt.wantErr {
					assert.ErrorIs(t, gotErr, errExitCooldownNotPassed)
				} else {
					assert.NoError(t, gotErr)
				}
			})

			assert.Equal(t, tt.wantStdout, gotStdout)
		})
	}
}
