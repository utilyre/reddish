FROM golang:1.22.4-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go generate -v ./...
RUN go build -v ./cmd/reddish

FROM debian:12.5-slim

COPY --from=builder /app/reddish /usr/local/bin

ENV STORAGE_SERVER_ADDR=0.0.0.0:6979
CMD reddish
