package main

import (
	"context"
	"flag"
	"fmt"
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

	file, err := os.Open(*inputPath)
	if err != nil {
		exitf("open input: %v", err)
	}
	defer file.Close()

	snapshot, err := store.LoadSnapshot(file)
	if err != nil {
		exitf("load input: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	report, err := service.NewReportGenerator().Generate(ctx, snapshot)
	if err != nil {
		exitf("generate report: %v", err)
	}

	if err := output.WriteText(os.Stdout, report); err != nil {
		exitf("write report: %v", err)
	}
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
