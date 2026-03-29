# API Documentation

## RAG Query Endpoint

### POST /api/v1/rag/query

Query internal documentation using natural language.

**Request:**
```json
{
  "query": "How do I implement authentication?",
  "collection": "internal_docs",
  "top_k": 5
}
```

**Response:**
```json
{
  "answer": "To implement authentication in our microservices...",
  "sources": [
    {
      "id": "auth-guide.md",
      "content": "Authentication Guide...",
      "type": "documentation",
      "metadata": {}
    }
  ],
  "metadata": {
    "tokens_used": 450
  }
}
```

## Test Generation Endpoint

### POST /api/v1/testgen/generate

Generate unit tests for TypeScript code.

**Request:**
```json
{
  "source_code": "export function add(a: number, b: number) { return a + b; }",
  "file_path": "src/utils/math.ts",
  "test_framework": "jest"
}
```

**Response:**
```json
{
  "test_code": "import { add } from '../utils/math';\n\ndescribe('add', () => {...});",
  "coverage": ["✓ add() covered"],
  "tokens_used": 320
}
```

## Index Document Endpoint

### POST /api/v1/rag/index

Index a document into the RAG system.

**Request:**
```json
{
  "id": "doc-123",
  "content": "This is documentation content...",
  "type": "documentation",
  "metadata": {
    "source": "internal-wiki",
    "author": "engineering-team"
  }
}
```

**Response:**
```json
{
  "status": "indexed",
  "id": "doc-123"
}
```
