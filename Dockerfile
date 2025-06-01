# NameTidy Dockerfile (Production)
# Multi-stage build for minimal final image

# Build stage
FROM golang:1.23-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags for smaller binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o nametidy .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests (if needed)
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 nametidy && \
    adduser -u 1000 -G nametidy -s /bin/sh -D nametidy

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/nametidy /usr/local/bin/nametidy

# Create workspace directory for file processing
RUN mkdir -p /workspace && \
    chown -R nametidy:nametidy /workspace

# Switch to non-root user
USER nametidy

# Set default workspace
WORKDIR /workspace

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD nametidy --help > /dev/null || exit 1

# Default command
ENTRYPOINT ["nametidy"]
CMD ["--help"]