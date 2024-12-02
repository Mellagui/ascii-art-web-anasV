# Stage 1: Build the Go backend
FROM debian:bullseye-slim AS builder

# Adding metadata
LABEL maintainer="Zone01 talents(myacoubi, anlazaar, amellagu)"
LABEL version="1.0.0"
LABEL description="An image to run the ASCII Art Web application."
LABEL website="http://localhost/8080"
LABEL build-date="Decembre 2024"

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

WORKDIR /app/backend

# Copy the built Go binary from the builder stage
COPY --from=builder /app/ascii-art-web ./ascii-art-web

# Copy the banners directory
COPY backend/banners ./banners
COPY backend/utils ./utils

# Copy other required files
COPY templates ../templates

# Copy the frontend files
COPY static ../static

# Expose the port for the Go server
EXPOSE 8080

# Command to run the Go server
CMD ["./ascii-art-web"]
