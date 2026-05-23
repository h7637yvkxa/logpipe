package metrics

import (
	"testing"
)

func TestNew_ZeroValues(t *testing.T) {
	c := New()
	snap := c.Snapshot()
	for key, val := range snap {
		if val != 0 {
			t.Errorf("expected 0 for %s, got %d", key, val)
		}
	}
}

func TestIncRead_GlobalAndPerSource(t *testing.T) {
	c := New()
	c.IncRead("svc-a")
	c.IncRead("svc-a")
	c.IncRead("svc-b")

	if got := c.LinesRead.Load(); got != 3 {
		t.Errorf("LinesRead: want 3, got %d", got)
	}

	snap := c.Snapshot()
	if snap["source.svc-a"] != 2 {
		t.Errorf("source.svc-a: want 2, got %d", snap["source.svc-a"])
	}
	if snap["source.svc-b"] != 1 {
		t.Errorf("source.svc-b: want 1, got %d", snap["source.svc-b"])
	}
}

func TestIncFiltered(t *testing.T) {
	c := New()
	c.IncFiltered()
	c.IncFiltered()
	if got := c.LinesFiltered.Load(); got != 2 {
		t.Errorf("LinesFiltered: want 2, got %d", got)
	}
}

func TestIncWritten(t *testing.T) {
	c := New()
	c.IncWritten()
	if got := c.LinesWritten.Load(); got != 1 {
		t.Errorf("LinesWritten: want 1, got %d", got)
	}
}

func TestIncParseError(t *testing.T) {
	c := New()
	c.IncParseError()
	c.IncParseError()
	c.IncParseError()
	if got := c.ParseErrors.Load(); got != 3 {
		t.Errorf("ParseErrors: want 3, got %d", got)
	}
}

func TestSnapshot_ContainsAllKeys(t *testing.T) {
	c := New()
	c.IncRead("api")
	c.IncFiltered()
	c.IncWritten()
	c.IncParseError()

	snap := c.Snapshot()
	required := []string{"lines_read", "lines_filtered", "lines_written", "parse_errors", "source.api"}
	for _, k := range required {
		if _, ok := snap[k]; !ok {
			t.Errorf("snapshot missing key %q", k)
		}
	}
}

func TestSourceCounter_ConcurrentSafety(t *testing.T) {
	c := New()
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			c.IncRead("shared-source")
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
	if got := c.LinesRead.Load(); got != 50 {
		t.Errorf("concurrent LinesRead: want 50, got %d", got)
	}
}
