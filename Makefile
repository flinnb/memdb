.PHONY: run test


run:
	go run cmd/memdb.go

test:
	go mod download
	go install gotest.tools/gotestsum
	gotestsum ./...
