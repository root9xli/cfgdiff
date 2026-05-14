package group

import (
	"sort"
	"strings"

	"github.com/cfgdiff/internal/diff"
)

// Options controls how changes are grouped.
type Options struct {
	// Separator is the delimiter used to split key prefixes (default: ".").
	Separator string
	// Depth is the number of prefix segments to group by (default: 1).
	Depth int
}

// DefaultOptions returns sensible defaults for grouping.
func DefaultOptions() Options {
	return Options{
		Separator: ".",
		Depth:     1,
	}
}

// Result holds changes that share a common key prefix.
type Result struct {
	Prefix  string
	Changes []diff.Change
}

// Apply groups a slice of changes by their key prefix.
// Keys without a separator are placed under the prefix "(root)".
func Apply(changes []diff.Change, opts Options) []Result {
	if opts.Separator == "" {
		opts.Separator = "."
	}
	if opts.Depth < 1 {
		opts.Depth = 1
	}

	buckets := make(map[string][]diff.Change)
	for _, c := range changes {
		prefix := extractPrefix(c.Key, opts.Separator, opts.Depth)
		buckets[prefix] = append(buckets[prefix], c)
	}

	results := make([]Result, 0, len(buckets))
	for prefix, cs := range buckets {
		results = append(results, Result{Prefix: prefix, Changes: cs})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Prefix < results[j].Prefix
	})
	return results
}

// extractPrefix returns the first `depth` segments of key joined by sep.
func extractPrefix(key, sep string, depth int) string {
	parts := strings.Split(key, sep)
	if len(parts) <= depth {
		return "(root)"
	}
	return strings.Join(parts[:depth], sep)
}
