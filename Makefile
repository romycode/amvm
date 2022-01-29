GOCOVER=go tool cover
GOTEST=go test

install:
	go install ./cmd/mvm

test:
	$(GOTEST) ./pkg/...

test/cover:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out

build/linux:
	@ mkdir -p dist/linux
	@ GOOS=linux GOARCH=amd64 go build -o ./dist/linux/mvm ./cmd/mvm

build/macos:
	@ mkdir -p dist/macos
	@ GOOS=darwin GOARCH=amd64 go build -o ./dist/macos/mvm ./cmd/mvm

 build: build/linux build/macos