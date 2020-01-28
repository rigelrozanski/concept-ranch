
build:
	go build ./quac/qu/...

install:
	go install ./quac/qu/...

test:
	go test ./quac/qu/...

.PHONY: build install test
