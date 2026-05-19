package filter

import (
	"fmt"
	"strings"
)

// Rule defines a single filter criterion applied to normalized log fields.
type Rule struct {
	Field  string `yaml:"field"`
	Match  string `yaml:"match"`
	Invert bool   `yaml:"invert"`
}

// Filter holds a set of rules and applies them to log entries.
type Filter struct {
	rules []Rule
}

// New creates a Filter from the provided rules. An empty rule set
// means every entry passes.
func New(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// Match reports whether the given log entry (a map of normalized fields)
// satisfies all filter rules.
func (f *Filter) Match(entry map[string]interface{}) bool {
	for _, r := range f.rules {
		val, ok := entry[r.Field]
		if !ok {
			if !r.Invert {
				return false
			}
			continue
		}

		strVal := strings.ToLower(fmt.Sprintf("%v", val))
		matched := strings.Contains(strVal, strings.ToLower(r.Match))

		if r.Invert {
			matched = !matched
		}
		if !matched {
			return false
		}
	}
	return true
}
