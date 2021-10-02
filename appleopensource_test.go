package appleopensource_test

import (
	"os"
	"testing"
)

func readTestdata(tb testing.TB, path string) []byte {
	tb.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		tb.Fatalf("not found %s: %v", path, err)
	}

	return data
}
