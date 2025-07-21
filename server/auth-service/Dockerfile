# ----------- STAGE 1: Build the Go binary -----------
FROM golang:1.24.3-alpine AS builder

# Install git and other needed tools
RUN apk --no-cache add git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first, for caching deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service .

# ----------- STAGE 2: Minimal runtime image -----------
FROM alpine:3.18

# Install CA certificates (needed for HTTPS requests)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the compiled Go binary
COPY --from=builder /app/auth-service .

# Expose the app port
EXPOSE 8083

# Declare environment variable (PORT is optional here)
ENV PORT=8083

# Run the Go service
CMD ["./auth-service"]
