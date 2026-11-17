// Package graph define as representações de grafos usadas ao longo da disciplina.
package graph

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidVertex é retornado quando uma aresta referencia um vértice fora do intervalo [0, Order()).
	ErrInvalidVertex = errors.New("graph: vértice inválido")

	ErrSelfLoop = errors.New("graph: laço não permitido em grafo simples")
)

// Neighborhood é o mínimo que BFS e DFS precisam: quantidade de vértices e, para cada
// um, os vizinhos que podem ser seguidos (em digrafos, só os arcos de saída).
type Neighborhood interface {
	// Order retorna a quantidade de vértices (n).
	Order() int

	// Neighbors retorna os vizinhos de v SEMPRE em ordem crescente.
	Neighbors(v int) []int
}

// Graph é o contrato comum a todas as representações de grafo simples e não direcionado.
type Graph interface {
	Neighborhood

	// AddEdge insere a aresta {u, v}. Como o grafo é não direcionado,
	// a ligação deve ser registrada nos dois sentidos.
	AddEdge(u, v int) error

	HasEdge(u, v int) bool

	String() string
}

// Edge representa uma aresta {U, V} (ou o arco U → V em um digrafo).
type Edge struct {
	U, V int
}

func (e Edge) String() string {
	return fmt.Sprintf("(%d, %d)", e.U, e.V)
}

// Build monta o grafo com o construtor informado, permitindo criar o mesmo
// grafo nas duas representações sem repetir código.
func Build(newGraph func(n int) Graph, n int, edges []Edge) (Graph, error) {
	g := newGraph(n)
	for _, e := range edges {
		if err := g.AddEdge(e.U, e.V); err != nil {
			return nil, fmt.Errorf("aresta %v: %w", e, err)
		}
	}
	return g, nil
}

// Size retorna a quantidade de arestas (m): cada aresta aparece na vizinhança
// das duas pontas, por isso a soma dos graus é dividida por 2.
func Size(g Graph) int {
	degrees := 0
	for v := range g.Order() {
		degrees += len(g.Neighbors(v))
	}
	return degrees / 2
}

// Edges retorna as arestas do grafo com U < V, em ordem crescente.
func Edges(g Graph) []Edge {
	var edges []Edge
	for u := range g.Order() {
		for _, v := range g.Neighbors(u) {
			if u < v {
				edges = append(edges, Edge{u, v})
			}
		}
	}
	return edges
}

func validate(n, u, v int) error {
	if u < 0 || u >= n || v < 0 || v >= n {
		return fmt.Errorf("%w: (%d, %d) com n = %d", ErrInvalidVertex, u, v, n)
	}
	if u == v {
		return fmt.Errorf("%w: (%d, %d)", ErrSelfLoop, u, v)
	}
	return nil
}

// FormatInts formata uma sequência no estilo do enunciado: [1, 2, 3].
func FormatInts(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = fmt.Sprint(x)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
