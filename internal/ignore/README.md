# cfgdiff ignore rules

The `ignore` package provides pattern-based key filtering for `cfgdiff`.
It allows users to exclude sensitive or irrelevant config keys from diff output and audit logs.

## Usage

Create a `.cfgignore` file in your project root (or pass `--ignore` flag to the CLI).
Each line defines one pattern. Empty lines and lines starting with `#` are ignored.

### Pattern syntax

| Pattern | Description |
|---|---|
| `DB_PASSWORD` | Exact key match |
| `*_SECRET` | Shell glob — any key ending in `_SECRET` |
| `secret.*` | Prefix match — any nested key under `secret.` |

### Example `.cfgignore`

```
# Credentials
DB_PASSWORD
DB_SECRET_KEY

# Tokens
*_TOKEN
*_API_KEY

# Nested secrets
secret.*
```

### Programmatic use

```go
// Load from file
rules, err := ignore.LoadFile(".cfgignore")

// Or define inline
rules := ignore.NewRules([]string{"DB_PASSWORD", "*_TOKEN"})

// Check a single key
if rules.Match("DB_PASSWORD") {
    // skip
}

// Filter a slice of keys
visible := rules.FilterKeys(allKeys)
```
