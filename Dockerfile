# Stage 1: Build the Go backend
FROM debian:bullseye-slim AS builder

# Install Go
RUN apt-get update && apt-get install -y golang-go

WORKDIR /app

# Copy go.mod and go.sum files from the backend directory
COPY backend/go.mod ./backend/

# Download dependencies
RUN cd backend && go mod download

# Copy the rest of the backend source code
COPY backend ./backend

# Build the Go binary
RUN cd backend && go build -o /app/ascii-art-web main.go

# Stage 2: Use Debian Bullseye (to match GLIBC versions)
FROM debian:bullseye-slim

WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/ascii-art-web /app/ascii-art-web

# Copy the banners directory
COPY ascii-art/backend/banners /backend/banners

# Copy the frontend files
COPY frontend /app/frontend

# Expose the port for the Go server
EXPOSE 8080

# Command to run the Go server
CMD ["./ascii-art-web"]