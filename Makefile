.PHONY: build run test clean docker-up docker-down

build:
	cd rust-engine && cargo build --release
	cd services && go build -o bin/api-service ./cmd/main.go

run-rust:
	cd rust-engine && cargo run --release

run-go:
	cd services && go run cmd/main.go

test:
	cd rust-engine && cargo test
	cd services && go test ./...

clean:
	cd rust-engine && cargo clean
	cd services && rm -rf bin/

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

proto:
	cd rust-engine && cargo build
