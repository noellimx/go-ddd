# Development Dockerfile for billing-service
FROM golang:1.26.1-alpine AS builder

ARG MAIN_PATH

WORKDIR /
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o service ${MAIN_PATH}

FROM alpine:latest
WORKDIR /app
COPY --from=builder /service ./service
EXPOSE 8080
CMD ["./service"]
