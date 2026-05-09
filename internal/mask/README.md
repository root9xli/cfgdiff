# mask

The `mask` package provides value-masking for sensitive config keys.

It scans config maps and replaces the values of keys that match
configurable patterns (e.g. `password`, `secret`, `token`) with a
masking character, keeping the key names intact.

## Usage

```go
import "github.com/yourorg/cfgdiff/internal/mask"

m := mask.New(mask.DefaultOptions())

// Check a single key
if m.IsSensitive("db_password") { ... }

// Mask a whole config map
masked := m.Apply(configMap)
```

## Options

| Field | Default | Description |
|-------|---------|-------------|
| `Char` | `"*"` | Replacement character |
| `ShowFirst` | `0` | Leading characters to reveal |
| `ShowLast` | `0` | Trailing characters to reveal |

```go
opts := mask.Options{
    Char:      "#",
    ShowFirst: 2,
    ShowLast:  2,
}
m := mask.New(opts)
m.MaskValue("abcdefgh") // => "ab####gh"
```

## Custom Patterns

```go
m := mask.NewWithPatterns(mask.DefaultOptions(), []string{
    `(?i)internal`,
    `(?i)credential`,
})
```

## Default Sensitive Patterns

- `password`
- `secret`
- `token`
- `api_key` / `apikey`
- `private_key`
- `auth`
