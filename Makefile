GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

test: tidy
	go test -race -v -coverprofile c.out ./...

tidy: goimports
	test -z "$$(goimports -l -d $(GO_FILES) | tee /dev/stderr)"
	go vet ./...

goimports:
	go get golang.org/x/tools/cmd/goimports
