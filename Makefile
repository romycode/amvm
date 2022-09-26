install:
	@ go install ./cmd/amvm

test:
	@ go test -count=1 ./pkg/...

test/cover:
	@ go test -v -coverprofile=coverage.out ./...
	@ go tool cover -func=coverage.out
	@ go tool cover -html=coverage.out

build/linux:
	@ mkdir -p dist/linux
	@ GOOS=linux GOARCH=amd64 go build -o ./dist/linux/amvm ./cmd/amvm

build/macos:
	@ mkdir -p dist/macos
	@ GOOS=darwin GOARCH=amd64 go build -o ./dist/macos/amvm ./cmd/amvm

build: build/linux build/macos
