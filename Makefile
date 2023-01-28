.PHONY: test generate lint snapshot release godoc

godoc:
	open "http://localhost:6060/pkg/github.com/jaimeteb/chatto/"
	godoc -http=:6060

test:
	go test -race ./... -cover -coverprofile=coverage.txt

test-docker:
	docker build -t jaimeteb/chatto-test -f utils/docker-test/Dockerfile .
	docker run --rm jaimeteb/chatto-test

generate:
	go generate ./...

lint:
	golangci-lint run

snapshot:
	goreleaser release --snapshot --rm-dist

release:
	goreleaser release --rm-dist

docs:
	mkdocs build --clean

cover:
	go tool cover -html=coverage.txt -o tmp/coverage.html
