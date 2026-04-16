package main_test

import (
	"os"
)

// readTestData returns testdata
func readTestData(filename string) string {
	buf, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return string(buf)
}
