# classify

The `classify` package assigns a **severity** and **domain** to each config diff change, making it easy to prioritise review and audit work.

## Severity levels

| Level      | Trigger (default)                              |
|------------|------------------------------------------------|
| `critical` | Key contains: `secret`, `password`, `token`, `key`, `cert` |
| `high`     | Key contains: `host`, `port`, `url`, `endpoint`, `dsn`, `db` |
| `low`      | Everything else                                |

## Domain detection

Domain is extracted from the key prefix:

- `db.password` → domain `db` (dot-notation)
- `APP_NAME` → domain `app` (underscore prefix, lowercased)
- `name` → domain `general`

## Usage

```go
opts := classify.DefaultOptions()
c := classify.New(opts)
results := c.Apply(changes)
classify.PrintResults(results)
```

## Custom patterns

```go
opts := classify.Options{
    CriticalPatterns: []string{"private", "secret"},
    HighPatterns:     []string{"host", "port"},
}
c := classify.New(opts)
```
