package patch

import (
	"fmt"

	"github.com/cfgdiff/internal/diff"
)

// Direction controls whether a patch is applied forward or in reverse.
type Direction int

const (
	Forward Direction = iota
	Reverse
)

// Patch represents a set of changes that can be applied to a config map.
type Patch struct {
	Changes []diff.Change
	Direction Direction
}

// New creates a Patch from a slice of diff changes.
func New(changes []diff.Change, dir Direction) *Patch {
	return &Patch{Changes: changes, Direction: dir}
}

// Apply applies the patch to the provided config map, returning a new map.
// It does not mutate the input.
func (p *Patch) Apply(cfg map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(cfg))
	for k, v := range cfg {
		out[k] = v
	}

	for _, c := range p.Changes {
		if err := p.applyChange(out, c); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (p *Patch) applyChange(cfg map[string]string, c diff.Change) error {
	switch p.Direction {
	case Forward:
		return applyForward(cfg, c)
	case Reverse:
		return applyReverse(cfg, c)
	default:
		return fmt.Errorf("unknown patch direction: %d", p.Direction)
	}
}

func applyForward(cfg map[string]string, c diff.Change) error {
	switch c.Type {
	case diff.Added:
		cfg[c.Key] = c.NewValue
	case diff.Removed:
		delete(cfg, c.Key)
	case diff.Modified:
		if _, ok := cfg[c.Key]; !ok {
			return fmt.Errorf("patch: key %q not found for modification", c.Key)
		}
		cfg[c.Key] = c.NewValue
	}
	return nil
}

func applyReverse(cfg map[string]string, c diff.Change) error {
	switch c.Type {
	case diff.Added:
		// Reverse of add is remove
		delete(cfg, c.Key)
	case diff.Removed:
		// Reverse of remove is add back
		cfg[c.Key] = c.OldValue
	case diff.Modified:
		if _, ok := cfg[c.Key]; !ok {
			return fmt.Errorf("patch: key %q not found for reverse modification", c.Key)
		}
		cfg[c.Key] = c.OldValue
	}
	return nil
}
