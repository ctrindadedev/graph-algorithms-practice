.PHONY: run test cover fmt vet check

# Uso: make run LISTA=lista01 ARGS="-extra"
LISTA ?= lista01
ARGS  ?=

run:
	go run ./cmd/$(LISTA) $(ARGS)

test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

fmt:
	gofmt -w .

vet:
	go vet ./...

check: vet test
