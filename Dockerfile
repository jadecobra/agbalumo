
# Build Stage for App
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy only necessary Go source code
COPY cmd cmd
COPY internal internal
COPY main.go .

# Build the application
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server main.go

# Build Stage for Litestream (to patch CVE-2026-33186 in gRPC dependency)
FROM golang:1.26-alpine AS litestream-builder

WORKDIR /src
RUN apk add --no-cache git && \
    git clone --depth 1 --branch v0.5.10 https://github.com/benbjohnson/litestream.git . && \
    go mod edit -replace google.golang.org/grpc=google.golang.org/grpc@v1.79.3 && \
    go mod tidy && \
    CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/litestream ./cmd/litestream

# Runtime Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates and runtime dependencies
# su-exec is for dropping root privileges in entrypoint
# libc6-compat is for glibc compatibility (litestream)
# We explicitly upgrade libcrypto3 and libssl3 to resolve CVE-2024-13176
RUN apk --no-cache upgrade && \
    apk --no-cache add ca-certificates tzdata wget bash su-exec libc6-compat && \
    apk add --no-cache --upgrade libcrypto3 libssl3

# Set time zone
ENV TZ=UTC

# Create appuser and setup directories
RUN adduser -D -u 1000 appuser && \
    mkdir -p /data && \
    chown appuser:appuser /data /app

# Copy binaries (Litestream from source-build with gRPC patch, app from local builder)
COPY --from=litestream-builder --chown=appuser:appuser /app/litestream /usr/local/bin/litestream
COPY --from=builder --chown=appuser:appuser /app/server .

# Copy UI assets (Templates & Static)
COPY --chown=appuser:appuser ui ui
COPY --chown=appuser:appuser config config

# Expose port
EXPOSE 8080

ENV AGBALUMO_ENV=production
ENV PORT=8080
ENV DATABASE_URL=/data/agbalumo.db

# Copy litestream config and entrypoint
COPY --chown=appuser:appuser etc/litestream.yml /etc/litestream.yml
COPY --chown=appuser:appuser scripts/entrypoint.sh /app/entrypoint.sh

# Ensure entrypoint is executable
RUN chmod +x /app/entrypoint.sh

# Explicitly set PATH to include common binary locations
ENV PATH="/usr/local/bin:/usr/bin:/bin:/app:${PATH}"

# Start as root to allow entrypoint to chown volumes, then drop to appuser
USER root

# Run the entrypoint script
CMD ["/app/entrypoint.sh"]
