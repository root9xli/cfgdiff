# summarize

The `summarize` package computes and formats aggregate statistics from a set of
config diff changes.

## Usage

```go
import (
    "os"
    "github.com/user/cfgdiff/internal/diff"
    "github.com/user/cfgdiff/internal/summarize"
)

changes := []diff.Change{ /* ... */ }

opts := summarize.DefaultOptions() // TopN: 5
stats := summarize.Compute(changes, opts)

// One-line string
fmt.Println(summarize.Format(stats))
// total=3 added=1 removed=1 modified=1

// Formatted block
summarize.Print(os.Stdout, stats)
// === Change Summary ===
//   Total    : 3
//   Added    : 1
//   Removed  : 1
//   Modified : 1
//   Top keys : port, host
```

## Options

| Field | Default | Description |
|-------|---------|-------------|
| `TopN` | `5` | Number of most-frequently changed keys to include in the report |

## Stats fields

| Field | Description |
|-------|-------------|
| `Added` | Keys present in new config but not in base |
| `Removed` | Keys present in base but not in new config |
| `Modified` | Keys whose value changed between configs |
| `Total` | Sum of all change types |
| `TopKeys` | Up to `TopN` keys ordered by change frequency |
