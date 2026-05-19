# logpipe

A structured log aggregator that tails multiple services, normalizes JSON log formats, and streams filtered output to stdout or file.

---

## Installation

```bash
go install github.com/yourname/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logpipe.git && cd logpipe && go build -o logpipe .
```

---

## Usage

Define your services in a `logpipe.yaml` config file:

```yaml
services:
  - name: api
    path: /var/log/api/app.log
  - name: worker
    path: /var/log/worker/app.log

filter:
  level: error
  fields:
    - timestamp
    - level
    - message

output:
  file: /tmp/aggregated.log
```

Then run:

```bash
logpipe --config logpipe.yaml
```

Stream filtered output directly to stdout:

```bash
logpipe --config logpipe.yaml --stdout
```

**Example output:**

```json
{"timestamp":"2024-01-15T10:23:01Z","service":"api","level":"error","message":"connection timeout"}
{"timestamp":"2024-01-15T10:23:04Z","service":"worker","level":"error","message":"job failed after 3 retries"}
```

---

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to config file | `logpipe.yaml` |
| `--stdout` | Stream output to stdout | `false` |
| `--level` | Minimum log level filter | `info` |

---

## License

MIT © yourname