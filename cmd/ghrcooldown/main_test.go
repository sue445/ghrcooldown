package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseRepositoryPath_Success(t *testing.T) {
	tests := []struct {
		path string
		want *repositoryPath
	}{
		{
			path: "user/repo",
			want: &repositoryPath{
				Owner: "user",
				Repo:  "repo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, err := parseRepositoryPath(tt.path)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_parseRepositoryPath_Error(t *testing.T) {
	tests := []struct {
		path string
	}{
		{
			path: "user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			_, err := parseRepositoryPath(tt.path)
			assert.Error(t, err)
		})
	}
}
