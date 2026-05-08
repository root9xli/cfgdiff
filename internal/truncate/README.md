# truncate

The `truncate` package provides value-truncation utilities for the `cfgdiff` CLI.
When diffing configs with long values (e.g. base64 blobs, JWT tokens, certificate
PEM data), output can become unreadable. This package lets you cap value length
before rendering.

## Usage

```go
import "github.com/cfgdiff/cfgdiff/internal/truncate"

// Use defaults (80 chars, "..." suffix)
tr := truncate.New(truncate.DefaultOptions())

// Truncate a single value
short := tr.Value("a very long configuration value that goes on forever")

// Truncate all values in a parsed config map
truncated := tr.Apply(configMap)

// Report how many values were shortened
fmt.Println(tr.Summary(configMap, truncated))
```

## Options

| Field     | Type   | Default | Description                          |
|-----------|--------|---------|--------------------------------------|
| `MaxLen`  | int    | 80      | Maximum number of characters allowed |
| `Suffix`  | string | `...`   | Appended when a value is cut         |
| `Enabled` | bool   | true    | Set false to disable all truncation  |

## Notes

- Keys are **never** truncated — only values.
- `Apply` always returns a **new** map; the original is not mutated.
- Setting `MaxLen` to `0` falls back to the default (80).
