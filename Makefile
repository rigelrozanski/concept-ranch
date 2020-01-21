
build:
	go build ./...

install:
	go install ./...

test:
	go test ./...

.PHONY: build install test
