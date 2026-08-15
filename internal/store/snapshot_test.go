package store

import (
	"strings"
	"testing"
)

func TestLoadSnapshotRejectsNullAccount(t *testing.T) {
	input := `{"accounts":[null],"usage":[]}`

	if _, err := LoadSnapshot(strings.NewReader(input)); err == nil {
		t.Fatal("LoadSnapshot() error = nil, want an error for a null account")
	}
}
