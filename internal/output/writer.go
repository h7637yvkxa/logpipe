package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/yourorg/logpipe/internal/config"
)

// Writer handles writing normalized log entries to stdout or a file.
type Writer struct {
	out    io.Writer
	closer io.Closer
}

// New creates a Writer based on the provided output config.
// If cfg.Type is "file", output is written to cfg.Path.
// Otherwise, output defaults to stdout.
func New(cfg config.OutputConfig) (*Writer, error) {
	if cfg.Type == "file" {
		if cfg.Path == "" {
			return nil, fmt.Errorf("output type 'file' requires a path")
		}
		f, err := os.OpenFile(cfg.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("opening output file %q: %w", cfg.Path, err)
		}
		return &Writer{out: f, closer: f}, nil
	}
	return &Writer{out: os.Stdout}, nil
}

// Write serializes the log entry map as a single JSON line.
func (w *Writer) Write(entry map[string]interface{}) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshalling log entry: %w", err)
	}
	_, err = fmt.Fprintf(w.out, "%s\n", data)
	return err
}

// Close releases any underlying file resource.
func (w *Writer) Close() error {
	if w.closer != nil {
		return w.closer.Close()
	}
	return nil
}
