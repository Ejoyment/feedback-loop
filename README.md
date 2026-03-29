# AI Developer Assistant Platform

A production-ready AI-powered developer tooling platform that reduces context switching and automates engineering workflows.

## Architecture

- **Go Backend**: REST API, RAG orchestration, LLM integration
- **Rust Engine**: High-performance vector search and embeddings
- **LangGraph**: Multi-agent workflow orchestration
- **Vector DB**: Semantic search over documentation and code

## Components

1. **RAG Query Service** (`/services/rag-service/`) - Natural language queries over internal docs
2. **Test Generator** (`/services/test-generator/`) - Automated unit test generation for TypeScript
3. **LLM Gateway** (`/services/llm-gateway/`) - Unified interface for OpenAI/Claude APIs
4. **Vector Engine** (`/rust-engine/`) - Fast embedding and similarity search

## Quick Start

```bash
# Build Rust engine
cd rust-engine && cargo build --release

# Start Go services
cd services && go run cmd/main.go

# Run test generator
cd services/test-generator && go run main.go
```

## Features

- Natural language documentation search
- Codebase schema querying
- Automated test generation for legacy code
- Multi-LLM support (OpenAI, Claude)
- LangGraph-based agent workflows
