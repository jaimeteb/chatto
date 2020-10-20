FROM golang:1.13-alpine AS builder
RUN apk add --update --no-cache ca-certificates git
RUN GO111MODULE=on go get -u github.com/jaimeteb/chatto
ENTRYPOINT ["chatto"]
