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

# Build the application
# CGO_ENABLED=0 since modernc.org/sqlite is pure Go (mostly) and for static binary
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server/main.go

# Runtime Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for external API calls
RUN apk --no-cache add ca-certificates tzdata

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

# Use non-root user
USER appuser

# Run the server
CMD ["./server"]
