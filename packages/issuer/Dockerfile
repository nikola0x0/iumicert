# Multi-stage build for IU-MiCert Issuer Backend
FROM golang:1.24.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /build

# Set target OS (but let GOARCH be determined by build platform or --platform flag)
ENV GOOS=linux

# Copy go mod files first for better caching
COPY packages/issuer/go.mod packages/issuer/go.sum ./

# Copy crypto dependency
COPY packages/crypto /build/../crypto

# Download dependencies
RUN go mod download

# Copy the issuer application
COPY packages/issuer/ .

# Build the binary from cmd/main.go
# Note: GOARCH is not set here, so it will match the build platform
# For CI/CD, Docker Buildx will set it via --platform flag
RUN CGO_ENABLED=1 go build -a -installsuffix cgo -ldflags="-w -s" -o micert ./cmd

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/micert .

# Copy configuration files
COPY --from=builder /build/config ./config
COPY --from=builder /build/.env.example ./.env.example

# Create necessary directories
RUN mkdir -p data/student_journeys/students \
    data/student_journeys/terms \
    data/verkle_terms \
    data/verkle_trees \
    publish_ready/receipts \
    publish_ready/roots \
    publish_ready/transactions \
    logs

# Set environment variables
ENV TZ=UTC
ENV PORT=8080

# Expose the API port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the binary
CMD ["./micert", "serve"]
