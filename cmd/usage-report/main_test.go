package main

import (
	"bytes"
	"os"
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

func TestRunRejectsNullAccount(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "snapshot.json")
	if err := os.WriteFile(inputPath, []byte(`{"accounts":[null],"usage":[]}`), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var output bytes.Buffer
	err := run(inputPath, time.Second, &output)
	if err == nil {
		t.Fatal("run() error = nil, want an error for a null account")
	}
	if got, want := err.Error(), "load input: account at index 0 is null"; got != want {
		t.Fatalf("run() error = %q, want %q", got, want)
	}
	if got := output.String(); got != "" {
		t.Fatalf("run() output = %q, want no report", got)
	}
}
