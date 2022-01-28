build:
	go build -o dist/mvm ./cmd/mvm

install:
	go install ./cmd/mvm

test:
	go test ./pkg/...