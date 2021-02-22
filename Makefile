.PHONY: test generate lint snapshot release godoc

godoc:
	open "http://localhost:6060/pkg/github.com/jaimeteb/chatto/"
	godoc -http=:6060

test:
	go test -race ./... -cover -coverprofile=coverage.txt

generate:
	go generate ./...

lint:
	golangci-lint run

snapshot:
	goreleaser release --snapshot --rm-dist

release:
	goreleaser release --rm-dist
