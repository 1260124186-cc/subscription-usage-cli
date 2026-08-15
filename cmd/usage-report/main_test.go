package main

import (
	"bytes"
	"errors"
	"io"
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

func TestRunWithInputClosesInputOnce(t *testing.T) {
	input := &countedInput{Reader: bytes.NewBufferString(`{"accounts":[],"usage":[]}`)}

	if err := runWithInput(input, time.Second, &bytes.Buffer{}); err != nil {
		t.Fatalf("runWithInput() error = %v", err)
	}
	if got, want := input.closes, 1; got != want {
		t.Fatalf("input Close() calls = %d, want %d", got, want)
	}
}

func TestRunWithInputReturnsCloseError(t *testing.T) {
	closeErr := errors.New("input close failed")
	input := &countedInput{
		Reader:   bytes.NewBufferString(`{"accounts":[],"usage":[]}`),
		closeErr: closeErr,
	}

	if err := runWithInput(input, time.Second, &bytes.Buffer{}); !errors.Is(err, closeErr) {
		t.Fatalf("runWithInput() error = %v, want the input close error", err)
	}
	if got, want := input.closes, 1; got != want {
		t.Fatalf("input Close() calls = %d, want %d", got, want)
	}
}

type countedInput struct {
	io.Reader
	closes   int
	closeErr error
}

func (input *countedInput) Close() error {
	input.closes++
	return input.closeErr
}
