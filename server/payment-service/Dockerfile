# -------- Build Stage --------
FROM golang:1.24.3-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service .

# -------- Runtime Stage --------
FROM alpine:latest

# Install certificates for HTTPS communication (Razorpay needs this)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the compiled Go binary from builder stage
COPY --from=builder /app/payment-service .

# Expose the service port
EXPOSE 8084

# Optional health check to monitor container status
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --quiet --tries=1 --spider http://localhost:8084/health || exit 1

# Run the Go binary
CMD ["./payment-service"]
