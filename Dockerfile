FROM golang:1.22.4 AS builder

RUN apt-get update
RUN apt-get install -y protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
RUN go install github.com/twitchtv/twirp/protoc-gen-twirp@v8.1.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go generate -v ./...
RUN go build -v ./cmd/reddish

FROM debian:12.5-slim

COPY --from=builder /app/reddish /usr/local/bin

ENV MODE=PROD
ENV LOG_LEVEL=INFO
ENV STORAGE_SERVER_ADDR=0.0.0.0:6979
CMD reddish
