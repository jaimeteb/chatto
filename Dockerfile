FROM golang:1.13-alpine AS builder
RUN apk add --update --no-cache ca-certificates git
RUN go get -u github.com/jaimeteb/chatto
ENTRYPOINT ["chatto"]
