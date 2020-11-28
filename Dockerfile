FROM golang:1.13-alpine AS builder
WORKDIR /chatto

COPY . .
RUN apk add --update --no-cache ca-certificates git && \
    go install chatto.go && \
    mkdir /data

CMD ["chatto", "--path", "data"]
