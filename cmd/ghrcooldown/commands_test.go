package main_test

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	main "github.com/sue445/ghrcooldown/cmd/ghrcooldown"
)

func TestCommandLatest(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/repos/hashicorp/terraform/releases?per_page=10",
		httpmock.NewStringResponder(200, readTestData("../../testdata/terraform-releases.json")),
	)

	got := captureStdout(t, func() {
		err := main.CommandLatest(t.Context(), &main.CommandLatestParams{
			GithubToken:      "DUMMY",
			GithubRepository: "hashicorp/terraform",
			CooldownDays:     0,
		})
		require.NoError(t, err)
	})

	assert.Equal(t, "v1.14.8\n", got)
}
