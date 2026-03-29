#!/bin/bash

# Script to index internal documentation into the RAG system

API_URL="http://localhost:8080/api/v1/rag/index"
DOCS_DIR="./docs"

echo "Indexing documentation from $DOCS_DIR..."

find "$DOCS_DIR" -type f \( -name "*.md" -o -name "*.txt" \) | while read -r file; do
    echo "Indexing: $file"
    
    content=$(cat "$file")
    filename=$(basename "$file")
    
    json_payload=$(jq -n \
        --arg id "$filename" \
        --arg content "$content" \
        --arg type "documentation" \
        '{id: $id, content: $content, type: $type, metadata: {source: $id}}')
    
    curl -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -d "$json_payload"
    
    echo ""
done

echo "Documentation indexing complete!"
