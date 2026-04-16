# ghrcooldown
A CLI tool and Go library to fetch the latest GitHub Releases respecting the cooldown period.

[![Latest Version](https://img.shields.io/github/v/tag/sue445/ghrcooldown)](https://github.com/sue445/ghrcooldown/tags)
[![test](https://github.com/sue445/ghrcooldown/actions/workflows/test.yml/badge.svg)](https://github.com/sue445/ghrcooldown/actions/workflows/test.yml)
[![Maintainability](https://qlty.sh/badges/bf40eb29-6803-4e4f-a75e-10e3ae7ecb66/maintainability.svg)](https://qlty.sh/gh/sue445/projects/ghrcooldown)
[![Coverage Status](https://coveralls.io/repos/github/sue445/ghrcooldown/badge.svg)](https://coveralls.io/github/sue445/ghrcooldown)
[![GoDoc](https://godoc.org/github.com/sue445/ghrcooldown?status.svg)](https://godoc.org/github.com/sue445/ghrcooldown)
[![Go Report Card](https://goreportcard.com/badge/github.com/sue445/ghrcooldown)](https://goreportcard.com/report/github.com/sue445/ghrcooldown)

## Usage as a CLI tool
```bash
$ ghrcooldown latest --repo hashicorp/terraform --cooldown-days 7
v1.14.8
```

### Install
TBD

### Commands
#### `ghrcooldown latest`
```
$ ghrcooldown latest --help
NAME:
   ghrcooldown latest - Print latest release version of the specified repository, respecting the provided cooldown period.

USAGE:
   ghrcooldown latest [options]

OPTIONS:
   --cooldown-days int      Cooldown days (default: 0)
   --github-api-url string  GitHub API Endpoint (e.g. https://<your-ghes-hostname>/api/v3) [$GITHUB_API_URL]
   --repo string            GitHub Repository Path (e.g. user/repo)
   --token string           GitHub token [$GITHUB_TOKEN]
   --help, -h               show help
```

## Usage as a library

### Install
```bash
go get -u github.com/sue445/ghrcooldown
```

### Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sue445/ghrcooldown"
)

func main() {
	// Initialize the client
	client, err := ghrcooldown.NewClient(&ghrcooldown.ClientParams{
		Token: os.Getenv("GITHUB_TOKEN"),
		// BaseURL: "", // Required if using GitHub Enterprise Server
	})

	// 7 days
	cooldown := 7 * 24 * time.Hour

	ctx := context.Background()

	// Example 1: Returns the latest tagName that has passed the 7-day cooldown period.
	tagName, err := client.GetLatestTagName(ctx, "hashicorp", "terraform", cooldown)
	if err != nil {
		log.Fatalf("failed to get latest tag: %v", err)
	}
	fmt.Printf("Latest tag: %s\n", tagName)

	// Example 2: Returns true if the cooldown has passed.
	hasPassed, err := client.HasCooldownPassed(ctx, "hashicorp", "terraform", "v1.14.8", cooldown)
	if err != nil {
		log.Fatalf("failed to check cooldown: %v", err)
	}
	fmt.Printf("Has v1.14.8 passed the cooldown? %v\n", hasPassed)
}
```

## Reference
https://godoc.org/github.com/sue445/ghrcooldown
