.PHONY: all build test vet fmt clean cross

VERSION ?= 1.0.0

all: build

build:
	go build -trimpath -o ignorekit .

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -s -w .

cross:
	VERSION=$(VERSION) ./scripts/build.sh

clean:
	rm -rf dist ignorekit ignorekit.exe
