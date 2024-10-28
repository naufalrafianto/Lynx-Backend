# Dockerfile
# Development stage
FROM golang:1.23-alpine AS development

# Install build essentials first
RUN apk update && \
    apk add --no-cache \
    gcc \
    musl-dev \
    git \
    make

# Set GOPATH and add it to PATH
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest


WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Expose port
EXPOSE 8080

# Command to run Air
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final production image
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]