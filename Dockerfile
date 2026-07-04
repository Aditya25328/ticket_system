# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency manifests and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ticket-system ./cmd/main.go

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/ticket-system .

# Create database directory so SQLite has a folder to write into
RUN mkdir -p database

EXPOSE 8080

CMD ["./ticket-system"]
