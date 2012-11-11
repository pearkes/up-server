dir:
	if [ ! -d "bin" ]; then mkdir "bin"; fi

build: dir
	echo "Building..."
	go build -o bin/up
	echo "Built to: './bin/up'"

run:
	bin/up run

test:
	go test

.PHONY: build run test

.SILENT: dir build run test
