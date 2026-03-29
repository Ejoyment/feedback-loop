FROM rust:1.76 as builder

WORKDIR /app
COPY Cargo.toml Cargo.lock ./
COPY proto ./proto
COPY build.rs ./
COPY src ./src

RUN apt-get update && apt-get install -y protobuf-compiler
RUN cargo build --release

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/target/release/vector-engine /usr/local/bin/

EXPOSE 50051
CMD ["vector-engine"]
