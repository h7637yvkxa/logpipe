package normalize

import (
	"encoding/json"
	"fmt"
	"time"
)

// Entry represents a normalized log entry.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Service   string            `json:"service"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Normalizer converts raw JSON log lines into a unified Entry format.
type Normalizer struct {
	service string
}

// New creates a Normalizer associated with the given service name.
func New(service string) *Normalizer {
	return &Normalizer{service: service}
}

// Normalize parses a raw JSON log line and returns a normalized Entry.
// It handles common field name variants (e.g. "msg"/"message", "ts"/"time"/"timestamp").
func (n *Normalizer) Normalize(line string) (*Entry, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("normalize: invalid JSON: %w", err)
	}

	entry := &Entry{
		Service: n.service,
		Fields:  make(map[string]interface{}),
	}

	// Resolve timestamp
	for _, key := range []string{"timestamp", "time", "ts"} {
		if v, ok := raw[key]; ok {
			switch val := v.(type) {
			case string:
				if t, err := time.Parse(time.RFC3339, val); err == nil {
					entry.Timestamp = t
				}
			case float64:
				entry.Timestamp = time.Unix(int64(val), 0).UTC()
			}
			delete(raw, key)
			break
		}
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	// Resolve level
	for _, key := range []string{"level", "severity", "lvl"} {
		if v, ok := raw[key]; ok {
			entry.Level = fmt.Sprintf("%v", v)
			delete(raw, key)
			break
		}
	}
	if entry.Level == "" {
		entry.Level = "info"
	}

	// Resolve message
	for _, key := range []string{"message", "msg"} {
		if v, ok := raw[key]; ok {
			entry.Message = fmt.Sprintf("%v", v)
			delete(raw, key)
			break
		}
	}

	// Remaining keys become extra fields
	for k, v := range raw {
		entry.Fields[k] = v
	}
	if len(entry.Fields) == 0 {
		entry.Fields = nil
	}

	return entry, nil
}
