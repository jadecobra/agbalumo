
# Build Stage for App
FROM golang:1.25-alpine AS builder

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

# Runtime Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates and upgrade all packages to get latest security fixes
RUN apk --no-cache upgrade && \
    apk --no-cache add ca-certificates tzdata wget bash

# Create appuser and setup directories
RUN adduser -D -u 1000 appuser && \
    mkdir -p /data && \
    chown appuser:appuser /data /app

# Copy binaries (Litestream from official image, app from local builder)
COPY --from=litestream/litestream:0.5.10 --chown=appuser:appuser /usr/local/bin/litestream /usr/local/bin/litestream
COPY --from=builder --chown=appuser:appuser /app/server .

# Copy UI assets (Templates & Static)
# Copy from local context, NOT builder, to allow fast UI updates
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

# Use non-root user
USER appuser

# Run the entrypoint script
CMD ["/app/entrypoint.sh"]
