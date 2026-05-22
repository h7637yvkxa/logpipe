// Package main is the entry point for the logpipe log aggregation tool.
// It wires together configuration loading, tailing, normalization, filtering,
// and output writing into a single streaming pipeline.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/logpipe/internal/config"
	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/normalize"
	"github.com/yourorg/logpipe/internal/output"
	"github.com/yourorg/logpipe/internal/tail"
)

func main() {
	cfgPath := flag.String("config", "logpipe.yaml", "path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logpipe: failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Build the tail manager across all configured sources.
	manager := tail.NewManager(cfg.Sources)

	// Build the normalization pipeline.
	normPipeline := normalize.NewPipeline(cfg.Sources)

	// Build the filter pipeline from configured rules.
	filterPipeline := filter.Pipeline(cfg.Filters)

	// Build the output writer.
	writer, err := output.New(cfg.Output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logpipe: failed to create output writer: %v\n", err)
		os.Exit(1)
	}
	defer writer.Close()

	// Handle OS signals for graceful shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	lines := manager.Lines()

	for {
		select {
		case line, ok := <-lines:
			if !ok {
				// Channel closed; all tailers have finished.
				return
			}

			// Normalize the raw log line into a structured entry.
			entry, err := normPipeline.Process(line)
			if err != nil {
				// Skip lines that cannot be parsed as JSON.
				continue
			}

			// Apply filter rules; skip entries that do not match.
			if !filterPipeline.Match(entry) {
				continue
			}

			// Write the normalized, filtered entry to the configured output.
			if err := writer.Write(entry); err != nil {
				fmt.Fprintf(os.Stderr, "logpipe: write error: %v\n", err)
			}

		case sig := <-sigCh:
			fmt.Fprintf(os.Stderr, "logpipe: received signal %s, shutting down\n", sig)
			return
		}
	}
}
