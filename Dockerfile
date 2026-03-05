# Build Stage
FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies (if needed, but modernc sqlite is pure go mostly)
# apk add --no-cache git 

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy only necessary Go source code
# This prevents cache invalidation when non-Go files (like UI, docs) change
COPY cmd cmd
COPY internal internal
COPY main.go .

# Build the application
# CGO_ENABLED=0 since modernc.org/sqlite is pure Go (mostly) and for static binary
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server main.go

# Runtime Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for external API calls
RUN apk --no-cache add ca-certificates tzdata wget bash

# Download Litestream
ADD https://github.com/benbjohnson/litestream/releases/download/v0.3.13/litestream-v0.3.13-linux-amd64.tar.gz /tmp/litestream.tar.gz
RUN tar -C /usr/local/bin -xzf /tmp/litestream.tar.gz && rm /tmp/litestream.tar.gz

# Create a non-root user
RUN adduser -D -g '' appuser

# Copy binary from builder
COPY --from=builder /app/server .

# Copy UI assets (Templates & Static)
# Copy from local context, NOT builder, to allow fast UI updates
COPY ui ui

# Expose port
EXPOSE 8080

ENV AGBALUMO_ENV=production
ENV PORT=8080
ENV DATABASE_URL=/data/agbalumo.db

# Copy litestream config and entrypoint
COPY etc/litestream.yml /etc/litestream.yml
COPY scripts/entrypoint.sh /app/entrypoint.sh

# Use non-root user
USER appuser

# Run the entrypoint script
CMD ["/app/entrypoint.sh"]
