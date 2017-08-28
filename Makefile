.PHONY: test vet release

# would be great to make the bash location portable but not sure how
SHELL = /bin/bash -o pipefail

DIFFER := $(GOPATH)/bin/differ
WRITE_MAILMAP := $(GOPATH)/bin/write_mailmap
BUMP_VERSION := $(GOPATH)/bin/bump_version
MEGACHECK := $(GOPATH)/bin/megacheck

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

$(BUMP_VERSION):
	go get github.com/Shyp/bump_version

$(DIFFER):
	go get github.com/kevinburke/differ

$(MEGACHECK):
	go get honnef.co/go/tools/cmd/megacheck

$(WRITE_MAILMAP):
	go get github.com/kevinburke/write_mailmap

vet: | $(MEGACHECK)
	go vet ./...
	$(MEGACHECK) --ignore='github.com/kevinburke/twilio-go/*.go:S1002' ./...

race-test: vet
	go test -race ./...

release: race-test | $(DIFFER) $(BUMP_VERSION)
	$(DIFFER) $(MAKE) authors
	$(BUMP_VERSION) minor http.go

force: ;

AUTHORS.txt: force | $(WRITE_MAILMAP)
	$(WRITE_MAILMAP) > AUTHORS.txt

authors: AUTHORS.txt
