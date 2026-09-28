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
├── cmd/                  # um executável por atividade (package main)
│   └── lista01/
├── internal/             # código compartilhado entre as atividades
│   ├── graph/            # interface Graph + matriz e lista de adjacência
│   ├── traversal/        # BFS, DFS recursiva e DFS com pilha
│   └── ds/               # estruturas auxiliares genéricas (fila, pilha)
└── docs/                 # enunciados e respostas teóricas de cada atividade
    └── lista01/
```

## 🚀 Como executar

> Pré-requisito: [Go 1.26+](https://go.dev/dl/)

```bash
# rodar uma atividade
go run ./cmd/lista01
make run LISTA=lista01 ARGS="-extra"

# rodar todos os testes
go test ./...
make check
```

## 🛠️ Tecnologias

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![Make](https://img.shields.io/badge/Make-427819?style=flat-square&logo=gnu&logoColor=white)
![Git](https://img.shields.io/badge/Git-F05032?style=flat-square&logo=git&logoColor=white)

## 👤 Autor

**Caio Trindade** — [@ctrindadedev](https://github.com/ctrindadedev)
