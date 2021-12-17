.ONESHELL:
.PHONY:

## test-only: run unit/mock tests without any dependent step
test-only:
	go test -v -cover ./...

## test: run unit/mock tests
test: generate
	go test -v -cover ./...

## test-race: run go unit/mock tests with race detection
test-race: generate
	go test -v -race ./...

## cover-report: run coverage and show html report
cover-report: cover
	go tool cover -html=coverage.nomocks.out

## cover: run unit/mock tests with coverage report. Generated mocks are filtered out of the report
cover: generate
	go test -v --race -coverprofile=coverage.out -coverpkg=./... ./...
	cat coverage.out | grep -v "mock" > coverage.nomocks.out
	go tool cover -func coverage.nomocks.out

## lint: runs golangci-lint
lint: generate
	golangci-lint run -v ./...

## lint-only: runs golangci-lint without any dependent step
lint-only:
	golangci-lint run -v ./...

## generate: runs go generate
generate:
	go generate -v ./...

## clean-mock: removes all generated mocks
clean-mock:
	find . -iname '*_mock.go' -exec rm {} \;

## regenerate: clear and regenerate mocks
regenerate: clean-mock generate

## update: runs go mod vendor and tidy
update: mod tidy

## mod: runs go mod vendor
mod:
	go mod vendor -v

## tidy: runs go mod tidy
tidy:
	go mod tidy -v

## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

