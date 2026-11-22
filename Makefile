.PHONY: help run build test cover fmt vet check clean

# Uso: make run LISTA=arvores_geradoras ARGS="-start 3"
LISTA ?= matriz_lista_bfs_dfs
ARGS  ?=
BIN   := bin

# Cada atividade fica em cmd/<unidade>/<lista>; LISTA é só o nome da pasta.
LISTAS := $(notdir $(patsubst %/,%,$(dir $(wildcard cmd/*/*/main.go))))
CMD    := $(firstword $(wildcard cmd/*/$(LISTA)))

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
	@echo "Listas: $(LISTAS)"

run:
	@test -n "$(CMD)" || { echo "lista '$(LISTA)' não encontrada em cmd/*/"; exit 1; }
	go run ./$(CMD) $(ARGS)

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
