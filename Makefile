.PHONY: test vet release

# would be great to make the bash location portable but not sure how
SHELL = /bin/bash

DIFFER := $(shell command -v differ)
WRITE_MAILMAP := $(shell command -v write_mailmap)
BUMP_VERSION := $(shell command -v bump_version)
MEGACHECK := $(shell command -v megacheck)

test: vet
	bazel test --test_output=errors --test_arg="-test.short" //...

ci:
	bazel --batch --host_jvm_args=-Dbazel.DigestFunction=SHA1 test \
		--experimental_repository_cache="$$HOME/.bzrepos" \
		--spawn_strategy=remote \
		--remote_rest_cache=https://remote.rest.stackmachine.com/cache \
		--test_output=errors \
		--strategy=Javac=remote \
		--profile=profile.out \
		--test_arg="-test.short" \
		--features=race //... 2>&1 | ts '[%Y-%m-%d %H:%M:%.S]'
	bazel analyze-profile --curses=no --noshow_progress profile.out

vet:
ifndef MEGACHECK
	go get -u honnef.co/go/tools/cmd/megacheck
endif
	go vet ./...
	megacheck --ignore='github.com/kevinburke/twilio-go/*.go:S1002' ./...

race-test: vet
	go test -race ./...

release: race-test
ifndef DIFFER
	go get -u github.com/kevinburke/differ
endif
	differ $(MAKE) authors
ifndef BUMP_VERSION
	go get github.com/Shyp/bump_version
endif
	bump_version minor http.go

authors:
ifndef WRITE_MAILMAP
	go get github.com/kevinburke/write_mailmap
endif
	write_mailmap > AUTHORS.txt
