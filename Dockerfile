# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN go build -o app main.go

# Run stage
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/app .
ENV PORT=8080
EXPOSE 8080
CMD ["./app"]
