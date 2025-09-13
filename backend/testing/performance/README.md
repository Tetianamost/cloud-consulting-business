# Performance Tests

This directory contains performance and load testing utilities for benchmarking system performance under various conditions.

## Test Files

Performance tests will be moved here from the backend root directory in subsequent tasks:

- Load testing utilities (`test_performance_*.go`, `test_load_*.go`)
- Concurrency tests
- Memory usage profiling
- CPU usage benchmarks
- Database performance tests

## Running Performance Tests

```bash
# Run specific performance test
cd backend && go run testing/performance/test_name.go

# Run performance benchmarks
cd backend && go test -bench=. ./testing/performance/

# Run with memory profiling
cd backend && go test -memprofile=mem.prof ./testing/performance/

# Run with CPU profiling  
cd backend && go test -cpuprofile=cpu.prof ./testing/performance/
```

## Requirements

- Sufficient system resources for load testing
- Test database with performance data
- Monitoring tools for metrics collection