# Redis-GO

A minimal Redis-like server implemented in Go, supporting basic RESP protocol commands and multiple connections.

## Features

- RESP protocol parsing and serialization
- Supports basic commands: `PING`, `SET`, `GET`, `HSET`, `HGET`
- In-memory key-value and hash storage
- Thread-safe operations

## Getting Started

### Prerequisites

- Go 1.18+
- (Optional) Docker for running a real Redis instance for comparison
- Docker run command: docker compose -f <path/to/your-compose-file.yml> up -d

### Running the Server

From the project directory, run:

```sh
go run *.go
