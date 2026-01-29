# Build Stage
FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies (if needed, but modernc sqlite is pure go mostly)
# apk add --no-cache git 

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
# CGO_ENABLED=0 since modernc.org/sqlite is pure Go (mostly) and for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server/main.go

# Runtime Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for external API calls (Gemini)
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/server .

# Copy UI assets (Templates & Static)
COPY --from=builder /app/ui ./ui

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]
