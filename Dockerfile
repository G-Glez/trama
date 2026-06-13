FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /trama ./cmd/api/
RUN CGO_ENABLED=0 GOOS=linux go build -o /trama-cli ./cmd/cli/

FROM alpine:3.21
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /app
COPY --from=builder /trama .
COPY --from=builder /trama-cli .
COPY entrypoint.sh /entrypoint.sh
COPY localdb /app/localdb
COPY migrations /app/migrations
RUN chmod +x /entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
