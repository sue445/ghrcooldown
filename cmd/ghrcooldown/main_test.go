package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	main "github.com/sue445/ghrcooldown/cmd/ghrcooldown"
)

func TestParseRepositoryPath_Success(t *testing.T) {
	tests := []struct {
		path string
		want *main.RepositoryPath
	}{
		{
			path: "user/repo",
			want: &main.RepositoryPath{
				Owner: "user",
				Repo:  "repo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, err := main.ParseRepositoryPath(tt.path)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestParseRepositoryPath_Error(t *testing.T) {
	tests := []struct {
		path string
	}{
		{
			path: "user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			_, err := main.ParseRepositoryPath(tt.path)
			assert.Error(t, err)
		})
	}
}
