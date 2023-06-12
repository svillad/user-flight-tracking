.PHONY: build run test

PROJECT?=user-flight-tracking

default: build

build: test build-local

build-local:
	go build -o ./app ./server

run: build
	./app

run-local: build-local
	./app

lint: get-linter
	golangci-lint run --timeout=5m

get-linter:
	command -v golangci-lint || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

test: generate
	go fmt ./...
	go test -vet all ./... --cover

cover:
	echo unit tests only
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

generate: remove-vendor get-generator
	go generate -x ./...
	# go mod vendor

get-generator:
	go install github.com/golang/mock/mockgen

regenerate: clean-mock generate

clean-mock:
	find mocks -iname '*._mock.go' -exec rm {} \;

go.mod:
	go mod init
	go mod tidy

vendor:
	go mod vendor -v

remove-vendor:
	rm -rf vendor/
