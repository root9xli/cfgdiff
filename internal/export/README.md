# export

The `export` package provides functionality to write diff results to various output formats.

## Supported Formats

| Format     | Flag value  | Description                          |
|------------|-------------|--------------------------------------|
| CSV        | `csv`       | Comma-separated values, easy to import into spreadsheets |
| Markdown   | `markdown`  | GitHub-flavoured markdown table      |
| JSON       | `json`      | Structured JSON array of changes     |

## Usage

```go
import "github.com/user/cfgdiff/internal/export"

exporter, err := export.New(export.FormatCSV, os.Stdout)
if err != nil {
    log.Fatal(err)
}

if err := exporter.Write(changes); err != nil {
    log.Fatal(err)
}
```

## Writing to a file

To write output directly to a file, open the file and pass it as the writer:

```go
f, err := os.Create("report.csv")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

exporter, err := export.New(export.FormatCSV, f)
if err != nil {
    log.Fatal(err)
}

if err := exporter.Write(changes); err != nil {
    log.Fatal(err)
}
```

## CLI

Pass `--export` and `--export-format` flags to the root command:

```
cfgdiff diff base.yaml head.yaml --export report.csv --export-format csv
cfgdiff diff base.yaml head.yaml --export report.md  --export-format markdown
```

## Output example (CSV)

```
key,type,old_value,new_value
db.host,added,,localhost
db.port,modified,5432,5433
```
