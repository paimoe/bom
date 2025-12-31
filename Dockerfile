# Start from the official Go image for building your binary
FROM golang:1.25.5 AS builder

LABEL org.opencontainers.image.source=https://github.com/paimoe/bom

WORKDIR /app

# Copy go mod and sum files, then download dependencies (better cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app

# Use a minimal image for running
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
COPY static static

ENV PORT=80
EXPOSE 80

# Run the binary
CMD ["./app"]
