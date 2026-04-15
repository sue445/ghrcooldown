package main_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// readTestData returns testdata
func readTestData(filename string) string {
	buf, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return string(buf)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	orgStdout := os.Stdout
	defer func() {
		os.Stdout = orgStdout
	}()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = w

	fn()

	w.Close()

	os.Stdout = orgStdout

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	return buf.String()
}
