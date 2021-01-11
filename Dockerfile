FROM golang:1.13-alpine AS builder
WORKDIR /chatto

COPY . .
RUN apk add --update --no-cache ca-certificates git
RUN go install chatto.go
RUN mkdir /data

CMD ["chatto", "-path", "data"]
