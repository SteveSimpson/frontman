FROM golang:1.24-bookworm AS gobase
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM gobase AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/frontman ./cmd/main.go

# Runtime Stage
FROM alpine:latest AS api
WORKDIR /app
COPY frontman_config.json .
COPY --from=builder /app/frontman .
ENTRYPOINT ["./frontman"]
