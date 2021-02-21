.PHONY: test generate lint snapshot release

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
