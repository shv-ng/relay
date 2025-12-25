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

```bash
make run
``` 

### Load Balancing Algorithms

- [x] **Round Robin** 
- [ ] **Least Connections** 
- [ ] **Weighted Round Robin** 
- [ ] **Random Selection**
- [ ] **IP Hash / Sticky Sessions** 

## Benchmarks

Benchmarks performed using [hey](https://github.com/rakyll/hey): `hey -n 10000 -c 100 http://localhost:8000/`

### Round Robin
```bash
> hey -n 10000 -c 100 http://localhost:8000/

Summary:
  Total:        0.5318 secs
  Slowest:      0.0671 secs
  Fastest:      0.0001 secs
  Average:      0.0049 secs
  Requests/sec: 18803.2185

  Total data:   290000 bytes
  Size/request: 29 bytes

Response time histogram:
  0.000 [1]     |
  0.007 [7557]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.013 [1753]  |■■■■■■■■■
  0.020 [441]   |■■
  0.027 [178]   |■
  0.034 [31]    |
  0.040 [6]     |
  0.047 [14]    |
  0.054 [3]     |
  0.060 [9]     |
  0.067 [7]     |


Latency distribution:
  10% in 0.0005 secs
  25% in 0.0013 secs
  50% in 0.0029 secs
  75% in 0.0066 secs
  90% in 0.0111 secs
  95% in 0.0157 secs
  99% in 0.0240 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0001 secs, 0.0001 secs, 0.0671 secs
  DNS-lookup:   0.0001 secs, 0.0000 secs, 0.0194 secs
  req write:    0.0002 secs, 0.0000 secs, 0.0485 secs
  resp wait:    0.0034 secs, 0.0001 secs, 0.0222 secs
  resp read:    0.0010 secs, 0.0000 secs, 0.0308 secs

Status code distribution:
  [200] 10000 responses
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


## Contributing

This is a learning project built in agile sprints. Feel free to:

- Fork and experiment with your own algorithms
- Optimize existing implementations
- Add metrics and monitoring
- Suggest improvements via issues

Built with ❤️ to learn Go and distributed systems.
