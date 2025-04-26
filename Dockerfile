# Start from a minimal Go image
FROM golang:1.24.2-alpine

# Set working directory inside container
WORKDIR /app

# Copy go mod files first and download deps (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the Go binary
RUN go build -o main .

# Expose internal port
EXPOSE 4444

# Run the binary
CMD ["./main"]
