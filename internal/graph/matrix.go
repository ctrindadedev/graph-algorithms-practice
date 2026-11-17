package graph

import (
	"fmt"
	"strings"
)

// AdjacencyMatrix representa o grafo por uma matriz n × n, onde
// cells[i][j] == 1 indica que existe aresta entre i e j (Questão 1).
type AdjacencyMatrix struct {
	n     int
	cells [][]int
}

// Garante em tempo de compilação que AdjacencyMatrix implementa Graph.
var _ Graph = (*AdjacencyMatrix)(nil)

func NewAdjacencyMatrix(n int) *AdjacencyMatrix {
	// PASSO 1 – criação e inicialização: make zera cada posição, ou seja,
	// inicialmente não há ligação entre nenhum par (i, j).
	cells := make([][]int, n)
	for i := range cells {
		cells[i] = make([]int, n)
	}
	return &AdjacencyMatrix{n: n, cells: cells}
}

func (m *AdjacencyMatrix) Order() int {
	return m.n
}

func (m *AdjacencyMatrix) AddEdge(u, v int) error {
	if err := validate(m.n, u, v); err != nil {
		return err
	}
	// PASSO 2 – inserção: registra u → v e, por ser não direcionado, v → u.
	m.cells[u][v] = 1
	m.cells[v][u] = 1
	return nil
}

func (m *AdjacencyMatrix) HasEdge(u, v int) bool {
	if validate(m.n, u, v) != nil {
		return false
	}
	return m.cells[u][v] == 1
}

func (m *AdjacencyMatrix) Neighbors(v int) []int {
	// Percorrer a linha v com j crescente já devolve os vizinhos ordenados.
	var out []int
	for j, cell := range m.cells[v] {
		if cell == 1 {
			out = append(out, j)
		}
	}
	return out
}

func (m *AdjacencyMatrix) String() string {
	// PASSO 3 – exibição: uma linha da matriz por vértice.
	var b strings.Builder
	for _, row := range m.cells {
		for j, cell := range row {
			if j > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprint(&b, cell)
		}
		b.WriteByte('\n')
	}
	return b.String()
}
