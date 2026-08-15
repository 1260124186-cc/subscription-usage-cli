package store

import (
	"strings"
	"testing"
)

func TestLoadSnapshotRejectsNullAccount(t *testing.T) {
	input := `{"accounts":[null],"usage":[]}`

	_, err := LoadSnapshot(strings.NewReader(input))
	if err == nil {
		t.Fatal("LoadSnapshot() error = nil, want an error for a null account")
	}
	if got, want := err.Error(), "account at index 0 is null"; got != want {
		t.Fatalf("LoadSnapshot() error = %q, want %q", got, want)
	}
}
