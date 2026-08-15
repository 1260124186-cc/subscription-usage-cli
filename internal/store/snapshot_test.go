package store

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestLoadSnapshotFileReturnsCloseError(t *testing.T) {
	closeErr := errors.New("input close failed")
	input := &closeErrorReader{
		Reader:   strings.NewReader(`{"accounts":[],"usage":[]}`),
		closeErr: closeErr,
	}

	if _, err := LoadSnapshotFile(input); !errors.Is(err, closeErr) {
		t.Fatalf("LoadSnapshotFile() error = %v, want the input close error", err)
	}
}

func TestLoadSnapshotFileReturnsReadAndCloseErrors(t *testing.T) {
	readErr := errors.New("input read failed")
	closeErr := errors.New("input close failed")
	input := &closeErrorReader{
		Reader:   errorReader{err: readErr},
		closeErr: closeErr,
	}

	if _, err := LoadSnapshotFile(input); !errors.Is(err, readErr) || !errors.Is(err, closeErr) {
		t.Fatalf("LoadSnapshotFile() error = %v, want both read and close errors", err)
	}
}

type closeErrorReader struct {
	io.Reader
	closeErr error
}

func (reader closeErrorReader) Close() error {
	return reader.closeErr
}

type errorReader struct {
	err error
}

func (reader errorReader) Read([]byte) (int, error) {
	return 0, reader.err
}
