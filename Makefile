.PHONY: test vet release

test: vet
	go test -short ./...

vet: 
	go vet ./...

race-test: vet
	go test -race ./...

release: race-test
	bump_version minor http.go
