FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8061

# Run the Go application directly instead of building a binary
CMD ["go", "run", "./cmd/server/main.go"]