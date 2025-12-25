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

<details>
  <summary>Round Robin</summary>

```bash
> hey -n 10000 -c 100 http://localhost:8080/

Summary:
  Total:        2.0686 secs
  Slowest:      0.1000 secs
  Fastest:      0.0003 secs
  Average:      0.0198 secs
  Requests/sec: 4834.2872

  Total data:   290000 bytes
  Size/request: 29 bytes

Response time histogram:
  0.000 [1]     |
  0.010 [3140]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.020 [2668]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.030 [2021]  |■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.040 [1126]  |■■■■■■■■■■■■■■
  0.050 [646]   |■■■■■■■■
  0.060 [277]   |■■■■
  0.070 [91]    |■
  0.080 [20]    |
  0.090 [8]     |
  0.100 [2]     |


Latency distribution:
  10% in 0.0035 secs
  25% in 0.0082 secs
  50% in 0.0169 secs
  75% in 0.0283 secs
  90% in 0.0407 secs
  95% in 0.0477 secs
  99% in 0.0620 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0001 secs, 0.0003 secs, 0.1000 secs
  DNS-lookup:   0.0001 secs, 0.0000 secs, 0.0178 secs
  req write:    0.0001 secs, 0.0000 secs, 0.0125 secs
  resp wait:    0.0192 secs, 0.0002 secs, 0.0999 secs
  resp read:    0.0003 secs, 0.0000 secs, 0.0130 secs

Status code distribution:
  [200] 10000 responses
```  
</details>

<details>
  <summary>Random</summary>

```bash
> hey -n 10000 -c 100 http://localhost:8080/

Summary:
  Total:        2.1078 secs
  Slowest:      0.1027 secs
  Fastest:      0.0003 secs
  Average:      0.0201 secs
  Requests/sec: 4744.2859

  Total data:   290000 bytes
  Size/request: 29 bytes

Response time histogram:
  0.000 [1]     |
  0.011 [3088]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.021 [2871]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.031 [1943]  |■■■■■■■■■■■■■■■■■■■■■■■■■
  0.041 [1144]  |■■■■■■■■■■■■■■■
  0.051 [583]   |■■■■■■■■
  0.062 [238]   |■■■
  0.072 [84]    |■
  0.082 [35]    |
  0.092 [10]    |
  0.103 [3]     |


Latency distribution:
  10% in 0.0037 secs
  25% in 0.0085 secs
  50% in 0.0169 secs
  75% in 0.0283 secs
  90% in 0.0407 secs
  95% in 0.0481 secs
  99% in 0.0639 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0001 secs, 0.0003 secs, 0.1027 secs
  DNS-lookup:   0.0001 secs, 0.0000 secs, 0.0213 secs
  req write:    0.0001 secs, 0.0000 secs, 0.0175 secs
  resp wait:    0.0196 secs, 0.0002 secs, 0.1016 secs
  resp read:    0.0002 secs, 0.0000 secs, 0.0135 secs

Status code distribution:
  [200] 10000 responses
```  
</details>

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
