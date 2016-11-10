.PHONY: test vet release

WRITE_MAILMAP := $(shell command -v write_mailmap)
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

authors:
ifndef WRITE_MAILMAP
	go get github.com/kevinburke/write_mailmap
endif
	write_mailmap > AUTHORS.txt
