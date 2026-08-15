package main

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/1260124186-cc/subscription-usage-cli/internal/output"
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
	writer := newSlowWriter(200 * time.Millisecond)

	started := time.Now()
	err := run(inputPath, time.Millisecond, writer)
	elapsed := time.Since(started)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("run() error = %v, want context deadline exceeded", err)
	}
	if elapsed >= 100*time.Millisecond {
		t.Fatalf("run() elapsed = %s, want less than 100ms", elapsed)
	}
	if got, want := writer.writes.Load(), int32(1); got != want {
		t.Fatalf("writer writes = %d, want %d", got, want)
	}
	select {
	case <-writer.exited:
	default:
		t.Fatal("write was still running after run() returned")
	}
	if writer.completed.Load() {
		t.Fatal("write completed after the context deadline")
	}
}

type slowWriter struct {
	delay     time.Duration
	writes    atomic.Int32
	completed atomic.Bool
	exited    chan struct{}
}

func newSlowWriter(delay time.Duration) *slowWriter {
	return &slowWriter{
		delay:  delay,
		exited: make(chan struct{}),
	}
}

func (w *slowWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (w *slowWriter) WriteContext(ctx context.Context, p []byte) (int, error) {
	w.writes.Add(1)
	defer close(w.exited)

	timer := time.NewTimer(w.delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-timer.C:
		w.completed.Store(true)
		return len(p), nil
	}
}

var _ output.ContextWriter = (*slowWriter)(nil)
