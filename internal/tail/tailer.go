package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Line represents a single log line from a source.
type Line struct {
	Source string
	Text   string
}

// Tailer tails a file and emits lines to a channel.
type Tailer struct {
	path   string
	source string
	out    chan<- Line
}

// New creates a new Tailer for the given file path and source label.
func New(path, source string, out chan<- Line) *Tailer {
	return &Tailer{
		path:   path,
		source: source,
		out:    out,
	}
}

// Run opens the file, seeks to the end, and streams new lines until ctx is cancelled.
func (t *Tailer) Run(ctx context.Context) error {
	f, err := os.Open(t.path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	reader := bufio.NewReader(f)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(200 * time.Millisecond)
				continue
			}
			return err
		}

		if len(line) > 0 {
			t.out <- Line{Source: t.source, Text: line}
		}
	}
}
