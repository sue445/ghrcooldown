package ghrcooldown_test

import (
	"os"
	"time"
)

// readTestData returns testdata
func readTestData(filename string) string {
	buf, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return string(buf)
}

func days(days int) time.Duration {
	return time.Duration(days) * 24 * time.Hour
}
