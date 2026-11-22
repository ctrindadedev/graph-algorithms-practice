<div align="center">

# Graph Algorithms Practice

**Implementações da disciplina de Grafos em Go**

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![UFRN](https://img.shields.io/badge/UFRN-DIMAP-003C7D?style=for-the-badge)
![Disciplina](https://img.shields.io/badge/Grafos-📚-6E40C9?style=for-the-badge)
![Status](https://img.shields.io/badge/status-em_andamento-yellow?style=for-the-badge)

</div>

---

## 🎯 Objetivo

Este repositório reúne as atividades práticas da disciplina de **Grafos** (UFRN/DIMAP)

Além de cumprir as listas, o projeto serve como estudo de **Go** e de boas práticas de engenharia de software:

- 🧩 **Código reutilizável** — representações e algoritmos ficam em pacotes compartilhados, e cada atividade só consome o que precisa.
- 🧪 **Testes como especificação** — as saídas esperadas do enunciado viram testes automatizados.
- 📝 **Código comentado** — cada trecho relevante explica o que acontece, como pedido nas atividades.
- 📈 **Evolução incremental** — novas listas acrescentam pacotes e comandos sem reescrever o que já existe.

## 🗂️ Estrutura

```
.
├── cmd/                      # um executável por atividade (package main)
│   ├── unidade1/
│   │   ├── lista01/          # matriz, lista, BFS e DFS
│   │   ├── lista02/          # árvores, árvores geradoras e arborescência
│   │   └── lista03/          # BFS/DFS em digrafo, ponderado e Dijkstra
│   └── unidade2/             # estudos dirigidos de caminhos mínimos (a fazer)
│       ├── dijkstra/
│       ├── bellman_ford/
│       └── floyd_warshall/
├── internal/                 # código compartilhado entre as atividades
│   ├── graph/                # interface Graph, matriz, lista, digrafo e lista ponderada
│   ├── traversal/            # BFS, DFS recursiva e DFS com pilha
│   ├── tree/                 # EhArvore, árvores geradoras BFS/DFS e arborescência
│   ├── shortest/             # caminho mínimo (Dijkstra)
│   ├── ds/                   # estruturas auxiliares genéricas (fila, pilha)
│   └── cli/                  # leitura de entrada dos menus
└── docs/                     # enunciados e respostas, na mesma divisão de cmd/
    ├── unidade1/
    │   ├── lista01/
    │   ├── lista02/
    │   └── lista03/
    └── unidade2/
        ├── dijkstra/
        ├── bellman_ford/
        └── floyd_warshall/
```

## 📋 Atividades

### Unidade I

| Lista | Conteúdo | Respostas |
|-------|----------|-----------|
| `lista01` | Matriz e lista de adjacência, BFS, DFS recursiva e DFS com pilha | [docs/unidade1/lista01/respostas.md](docs/unidade1/lista01/respostas.md) |
| `lista02` | Árvores, árvores geradoras por BFS/DFS e arborescência (Aula 9) | [docs/unidade1/lista02/respostas.md](docs/unidade1/lista02/respostas.md) |
| `lista03` | BFS e DFS em grafos direcionados, ponderados e direcionados ponderados; Dijkstra | [docs/unidade1/lista03/respostas.md](docs/unidade1/lista03/respostas.md) |

### Unidade II – Caminhos mínimos

Avaliação da unidade: três estudos dirigidos, entrega pelo SIGAA até **01/12/2026, 23h59**.

| Estudo dirigido | Conteúdo | Status |
|-----------------|----------|--------|
| `dijkstra` | Caminhos mínimos em grafos ponderados com pesos não negativos | a fazer |
| `bellman_ford` | Caminhos mínimos com pesos negativos (e detecção de ciclo negativo) | a fazer |
| `floyd_warshall` | Caminhos mínimos entre todos os pares de vértices | a fazer |

## 🚀 Como executar

> Pré-requisito: [Go 1.26+](https://go.dev/dl/)

```bash
# rodar uma atividade
go run ./cmd/unidade1/lista01           # grafo do enunciado
go run ./cmd/unidade1/lista01 -extra    # com a aresta (2,5) (Questão 10)
go run ./cmd/unidade1/lista01 -input    # grafo digitado pelo usuário (Questão 13)
go run ./cmd/unidade1/lista02
make run LISTA=lista02 ARGS="-start 3"
go run ./cmd/unidade1/lista03           # origem A
go run ./cmd/unidade1/lista03 -start D  # outra origem (mostra inalcançáveis)

# compilar todas as listas em bin/
make build

# rodar todos os testes
go test ./...
make check
```

### Makefile

| Alvo | O que faz |
|------|-----------|
| `make` / `make help` | Lista os alvos e as listas disponíveis |
| `make run` | Executa uma lista (`LISTA=lista01` por padrão; `ARGS` repassa as flags) |
| `make build` | Compila todas as listas em `bin/` |
| `make test` | Roda todos os testes |
| `make cover` | Roda os testes com relatório de cobertura |
| `make fmt` | Formata o código com `gofmt` |
| `make vet` | Análise estática com `go vet` |
| `make check` | `vet` + `test` |
| `make clean` | Remove `bin/` e `coverage.out` |

```bash
make run                                # lista01, grafo do enunciado
make run ARGS="-extra"                  # lista01 com a aresta (2,5) (Questão 10)
make run ARGS="-input"                  # lista01 com grafo digitado (Questão 13)
make run LISTA=lista02                  # lista02
make run LISTA=lista02 ARGS="-start 3"  # lista02 com outra raiz
make run LISTA=lista03                  # lista03
```

## 🛠️ Tecnologias

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![Make](https://img.shields.io/badge/Make-427819?style=flat-square&logo=gnu&logoColor=white)
![Git](https://img.shields.io/badge/Git-F05032?style=flat-square&logo=git&logoColor=white)

## 👤 Autor

**Caio Trindade** — [@ctrindadedev](https://github.com/ctrindadedev)
