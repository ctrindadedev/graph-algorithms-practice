package graph

import (
	"fmt"
	"slices"
	"strings"
)

// AdjacencyList representa o grafo por uma coleção de vizinhos para cada
// vértice (Questão 2): adj[i] contém os vizinhos de i.
type AdjacencyList struct {
	adj [][]int
}

// Garante em tempo de compilação que AdjacencyList implementa Graph.
var _ Graph = (*AdjacencyList)(nil)

func NewAdjacencyList(n int) *AdjacencyList {
	// PASSO 1 – criação: cada vértice começa sem vizinhos registrados.
	return &AdjacencyList{adj: make([][]int, n)}
}

func (l *AdjacencyList) Order() int {
	return len(l.adj)
}

func (l *AdjacencyList) AddEdge(u, v int) error {
	if err := validate(l.Order(), u, v); err != nil {
		return err
	}
	// PASSO 2 – inserção nos dois sentidos: v vira vizinho de u e u vira
	// vizinho de v. A inserção é ordenada e ignora arestas repetidas.
	l.adj[u] = insertSorted(l.adj[u], v)
	l.adj[v] = insertSorted(l.adj[v], u)
	return nil
}

func insertSorted(xs []int, x int) []int {
	i, found := slices.BinarySearch(xs, x)
	if found {
		return xs
	}
	return slices.Insert(xs, i, x)
}

func (l *AdjacencyList) HasEdge(u, v int) bool {
	if validate(l.Order(), u, v) != nil {
		return false
	}
	_, found := slices.BinarySearch(l.adj[u], v)
	return found
}

// Neighbors devolve uma cópia para que quem chama não altere a estrutura interna.
func (l *AdjacencyList) Neighbors(v int) []int {
	return slices.Clone(l.adj[v])
}

func (l *AdjacencyList) String() string {
	// PASSO 3 – exibição no formato "0 → [1, 2]".
	var b strings.Builder
	for i, neighbors := range l.adj {
		fmt.Fprintf(&b, "%d → %s\n", i, FormatInts(neighbors))
	}
	return b.String()
}
