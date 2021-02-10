target: ;

test:
	go test ./... -cover -coverprofile=coverage.txt
