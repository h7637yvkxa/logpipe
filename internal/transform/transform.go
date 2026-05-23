package transform

import (
	"strings"

	"github.com/user/logpipe/internal/normalize"
)

// Rule defines a single field transformation.
type Rule struct {
	// Field is the log entry field to transform.
	Field string
	// Op is the operation: "uppercase", "lowercase", "truncate", "replace".
	Op string
	// Arg is an optional argument (e.g. max length for truncate, or "old:new" for replace).
	Arg string
}

// Transformer applies a set of Rules to log entries.
type Transformer struct {
	rules []Rule
}

// New creates a Transformer with the given rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply returns a new entry with all rules applied.
// Fields not targeted by any rule are left unchanged.
func (t *Transformer) Apply(entry normalize.Entry) normalize.Entry {
	out := normalize.Entry{
		Timestamp: entry.Timestamp,
		Level:     entry.Level,
		Message:   entry.Message,
		Source:    entry.Source,
		Extra:     make(map[string]any, len(entry.Extra)),
	}
	for k, v := range entry.Extra {
		out.Extra[k] = v
	}

	for _, r := range t.rules {
		switch r.Field {
		case "level":
			out.Level = applyOp(out.Level, r)
		case "message":
			out.Message = applyOp(out.Message, r)
		case "source":
			out.Source = applyOp(out.Source, r)
		default:
			if v, ok := out.Extra[r.Field]; ok {
				if s, ok2 := v.(string); ok2 {
					out.Extra[r.Field] = applyOp(s, r)
				}
			}
		}
	}
	return out
}

func applyOp(s string, r Rule) string {
	switch r.Op {
	case "uppercase":
		return strings.ToUpper(s)
	case "lowercase":
		return strings.ToLower(s)
	case "truncate":
		max := 0
		if _, err := parseIntArg(r.Arg, &max); err == nil && max > 0 && len(s) > max {
			return s[:max]
		}
		return s
	case "replace":
		parts := strings.SplitN(r.Arg, ":", 2)
		if len(parts) == 2 {
			return strings.ReplaceAll(s, parts[0], parts[1])
		}
		return s
	default:
		return s
	}
}

func parseIntArg(s string, out *int) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, &parseErr{}
		}
		n = n*10 + int(c-'0')
	}
	*out = n
	return n, nil
}

type parseErr struct{}

func (e *parseErr) Error() string { return "invalid integer" }
