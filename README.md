# Relay - Simple Load Balancer in Go

A lightweight HTTP load balancer built from scratch in Go to understand load balancing algorithms and their performance characteristics.

## Description

This project implements a simple load balancer that distributes incoming HTTP requests across multiple backend servers using various algorithms. Each algorithm is benchmarked to compare performance and behavior under load.

## Motivation

Load balancers are critical infrastructure components, but their inner workings can seem like magic. This project demystifies load balancing by:

- Building a functional load balancer from first principles
- Implementing and comparing different distribution algorithms
- Measuring real performance differences with benchmarks
- Learning Go's HTTP server and client capabilities

Perfect for understanding how production load balancers (nginx, HAProxy, cloud load balancers) work under the hood.

## Quick Start

**Prerequisites**: Go 1.23+ installed

**Run the backend server**:
```bash
cd server
PORT=8000 go run .
```

Test the endpoints:
```bash
curl http://localhost:8000/health  # Returns 200 OK
curl http://localhost:8000/        # Returns "Hello World"
```

## Usage

### Running Backend Servers

Start multiple backend instances on different ports:
```bash
# Terminal 1
cd server && PORT=8001 go run .

# Terminal 2
cd server && PORT=8002 go run .

# Terminal 3
cd server && PORT=8003 go run .
```

Each server responds with "Hello World" on `/` and returns 200 OK on `/health`.

### Running the Load Balancer

*(Coming in Sprint 2)*

### Load Balancing Algorithms

*(Will be added as implemented)*

- **Round Robin** - Sprint 2
- **Least Connections** - Sprint 5
- **Weighted Round Robin** - Sprint 6
- **Random Selection** - Sprint 7
- **IP Hash / Sticky Sessions** - Sprint 8

## Benchmarks

Benchmarks performed using [hey](https://github.com/rakyll/hey): `hey -n 10000 -c 100 http://localhost:8000/`

### Round Robin
*(Coming in Sprint 4)*

### Least Connections
*(Coming in Sprint 5)*

### Weighted Round Robin
*(Coming in Sprint 6)*

### Random Selection
*(Coming in Sprint 7)*

### IP Hash
*(Coming in Sprint 8)*

## Project Structure
```
.
├── main.go           # Entry point
├── server/
│   └── main.go       # Backend HTTP server (Sprint 0 ✓)
├── lb/               # Load balancer package (Sprint 2+)
├── Makefile          # Build commands
├── go.mod
└── README.md
```

## Development

**Build**:
```bash
make build
```

**Run tests**:
```bash
make test
```

**Development mode** (with hot reload):
```bash
make dev
```

## Sprint Progress

- [x] **Sprint 0**: Basic backend server with `/health` and `/` endpoints
- [ ] **Sprint 1**: Multiple backend instances
- [ ] **Sprint 2**: Round robin load balancer
- [ ] **Sprint 3**: Health checking
- [ ] **Sprint 4**: Benchmarking setup
- [ ] **Sprint 5+**: Additional algorithms

## Contributing

This is a learning project built in agile sprints. Feel free to:

- Fork and experiment with your own algorithms
- Optimize existing implementations
- Add metrics and monitoring
- Suggest improvements via issues

Built with ❤️ to learn Go and distributed systems.
