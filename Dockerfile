# build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/url-shortener ./cmd/url-shortener

# run stage
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/url-shortener /app/url-shortener
COPY config /app/config

EXPOSE 8080

ENV CONFIG_PATH=/app/config/local.yaml
CMD ["/app/url-shortener"]