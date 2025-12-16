# Load Testing Guide

This directory contains load testing configurations using Artillery for performance and stress testing of the application.

## Overview

Load tests validate the system's ability to handle expected and peak traffic loads. We use Artillery, a modern, powerful, and easy-to-use load testing toolkit.

## Prerequisites

Install Artillery globally:

```bash
npm install -g artillery@latest
```

Or use it via npx (no installation needed):

```bash
npx artillery@latest --version
```

## Test Configurations

### 1. Backend API Load Tests (`artillery-backend.yml`)

Tests the Go backend API service under various load conditions.

**Scenarios:**
- Health check monitoring
- User authentication flow
- Project CRUD operations
- Blueprint workflow
- Bid operations

**Load Phases:**
1. **Warm up** (30s): 5 requests/sec
2. **Ramp up** (60s): 10 → 50 requests/sec
3. **Sustained load** (120s): 50 requests/sec
4. **Spike test** (30s): 100 requests/sec
5. **Cool down** (30s): 10 requests/sec

**Performance Targets:**
- Max error rate: 1%
- P95 response time: < 500ms
- P99 response time: < 1000ms

### 2. AI Service Load Tests (`artillery-ai-service.yml`)

Tests the Python FastAPI AI service, including inference endpoints.

**Scenarios:**
- Health check
- Blueprint analysis
- Bid generation
- Mixed operations

**Load Phases:**
1. **Warm up** (20s): 2 requests/sec
2. **Ramp up** (60s): 5 → 15 requests/sec
3. **Sustained load** (90s): 15 requests/sec
4. **Spike test** (20s): 30 requests/sec

**Performance Targets:**
- Max error rate: 2% (AI operations can have transient issues)
- P95 response time: < 2s
- P99 response time: < 5s

## Running Load Tests

### Run Tests Locally

Make sure all services are running first:

```bash
# Start all services
cd /home/runner/work/fantastic-octo-memory/fantastic-octo-memory
make dev
```

Then run the load tests:

```bash
# Test backend API
artillery run load-tests/artillery-backend.yml

# Test AI service
artillery run load-tests/artillery-ai-service.yml

# Run with custom target URL
artillery run --target http://your-api.com load-tests/artillery-backend.yml
```

### Run Tests with Detailed Output

```bash
# Run with JSON output
artillery run --output report.json load-tests/artillery-backend.yml

# Generate HTML report
artillery report report.json

# Run in debug mode
DEBUG=http artillery run load-tests/artillery-backend.yml
```

### Run Tests Against Production

```bash
# Backend
artillery run --target https://api.yourdomain.com load-tests/artillery-backend.yml

# AI Service
artillery run --target https://ai.yourdomain.com load-tests/artillery-ai-service.yml
```

## Understanding Results

Artillery provides detailed metrics:

### Key Metrics

- **HTTP codes**: Distribution of response status codes
- **Response times**: 
  - `min`: Fastest response
  - `max`: Slowest response
  - `median`: 50th percentile
  - `p95`: 95th percentile (95% of requests faster than this)
  - `p99`: 99th percentile
- **RPS**: Requests per second
- **Errors**: Count of failed requests
- **Scenarios**: Number of virtual users completed

### Sample Output

```
Summary report @ 14:30:15
  Scenarios launched:  1000
  Scenarios completed: 995
  Requests completed:  4980
  Mean response/sec:   16.6
  Response time (msec):
    min: 12
    max: 1234
    median: 98
    p95: 245
    p99: 567
  Scenario counts:
    Health Check: 100 (100%)
  Codes:
    200: 4980
  Errors:
    ETIMEDOUT: 5
```

## CI/CD Integration

Load tests can be integrated into CI/CD pipelines for continuous performance monitoring.

### GitHub Actions Example

```yaml
- name: Run Load Tests
  run: |
    npm install -g artillery
    artillery run --output report.json load-tests/artillery-backend.yml
    artillery report report.json
  
- name: Upload Load Test Report
  uses: actions/upload-artifact@v3
  with:
    name: load-test-report
    path: report.json
```

## Advanced Usage

### Custom Environment Variables

```bash
# Set custom variables
artillery run \
  -e production \
  --target https://api.yourdomain.com \
  load-tests/artillery-backend.yml
```

### Running Specific Scenarios

Edit the YAML file to comment out scenarios you don't want to run, or adjust weights.

### Performance Benchmarking

Run tests multiple times and compare results:

```bash
# Baseline
artillery run --output baseline.json load-tests/artillery-backend.yml

# After optimization
artillery run --output optimized.json load-tests/artillery-backend.yml

# Compare (manual analysis of JSON)
```

## Monitoring During Load Tests

While running load tests, monitor your services:

### Docker Stats

```bash
docker stats
```

### Application Logs

```bash
# Backend logs
docker logs -f construction-backend

# AI service logs
docker logs -f construction-ai-service
```

### System Resources

```bash
# CPU and memory
htop

# Network
iftop
```

## Troubleshooting

### Test Failures

**Error: ECONNREFUSED**
- Services are not running
- Wrong target URL
- Firewall blocking requests

**Error: ETIMEDOUT**
- Services are overloaded
- Network issues
- Timeout too short for operation

**High Error Rate**
- Check service logs
- Verify database connections
- Check resource limits (CPU, memory, connections)

### Performance Issues

If tests show poor performance:

1. **Check service health**: Ensure all services are running properly
2. **Review logs**: Look for errors or warnings
3. **Check resources**: CPU, memory, database connections
4. **Database**: Check query performance, indexes
5. **Cache**: Verify Redis is working
6. **Network**: Check for network bottlenecks

## Best Practices

1. **Start small**: Begin with low load and gradually increase
2. **Warm up**: Always include a warm-up phase
3. **Realistic scenarios**: Test actual user workflows
4. **Monitor**: Watch system metrics during tests
5. **Baseline**: Establish baseline metrics before optimization
6. **Regular testing**: Run load tests regularly, not just before release
7. **Production-like environment**: Test on environments similar to production

## Performance Targets

### Backend API

| Metric | Target | Critical |
|--------|--------|----------|
| Availability | > 99.9% | > 99% |
| Response Time (P95) | < 500ms | < 1000ms |
| Response Time (P99) | < 1000ms | < 2000ms |
| Error Rate | < 0.1% | < 1% |
| Throughput | > 100 RPS | > 50 RPS |

### AI Service

| Metric | Target | Critical |
|--------|--------|----------|
| Availability | > 99.5% | > 99% |
| Response Time (P95) | < 2s | < 5s |
| Response Time (P99) | < 5s | < 10s |
| Error Rate | < 0.5% | < 2% |
| Throughput | > 20 RPS | > 10 RPS |

## Additional Resources

- [Artillery Documentation](https://www.artillery.io/docs)
- [Load Testing Best Practices](https://www.artillery.io/docs/guides/guides/best-practices)
- [Performance Testing Guide](https://martinfowler.com/articles/practical-test-pyramid.html)

## Next Steps

1. Run baseline tests on current system
2. Identify performance bottlenecks
3. Optimize based on findings
4. Re-run tests to verify improvements
5. Set up continuous load testing in CI/CD
6. Monitor production metrics

## Support

For issues or questions about load testing:
- Check Artillery docs: https://www.artillery.io/docs
- Review application logs
- Open an issue in the repository
