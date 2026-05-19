package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/tail"
)

func TestManager_AggregatesMultipleSources(t *testing.T) {
	pathA := writeTempLog(t, "")
	pathB := writeTempLog(t, "")

	sources := []tail.Source{
		{Name: "svc-a", Path: pathA},
		{Name: "svc-b", Path: pathB},
	}

	mgr := tail.NewManager(sources)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		_ = mgr.Run(ctx)
	}()

	time.Sleep(60 * time.Millisecond)

	for _, p := range []string{pathA, pathB} {
		f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			t.Fatalf("open %s: %v", p, err)
		}
		_, _ = f.WriteString(`{"msg":"test"}` + "\n")
		f.Close()
	}

	received := map[string]bool{}
	timeout := time.After(2 * time.Second)
	for len(received) < 2 {
		select {
		case line := <-mgr.Lines():
			received[line.Source] = true
		case <-timeout:
			t.Fatalf("timed out; only received sources: %v", received)
		}
	}

	if !received["svc-a"] || !received["svc-b"] {
		t.Errorf("did not receive lines from all sources: %v", received)
	}

	cancel()
}

func TestManager_Lines_ChannelNotNil(t *testing.T) {
	mgr := tail.NewManager(nil)
	if mgr.Lines() == nil {
		t.Error("expected non-nil Lines channel")
	}
}
