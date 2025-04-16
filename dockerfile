# Stage 1: Build environment
FROM rust:1.86-slim AS builder

# Install required system dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy project files
COPY Cargo.toml Cargo.lock ./

# Copy source code and build
COPY src ./src
RUN cargo build --release

# Stage 2: Runtime environment (updated for OpenSSL 3.x)
FROM debian:bookworm-slim

# Install OpenSSL 3.x
RUN apt-get update && apt-get install -y \
    libssl3 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/target/release/bot /app/bot
ENTRYPOINT ["/app/bot"]
