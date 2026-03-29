# Multi-stage build for Go backend
FROM golang:1.22-alpine AS build

WORKDIR /app

# Cache modules (backend lives in ./backend in the repo)
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/. .

# Build the binary
RUN go build -o app ./internal/cmd

# Minimal runtime image
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=build /app/app ./app

# Expose default port (Railway will set PORT env)
EXPOSE 8080

# Default address is controlled by HTTP_ADDR or PORT in code
ENV HTTP_ADDR=""

CMD ["./app"]
