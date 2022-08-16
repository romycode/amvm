GO_COVER=go tool cover
GO_TEST=go test

install:
	go install ./cmd/amvm

test:
	$(GOTEST) ./pkg/...

test/cover:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out

build/linux:
	@ mkdir -p dist/linux
	@ GOOS=linux GOARCH=amd64 go build -o ./dist/linux/amvm ./cmd/amvm

build/macos:
	@ mkdir -p dist/macos
	@ GOOS=darwin GOARCH=amd64 go build -o ./dist/macos/amvm ./cmd/amvm

build: build/linux build/macos
