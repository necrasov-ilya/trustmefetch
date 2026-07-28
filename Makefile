.PHONY: build test check install clean

build:
	go build -o bin/trustmefetch ./cmd/trustmefetch

test:
	go test ./...

check:
	gofmt -w cmd internal
	go vet ./...
	go test ./...

install:
	./scripts/install-local.sh

clean:
	rm -rf bin dist

