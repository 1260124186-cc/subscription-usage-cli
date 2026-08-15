package main

import (
	"bytes"
	"path/filepath"
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

func TestRunStopsWritingAtTimeout(t *testing.T) {
	inputPath := filepath.Join("..", "..", "examples", "sample.json")

	if err := run(inputPath, time.Millisecond, slowWriter{}); err == nil {
		t.Fatal("run() error = nil, want context deadline exceeded")
	}
}

type slowWriter struct{}

func (slowWriter) Write(p []byte) (int, error) {
	time.Sleep(5 * time.Millisecond)
	return len(p), nil
}
