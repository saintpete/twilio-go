vet: 
	go vet ./...

test: vet
	go test -short ./...

race-test: vet
	go test -race ./...

release: race-test
	bump_version minor http.go
