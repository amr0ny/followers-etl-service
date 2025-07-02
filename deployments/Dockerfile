FROM golang:1.23.10-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY migrations ./
RUN go mod download

COPY . .

RUN go build -o main ./


FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
CMD ["./main"]