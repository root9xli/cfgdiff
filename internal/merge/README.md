# merge

The `merge` package provides functionality to combine two flat config maps
(as produced by `internal/parser`) using a configurable conflict resolution
strategy.

## Usage

```go
import "cfgdiff/internal/merge"

result, err := merge.Merge(base, override, merge.StrategyOverride)
if err != nil {
    log.Fatal(err)
}

for _, c := range result.Conflicts {
    fmt.Printf("conflict: %s  base=%v  override=%v\n", c.Key, c.BaseVal, c.OverVal)
}
```

## Strategies

| Strategy | Behaviour |
|---|---|
| `StrategyBase` | Keep the base value when a conflict occurs. |
| `StrategyOverride` | Replace the base value with the override value. |
| `StrategyError` | Return an error immediately on the first conflict. |

## Conflict reporting

Regardless of strategy (except `StrategyError`), all conflicts are collected
in `Result.Conflicts` so callers can log or display them without stopping
the merge.

## Notes

- Both input maps must be **flat** (dot-separated keys). Use `parser.Parse`
  to obtain a flat map from a config file.
- Input maps are **not mutated**; the merged data is returned in a new map.
