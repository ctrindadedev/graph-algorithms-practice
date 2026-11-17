package tree

import (
	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/traversal"
)

// Report reúne os critérios de verificação de uma árvore geradora.
type Report struct {
	Vertices     int  // vértices alcançados a partir da raiz usando só arestas de T
	Edges        int  // |T|
	HasVMinusOne bool // |T| = |V| - 1
	EdgesInGraph bool // toda aresta de T pertence ao grafo original
	Acyclic      bool // T não forma ciclo
	Height       int  // maior distância, em arestas, da raiz até um vértice de T
	MaxWidth     int  // maior quantidade de vértices em um mesmo nível
}

// Valid informa se T é uma árvore geradora de um grafo com n vértices.
func (r Report) Valid(n int) bool {
	return r.Vertices == n && r.HasVMinusOne && r.EdgesInGraph && r.Acyclic
}

// Check avalia T como árvore geradora de g com raiz root.
func Check(g graph.Graph, root int, t []graph.Edge) Report {
	n := g.Order()
	r := Report{Edges: len(t), HasVMinusOne: len(t) == n-1, EdgesInGraph: true}

	tg := graph.NewAdjacencyList(n)
	duplicated := false
	for _, e := range t {
		if !g.HasEdge(e.U, e.V) {
			r.EdgesInGraph = false
		}
		if tg.HasEdge(e.U, e.V) {
			duplicated = true
		}
		tg.AddEdge(e.U, e.V)
	}

	// Uma floresta com c componentes sobre n vértices tem exatamente n - c
	// arestas; qualquer aresta a mais fecha um ciclo.
	components := 0
	seen := make([]bool, n)
	for v := range n {
		if !seen[v] {
			components++
			for _, u := range traversal.BFS(tg, v, nil) {
				seen[u] = true
			}
		}
	}
	r.Acyclic = !duplicated && len(t) == n-components

	// Os passos da BFS na árvore a partir da raiz fornecem os níveis.
	level := make([]int, n)
	widths := map[int]int{}
	traversal.BFS(tg, root, func(s traversal.Step) {
		r.Vertices++
		widths[level[s.Vertex]]++
		r.Height = max(r.Height, level[s.Vertex])
		for _, v := range s.Added {
			level[v] = level[s.Vertex] + 1
		}
	})
	for _, w := range widths {
		r.MaxWidth = max(r.MaxWidth, w)
	}
	return r
}
