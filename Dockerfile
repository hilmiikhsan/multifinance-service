# Stage 1: Build
FROM golang:1.22.8-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set build context
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Debug build platform
RUN echo "Building for TARGETPLATFORM=$TARGETPLATFORM, TARGETARCH=$TARGETARCH"

# Build the binary for the target platform
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o multifinance-service ./cmd/bin/main.go

# Stage 2: Final Image
FROM alpine:latest

# Install runtime dependencies and debugging tools
RUN apk add --no-cache ca-certificates file

# Set working directory
WORKDIR /app

# Copy compiled binary from the builder stage
COPY --from=builder /app/multifinance-service .

# Debug binary details
RUN file multifinance-service

# Expose required ports
EXPOSE 9090
EXPOSE 7000

# Run the application
CMD ["./multifinance-service"]
