# Adewale Go Chat

A concurrent terminal chat server written in Go. This repository is an independent implementation built to study TCP networking, goroutine-based concurrency, shared-state synchronization, room routing and abuse controls.

## Features

- Concurrent TCP clients
- Named chat rooms
- Room-scoped broadcasts
- Direct messages
- `/users`, `/rooms`, `/join`, `/msg`, `/help`, `/quit`
- Username validation
- Sliding-window per-client rate limiting
- Graceful process shutdown
- Concurrency-safe shared state
- Unit tests for room state and rate limiting

## Run

```bash
go test ./...
go run . -addr :8080
```

Connect from another terminal:

```bash
nc localhost 8080
```

## Architecture

`main.go` owns listener lifecycle and signal handling. `server.Server` accepts connections, `Client` owns each connection, and `Hub` synchronizes users and rooms. The design keeps network I/O at the client boundary and protects shared state with `sync.RWMutex`.

## Project provenance

The project direction was inspired by Uthman Oladele's open-source GO-CHAT project (`codetesla51/go-chat-server`), which is MIT licensed. The code in this repository is a new implementation rather than a verbatim copy of that source tree.

Upstream reference: https://github.com/codetesla51/go-chat-server

Upstream copyright and MIT license remain applicable to any upstream material separately reused from that project.

## Author

Adewale Babalola — Backend Engineer & Electrical/Electronics Engineer
