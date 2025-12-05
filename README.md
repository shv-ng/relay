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
  Total:        0.4345 secs
  Slowest:      0.0557 secs
  Fastest:      0.0001 secs
  Average:      0.0039 secs
  Requests/sec: 23016.9177

  Total data:   120000 bytes
  Size/request: 12 bytes

Response time histogram:
  0.000 [1]     |
  0.006 [7695]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.011 [1790]  |■■■■■■■■■
  0.017 [335]   |■■
  0.022 [77]    |
  0.028 [6]     |
  0.033 [25]    |
  0.039 [52]    |
  0.045 [15]    |
  0.050 [3]     |
  0.056 [1]     |


Latency distribution:
  10% in 0.0004 secs
  25% in 0.0010 secs
  50% in 0.0024 secs
  75% in 0.0053 secs
  90% in 0.0091 secs
  95% in 0.0113 secs
  99% in 0.0239 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0001 secs, 0.0001 secs, 0.0557 secs
  DNS-lookup:   0.0002 secs, 0.0000 secs, 0.0387 secs
  req write:    0.0001 secs, 0.0000 secs, 0.0180 secs
  resp wait:    0.0028 secs, 0.0001 secs, 0.0179 secs
  resp read:    0.0007 secs, 0.0000 secs, 0.0135 secs

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
