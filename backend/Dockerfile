# Use a more basic approach - build from source
FROM alpine:latest

# Install Go and other dependencies
RUN apk add --no-cache go git ca-certificates tzdata

# Set Go environment
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
ENV GO111MODULE=on

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/server/main.go

EXPOSE 8061

# Run the built binary
CMD ["./main"]