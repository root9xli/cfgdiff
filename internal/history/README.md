# history

The `history` package provides persistent storage and retrieval of diff history entries for `cfgdiff`.

## Overview

Every time a diff is run, an entry can be recorded to disk. Each entry captures:

- A unique ID (nanosecond Unix timestamp)
- The timestamp of the diff
- The paths of the two files compared
- The full list of changes
- An aggregated summary (added / removed / modified counts)

## Usage

```go
store, err := history.NewStore(".cfgdiff/history")
if err != nil { ... }

// Record a diff
entry, err := store.Record("prod.env", "staging.env", changes)

// List all entries (newest first)
entries, err := store.List()

// Retrieve a specific entry
entry, err := store.Get(id)
```

## Printing

```go
// Tabular list
history.PrintList(os.Stdout, entries)

// Full detail of one entry
history.PrintEntry(os.Stdout, entry)
```

## Storage

Entries are stored as individual JSON files named `<id>.json` inside the configured directory.
