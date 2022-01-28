GOCOVER=go tool cover
GOTEST=go test

build:
	go build -o dist/mvm ./cmd/mvm

install:
	go install ./cmd/mvm

test:
	$(GOTEST) ./pkg/...

test/cover:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out