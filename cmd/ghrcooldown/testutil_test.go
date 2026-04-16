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

	defer func() {
		err := r.Close()
		require.NoError(t, err)
	}()

	var buf bytes.Buffer
	copyErrCh := make(chan error, 1)
	go func() {
		_, err := io.Copy(&buf, r)
		copyErrCh <- err
	}()

	os.Stdout = w
	defer func() {
		os.Stdout = orgStdout
		err := w.Close()
		require.NoError(t, err)
	}()

	fn()

	os.Stdout = orgStdout
	err = w.Close()
	require.NoError(t, err)

	err = <-copyErrCh
	require.NoError(t, err)
	return buf.String()
}
