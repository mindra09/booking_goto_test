

FROM golang:1.23-alpine AS development  

WORKDIR /app

# Install air for hot reload (compatible with Go 1.23)
RUN go install github.com/cosmtrek/air@v1.49.0

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Command to run air for hot reload
CMD ["air", "-c", ".air.toml"]