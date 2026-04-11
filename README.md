# ghrcooldown
A Go library to fetch the latest GitHub release respecting the cooldown period.

[![test](https://github.com/sue445/ghrcooldown/actions/workflows/test.yml/badge.svg)](https://github.com/sue445/ghrcooldown/actions/workflows/test.yml)

## Example

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
