# TSDB - High-Performance Time-Series Database in Go

A lightweight, high-performance time-series database written in Go, inspired by TDengine. Designed for IoT, monitoring, and real-time analytics workloads.

## Features

- **Columnar Storage**: Optimized for time-series data with columnar layout
- **High-Performance Ingestion**: Fast data writes with in-memory buffering
- **Compression**: Built-in zstd compression for efficient storage
- **RESTful API**: Simple HTTP API for data operations
- **Concurrent Access**: Thread-safe operations with fine-grained locking
- **Cross-Platform**: Runs on Linux, macOS, and Windows

## Quick Start

### Build from Source

```bash
go build -o tsdb ./cmd/tsdb
```

### Run Server

```bash
./tsdb -port 6041 -data ./data
```

### API Usage

#### Create Table

```bash
curl -X POST http://localhost:6041/api/create \
  -H "Content-Type: application/json" \
  -d '{
    "table": "sensors",
    "columns": ["ts", "temperature", "humidity", "device_id"],
    "types": ["timestamp", "float", "float", "string"]
  }'
```

#### Insert Data

```bash
curl -X POST http://localhost:6041/api/insert \
  -H "Content-Type: application/json" \
  -d '{
    "table": "sensors",
    "values": [
      {
        "ts": "2026-03-17T10:00:00Z",
        "temperature": 23.5,
        "humidity": 65.2,
        "device_id": "sensor001"
      }
    ]
  }'
```

#### Query Data

```bash
curl -X POST http://localhost:6041/api/query \
  -H "Content-Type: application/json" \
  -d '{
    "table": "sensors",
    "start_time": "2026-03-17T00:00:00Z",
    "end_time": "2026-03-17T23:59:59Z",
    "columns": ["ts", "temperature"]
  }'
```

## Architecture

```
┌─────────────────────────────────────┐
│         HTTP API Layer              │
│  (RESTful endpoints for CRUD)       │
└─────────────────────────────────────┘
                 │
┌─────────────────────────────────────┐
│       Storage Engine                │
│  - Columnar storage                 │
│  - In-memory buffering              │
│  - Compression (zstd)               │
│  - Concurrent access control        │
└─────────────────────────────────────┘
                 │
┌─────────────────────────────────────┐
│         Disk Persistence            │
│  (Binary format with compression)   │
└─────────────────────────────────────┘
```

## Performance Optimizations

- **Columnar Layout**: Data stored by column for better compression and query performance
- **Zero-Copy Reads**: Efficient memory access patterns
- **Fast Compression**: zstd compression with fastest preset for minimal CPU overhead
- **Lock-Free Reads**: Read operations use RWMutex for concurrent access
- **Batch Writes**: Support for bulk inserts

## Benchmarks

Run benchmarks:

```bash
go test -bench=. -benchmem ./internal/storage
```

## Development

### Project Structure

```
.
├── cmd/
│   └── tsdb/           # Main server entry point
├── internal/
│   ├── api/            # HTTP API handlers
│   └── storage/        # Storage engine
├── .github/
│   └── workflows/      # CI/CD pipelines
└── go.mod
```

### Testing

```bash
go test ./...
```

## Roadmap

- [ ] SQL query parser
- [ ] Distributed clustering
- [ ] Stream processing
- [ ] Data replication
- [ ] Advanced compression algorithms
- [ ] Grafana integration

## License

MIT License

## Contributing

Contributions welcome! Please open an issue or submit a PR.
