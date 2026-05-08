package summarize

import (
	"fmt"
	"sort"

	"github.com/user/cfgdiff/internal/diff"
)

// Stats holds aggregated statistics about a set of config changes.
type Stats struct {
	Added    int
	Removed  int
	Modified int
	Total    int
	TopKeys  []string
}

// Options controls summarize behaviour.
type Options struct {
	// TopN limits how many frequently changed keys are reported.
	TopN int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{TopN: 5}
}

// Compute derives a Stats summary from a slice of diff.Change values.
func Compute(changes []diff.Change, opts Options) Stats {
	if opts.TopN <= 0 {
		opts.TopN = DefaultOptions().TopN
	}

	freq := make(map[string]int)
	s := Stats{}

	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			s.Added++
		case diff.Removed:
			s.Removed++
		case diff.Modified:
			s.Modified++
		}
		freq[c.Key]++
	}

	s.Total = s.Added + s.Removed + s.Modified
	s.TopKeys = topN(freq, opts.TopN)
	return s
}

// Format returns a human-readable one-line summary string.
func Format(s Stats) string {
	return fmt.Sprintf("total=%d added=%d removed=%d modified=%d",
		s.Total, s.Added, s.Removed, s.Modified)
}

// topN returns the top-n keys by frequency, sorted deterministically.
func topN(freq map[string]int, n int) []string {
	type kv struct {
		key   string
		count int
	}
	var pairs []kv
	for k, v := range freq {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].key < pairs[j].key
	})
	var keys []string
	for i, p := range pairs {
		if i >= n {
			break
		}
		keys = append(keys, p.key)
	}
	return keys
}
