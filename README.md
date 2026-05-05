# logslice

A fast log file slicer and filter tool that supports structured and unstructured log formats.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

---

## Usage

```bash
# Slice logs between two timestamps
logslice --from "2024-01-15T10:00:00" --to "2024-01-15T11:00:00" app.log

# Filter by log level
logslice --level ERROR app.log

# Filter by keyword and output to file
logslice --grep "connection refused" --output errors.log app.log

# Parse structured JSON logs
logslice --format json --level WARN service.log
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start timestamp for slicing |
| `--to` | End timestamp for slicing |
| `--level` | Filter by log level (INFO, WARN, ERROR) |
| `--grep` | Filter lines matching a pattern |
| `--format` | Log format: `text` (default) or `json` |
| `--output` | Write results to a file instead of stdout |

---

## Features

- Fast line-by-line streaming — handles large log files with low memory usage
- Supports structured (JSON) and unstructured (plain text) log formats
- Flexible timestamp parsing with automatic format detection
- Composable filters for level, pattern, and time range

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)