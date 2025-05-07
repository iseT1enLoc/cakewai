# Build Stage
FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Ensure a statically compiled binary
ENV CGO_ENABLED=0

# Debug: Print files to verify main.go exists
RUN ls -lah && go build -o main .

# Minimal Production Image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080

CMD ["./main"]
