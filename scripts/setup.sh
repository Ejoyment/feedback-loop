#!/bin/bash

echo "Setting up AI Developer Assistant Platform..."

# Check dependencies
command -v cargo >/dev/null 2>&1 || { echo "Rust/Cargo required. Install from https://rustup.rs/"; exit 1; }
command -v go >/dev/null 2>&1 || { echo "Go required. Install from https://go.dev/"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "Docker required for Qdrant"; exit 1; }

# Install protobuf compiler
if ! command -v protoc >/dev/null 2>&1; then
    echo "Installing protobuf compiler..."
    # Platform-specific installation would go here
fi

# Start Qdrant
echo "Starting Qdrant vector database..."
docker run -d -p 6333:6333 -p 6334:6334 --name qdrant qdrant/qdrant:latest

# Build Rust engine
echo "Building Rust vector engine..."
cd rust-engine && cargo build --release
cd ..

# Build Go services
echo "Building Go services..."
cd services && go mod download && go build -o bin/api-service ./cmd/main.go
cd ..

echo "Setup complete! Copy .env.example to .env and add your API keys."
echo "Run 'make run-rust' in one terminal and 'make run-go' in another."
