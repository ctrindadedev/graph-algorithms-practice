.PHONY: help run build test cover fmt vet check clean

# Uso: make run LISTA=lista02 ARGS="-start 3"
LISTA ?= lista01
ARGS  ?=
BIN   := bin

.DEFAULT_GOAL := help

help:
	@echo "Alvos disponíveis:"
	@echo "  run    executa uma lista (LISTA=$(LISTA), ARGS=\"$(ARGS)\")"
	@echo "  build  compila todas as listas em $(BIN)/"
	@echo "  test   roda todos os testes"
	@echo "  cover  roda os testes com relatório de cobertura"
	@echo "  fmt    formata o código com gofmt"
	@echo "  vet    analisa o código com go vet"
	@echo "  check  vet + test"
	@echo "  clean  remove binários e relatórios de cobertura"
	@echo "Listas: $(notdir $(wildcard cmd/*))"

run:
	go run ./cmd/$(LISTA) $(ARGS)

build:
	go build -o $(BIN)/ ./cmd/...

test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

fmt:
	gofmt -w .

vet:
	go vet ./...

check: vet test

clean:
	rm -rf $(BIN) coverage.out
