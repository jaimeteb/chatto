FROM golang:1.13-alpine AS build
WORKDIR /chatto

COPY . .
RUN apk add --update --no-cache ca-certificates git
RUN CGO_ENABLED=0 go build -o /go/bin/chatto cmd/chatto/main.go
RUN mkdir /chatto/data

FROM scratch

COPY --from=build /go/bin/chatto /go/bin/chatto
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /chatto/data /chatto/data

ENTRYPOINT ["/go/bin/chatto"]
CMD ["-path", "/chatto/data"]
