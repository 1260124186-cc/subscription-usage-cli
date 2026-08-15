package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/1260124186-cc/subscription-usage-cli/internal/output"
	"github.com/1260124186-cc/subscription-usage-cli/internal/service"
	"github.com/1260124186-cc/subscription-usage-cli/internal/store"
)

func main() {
	inputPath := flag.String("input", "", "path to a usage snapshot JSON file")
	timeout := flag.Duration("timeout", 5*time.Second, "maximum report generation duration")
	flag.Parse()

	if *inputPath == "" {
		exitf("missing required -input")
	}

	if err := run(*inputPath, *timeout, os.Stdout); err != nil {
		exitf("%v", err)
	}
}

func run(inputPath string, timeout time.Duration, stdout io.Writer) (err error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	return runWithInput(file, timeout, stdout)
}

func runWithInput(input io.ReadCloser, timeout time.Duration, stdout io.Writer) (err error) {
	snapshot, err := store.LoadSnapshotFile(input)
	if err != nil {
		return fmt.Errorf("load input: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	report, err := service.NewReportGenerator().Generate(ctx, snapshot)
	if err != nil {
		return fmt.Errorf("generate report: %w", err)
	}

	if err := output.WriteText(stdout, report); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	return nil
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
