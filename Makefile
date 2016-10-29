.PHONY: test vet release
BUMP_VERSION := $(shell command -v bump_version)

test: vet
	go test -short ./...

vet: 
	go vet ./...

race-test: vet
	go test -race ./...

release: race-test
ifndef BUMP_VERSION
	go get github.com/Shyp/bump_version
endif
	bump_version minor http.go
