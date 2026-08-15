package store

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestLoadSnapshotFileReturnsCloseError(t *testing.T) {
	input := &closeErrorReader{Reader: strings.NewReader(`{"accounts":[],"usage":[]}`)}

	if _, err := LoadSnapshotFile(input); err == nil {
		t.Fatal("LoadSnapshotFile() error = nil, want the input close error")
	}
}

type closeErrorReader struct {
	io.Reader
}

func (closeErrorReader) Close() error {
	return errors.New("input close failed")
}
