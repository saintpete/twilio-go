.PHONY: test vet release

# would be great to make the bash location portable but not sure how
SHELL = /bin/bash

WRITE_MAILMAP := $(shell command -v write_mailmap)
BUMP_VERSION := $(shell command -v bump_version)
STATICCHECK := $(shell command -v staticcheck)

test: vet
	go test -short ./...

vet: 
ifndef STATICCHECK
	go get -u honnef.co/go/staticcheck/cmd/staticcheck
endif
	go vet ./...
	staticcheck ./...

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
