# Deployment Guide

## Prerequisites

- Docker & Docker Compose
- Go 1.22+
- Rust 1.76+
- OpenAI and/or Anthropic API keys

## Quick Start (Docker)

1. Copy environment file:
```bash
cp .env.example .env
```

2. Add your API keys to `.env`

3. Start all services:
```bash
make docker-up
```

4. Verify health:
```bash
curl http://localhost:8080/health
```

## Local Development

### Terminal 1: Start Qdrant
```bash
docker run -p 6333:6333 -p 6334:6334 qdrant/qdrant:latest
```

### Terminal 2: Start Rust Vector Engine
```bash
make run-rust
```

### Terminal 3: Start Go API Service
```bash
make run-go
```

## Indexing Documentation

```bash
# Place your docs in ./docs directory
bash scripts/index_docs.sh
```

## Testing

```bash
make test
```

## Production Deployment

### Kubernetes
- Helm charts in `/k8s` directory
- Configure secrets for API keys
- Set resource limits based on load

### Scaling Considerations
- Vector engine: CPU-bound, scale horizontally
- API service: Stateless, auto-scale based on requests
- Qdrant: Persistent storage, consider managed service
