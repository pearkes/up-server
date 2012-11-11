build:
	go build -o bin/up

run:
	bin/up run

test:
	go test

.PHONY: build run test
