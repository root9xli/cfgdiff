# cfgdiff

> CLI tool to diff and audit config file changes across environments

---

## Installation

```bash
go install github.com/yourusername/cfgdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/cfgdiff.git && cd cfgdiff && go build -o cfgdiff .
```

---

## Usage

Compare config files across two environments:

```bash
cfgdiff --base config.staging.yaml --target config.production.yaml
```

Audit changes and output a structured diff report:

```bash
cfgdiff --base config.staging.yaml --target config.production.yaml --format json --output report.json
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--base` | Base config file path | required |
| `--target` | Target config file path | required |
| `--format` | Output format (`text`, `json`, `yaml`) | `text` |
| `--output` | Write output to file instead of stdout | — |
| `--ignore` | Comma-separated keys to ignore | — |

---

## Supported Formats

- YAML
- JSON
- TOML
- `.env`

---

## License

[MIT](LICENSE) © 2024 yourusername