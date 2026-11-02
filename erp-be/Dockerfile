# Stage 1: Build binary Go
FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Stage 2: Runtime image minimal
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Buat direktori untuk file SQLite database
RUN mkdir -p /app/storage

COPY --from=builder /app/main /app/main

EXPOSE 8000

CMD ["./main"]
