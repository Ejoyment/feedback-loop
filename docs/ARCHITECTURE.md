# System Architecture

## Overview

This platform implements a production-grade AI developer assistant that reduces context switching through RAG-based documentation queries and automated test generation.

## Core Components

### 1. Rust Vector Engine (Port 50051)
- High-performance embedding generation using fastembed
- gRPC service for low-latency vector operations
- Qdrant integration for semantic search
- Handles 1000+ embeddings/sec

### 2. Go API Service (Port 8080)
- REST API for RAG queries and test generation
- LLM gateway supporting OpenAI and Claude
- Request orchestration and caching
- Middleware for auth, logging, rate limiting

### 3. LangGraph Workflow Engine
- Multi-agent orchestration
- State management across agent steps
- Conditional routing based on query intent
- Error recovery and retry logic

### 4. RAG Service
- Semantic search over internal docs and schemas
- Context-aware answer generation
- Source attribution and relevance scoring
- Supports multiple collections (docs, code, APIs)

### 5. Test Generator
- Parses TypeScript source code
- Queries RAG for testing patterns
- Generates Jest/Vitest/Mocha tests
- Coverage analysis and validation

## Data Flow

```
User Query → API Gateway → LangGraph Workflow
                              ↓
                         Analyze Intent
                              ↓
                    RAG Service (Query Docs)
                              ↓
                    Vector Engine (Embed + Search)
                              ↓
                         Qdrant (Retrieve)
                              ↓
                    LLM Gateway (Generate Answer)
                              ↓
                         Response to User
```

## Technology Stack

- Go 1.22: API services, business logic
- Rust 1.76: Performance-critical vector operations
- Qdrant: Vector database
- OpenAI GPT-4 / Claude 3: LLM providers
- gRPC: Inter-service communication
- Docker: Containerization
