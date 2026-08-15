package main

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunWritesUsageReport(t *testing.T) {
	var output bytes.Buffer
	inputPath := filepath.Join("..", "..", "examples", "sample.json")

	if err := run(inputPath, time.Second, &output); err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if got, want := output.String(), "account=acme plan=starter used=125 included=100 overage=25 charge_cents=175\naccount=beta plan=team used=4 included=10 overage=0 charge_cents=0\ntotal_charge_cents=175\n"; got != want {
		t.Fatalf("run() output = %q, want %q", got, want)
	}
}

func TestRunReturnsOutputWriteError(t *testing.T) {
	inputPath := filepath.Join("..", "..", "examples", "sample.json")
	writeErr := errors.New("report destination unavailable")

	err := run(inputPath, time.Second, failWriter{err: writeErr})
	if !errors.Is(err, writeErr) {
		t.Fatalf("run() error = %v, want %v", err, writeErr)
	}
	if !strings.Contains(err.Error(), "write report") {
		t.Fatalf("run() error = %q, want write report context", err)
	}
}

type failWriter struct {
	err error
}

func (w failWriter) Write([]byte) (int, error) {
	return 0, w.err
}
